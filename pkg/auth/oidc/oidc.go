/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"golang.org/x/oauth2"
)

type openIDConfiguration struct {
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

// DoRequestFunc is an alias for any function that takes an http request and returns a response and error
type DoRequestFunc func(request *http.Request) (*http.Response, error)

// OpenIDStrategy implementation of OpenID strategy
type OpenIDStrategy struct {
	*oauth2.Config
	Options Options
}

// Options is the configuration used to construct a new OIDC authentication strategy and token refresher
type Options struct {
	// IssuerURL is the endpoint which to call for token acquisition and other oauth configurations
	IssuerURL string

	// AuthorizationEndpoint is the oauth endpoint for authorization.
	// If this property is not set the IssuerURL should be set in order to fetch the configuration
	AuthorizationEndpoint string

	// TokenEndpoint is the oauth endpoint for fetching a token.
	// If this property is not set the IssuerURL should be set in order to fetch the configuration
	TokenEndpoint string

	// ClientID is the id of the oauth client used to verify the tokens
	ClientID string

	// ClientID is the id of the oauth client used to verify the tokens
	ClientSecret string

	HTTPClient *http.Client
}

// NewOpenIDStrategy returns OpenId auth strategy
func NewOpenIDStrategy(options Options) (auth.AuthenticationStrategy, *openIDConfiguration) {
	openIDConfig, err := fetchOpenidConfiguration(options.IssuerURL, options.HTTPClient.Do)
	if err != nil {
		panic(fmt.Errorf("Error occured while fetching openid configuration: %s", err))
	}

	config := &oauth2.Config{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  openIDConfig.AuthorizationEndpoint,
			TokenURL: openIDConfig.TokenEndpoint,
		},
	}

	return &OpenIDStrategy{
		Config:  config,
		Options: options,
	}, openIDConfig
}

// Authenticate is used to perform authentication action for OpenID strategy
func (s *OpenIDStrategy) Authenticate(user, password string) (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.Options.HTTPClient)
	token, err := s.PasswordCredentialsToken(ctx, user, password)
	if err != nil {
		return nil, err
	}

	resultToken := &auth.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry,
		TokenType:    token.TokenType,
	}

	return resultToken, err
}

type OIDCRefresher struct {
	*oauth2.Config
	// ClientID              string
	// ClientSecret          string
	// AuthorizationEndpoint string
	// TokenEndpoint         string
	HTTPClient *http.Client
}

func NewTokenRefresher(options Options) (auth.TokenRefresher, error) {
	if options.AuthorizationEndpoint == "" || options.TokenEndpoint == "" {
		openIDConfig, err := fetchOpenidConfiguration(options.IssuerURL, options.HTTPClient.Do)
		if err != nil {
			return nil, err
		}
		options.AuthorizationEndpoint = openIDConfig.AuthorizationEndpoint
		options.TokenEndpoint = openIDConfig.TokenEndpoint
	}

	return &OIDCRefresher{
		Config: &oauth2.Config{
			ClientID:     options.ClientID,
			ClientSecret: options.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  options.AuthorizationEndpoint,
				TokenURL: options.TokenEndpoint,
			},
		},
		HTTPClient: options.HTTPClient,
	}, nil
}

// Refresh tries to refresh the access token if it has expired and refresh token is provided
func (r *OIDCRefresher) Refresh(old auth.Token) (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, r.HTTPClient)
	token := &oauth2.Token{
		AccessToken:  old.AccessToken,
		RefreshToken: old.RefreshToken,
		Expiry:       old.ExpiresIn,
		TokenType:    old.TokenType,
	}
	refresher := r.TokenSource(ctx, token)
	refreshedToken, err := refresher.Token()
	if err != nil {
		return nil, err
	}

	old.AccessToken = refreshedToken.AccessToken
	old.RefreshToken = refreshedToken.RefreshToken
	old.ExpiresIn = refreshedToken.Expiry
	old.TokenType = refreshedToken.TokenType
	return &old, nil
}

func (r *OIDCRefresher) Client(reuseToken *auth.Token) *http.Client {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, r.HTTPClient)
	token := &oauth2.Token{
		AccessToken:  reuseToken.AccessToken,
		RefreshToken: reuseToken.RefreshToken,
		Expiry:       reuseToken.ExpiresIn,
		TokenType:    reuseToken.TokenType,
	}

	return r.Config.Client(ctx, token)
}

func fetchOpenidConfiguration(issuerURL string, readConfigurationFunc DoRequestFunc) (*openIDConfiguration, error) {
	req, err := http.NewRequest(http.MethodGet, issuerURL+"/.well-known/openid-configuration", nil)
	if err != nil {
		return nil, err
	}

	response, err := readConfigurationFunc(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Error getting OpenID configuration")
	}

	var configuration *openIDConfiguration
	if err = httputil.UnmarshalResponse(response, &configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}
