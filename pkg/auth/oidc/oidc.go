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
	"github.com/Peripli/service-manager-cli/pkg/smclient"
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

	httpClient *http.Client
}

// NewOpenIDStrategy returns OpenId auth strategy
func NewOpenIDStrategy(config *smclient.ClientConfig, httpClient *http.Client) (auth.AuthenticationStrategy, *smclient.ClientConfig, error) {
	openIDConfig, err := fetchOpenidConfiguration(config.IssuerURL, httpClient.Do)
	if err != nil {
		return nil, nil, fmt.Errorf("Error occurred while fetching openid configuration: %s", err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  openIDConfig.AuthorizationEndpoint,
			TokenURL: openIDConfig.TokenEndpoint,
		},
	}

	config.AuthorizationEndpoint = openIDConfig.AuthorizationEndpoint
	config.TokenEndpoint = openIDConfig.TokenEndpoint

	return &OpenIDStrategy{
		Config:     oauthConfig,
		httpClient: httpClient,
	}, config, nil
}

// Authenticate is used to perform authentication action for OpenID strategy
func (s *OpenIDStrategy) Authenticate(user, password string) (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)
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

// Refresher implements TokenRefresher interface
type Refresher struct {
	*oauth2.Config
	HTTPClient *http.Client
}

// NewTokenRefresher returns new oidc token refresher
func NewTokenRefresher(config *smclient.ClientConfig, httpClient *http.Client) (auth.TokenRefresher, error) {
	if config.AuthorizationEndpoint == "" || config.TokenEndpoint == "" {
		openIDConfig, err := fetchOpenidConfiguration(config.IssuerURL, httpClient.Do)
		if err != nil {
			return nil, fmt.Errorf("Error occurred while fetching openid configuration: %s", err)
		}
		config.AuthorizationEndpoint = openIDConfig.AuthorizationEndpoint
		config.TokenEndpoint = openIDConfig.TokenEndpoint
	}

	return &Refresher{
		Config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthorizationEndpoint,
				TokenURL: config.TokenEndpoint,
			},
		},
		HTTPClient: httpClient,
	}, nil
}

// Refresh tries to refresh the access token if it has expired and refresh token is provided
func (r *Refresher) Refresh(old auth.Token) (*auth.Token, error) {
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

// Client returns http client with automatic token refreshing mechanism
func (r *Refresher) Client(reuseToken *auth.Token) *http.Client {
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
		return nil, errors.New("Unexpected status code")
	}

	var configuration *openIDConfiguration
	if err = httputil.UnmarshalResponse(response, &configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}
