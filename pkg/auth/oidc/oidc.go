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
	token *oauth2.Token

	httpClient *http.Client
}

// NewOpenIDStrategy returns OpenId auth strategy
func NewOpenIDStrategy(config *auth.Options, httpClient *http.Client) (auth.AuthenticationStrategy, *auth.Options, error) {
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
		token:      nil,
	}, config, nil
}

// Authenticate is used to perform authentication action for OpenID strategy
func (s *OpenIDStrategy) Authenticate(user, password string) (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)
	token, err := s.PasswordCredentialsToken(ctx, user, password)
	if err != nil {
		return nil, err
	}

	s.token = token

	resultToken := &auth.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry,
		TokenType:    token.TokenType,
	}

	return resultToken, err
}

// Token returns access token, if it expires will refresh it and return it
func (s *OpenIDStrategy) Token() (*auth.Token, error) {
	if s.token == nil {
		return nil, errors.New("strategy is not authenticated")
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)

	refresher := s.Config.TokenSource(ctx, s.token)
	refreshedToken, err := refresher.Token()
	if err != nil {
		return nil, err
	}

	s.token = refreshedToken
	resultToken := &auth.Token{
		AccessToken:  refreshedToken.AccessToken,
		RefreshToken: refreshedToken.RefreshToken,
		ExpiresIn:    refreshedToken.Expiry,
		TokenType:    refreshedToken.TokenType,
	}

	return resultToken, nil
}

// Client returns http client with automatic token refreshing mechanism
func (s *OpenIDStrategy) Client() *http.Client {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)
	return s.Config.Client(ctx, s.token)
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
