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

// DoRequestFunc is an alias for any function that takes an http request and returns a response and error
type DoRequestFunc func(request *http.Request) (*http.Response, error)

// OpenIDStrategy implementation of OpenID strategy
type OpenIDStrategy struct {
	*oauth2.Config
	Options Options
}

type openIDConfiguration struct {
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

// Options is the configuration used to construct a new OIDC authenticator
type Options struct {
	IssuerURL string
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

type TokenRefresher struct {
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	HTTPClient            *http.Client
}

func NewTokenRefresher(clientID, clientSecret, authorizationEndpoint, tokenEndpoint string, httpClient *http.Client) auth.TokenRefresher {
	return &TokenRefresher{
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		AuthorizationEndpoint: authorizationEndpoint,
		TokenEndpoint:         tokenEndpoint,
		HTTPClient:            httpClient,
	}
}

// Refresh tries to refresh the access token if it has expired and refresh token is provided
func (r *TokenRefresher) Refresh(old auth.Token) (*auth.Token, error) {
	config := &oauth2.Config{
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  r.AuthorizationEndpoint,
			TokenURL: r.TokenEndpoint,
		},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, r.HTTPClient)
	token := &oauth2.Token{
		AccessToken:  old.AccessToken,
		RefreshToken: old.RefreshToken,
		Expiry:       old.ExpiresIn,
		TokenType:    old.TokenType,
	}
	refresher := config.TokenSource(ctx, token)
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

func (r *TokenRefresher) Client(reuseToken *auth.Token) *http.Client {
	config := &oauth2.Config{
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  r.AuthorizationEndpoint,
			TokenURL: r.TokenEndpoint,
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, r.HTTPClient)
	token := &oauth2.Token{
		AccessToken:  reuseToken.AccessToken,
		RefreshToken: reuseToken.RefreshToken,
		Expiry:       reuseToken.ExpiresIn,
		TokenType:    reuseToken.TokenType,
	}

	return config.Client(ctx, token)
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
