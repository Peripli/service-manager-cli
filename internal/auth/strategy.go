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

package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"golang.org/x/oauth2"
)

const (
	defaultClientID     = "smctl"
	defaultClientSecret = "smctl"
)

// Token contains the structure of a typical UAA response token
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type openIDConfiguration struct {
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

// AuthenticationStrategy should be implemented for different authentication strategies
//go:generate counterfeiter . AuthenticationStrategy
type AuthenticationStrategy interface {
	Authenticate(ctx context.Context, issuerURL, user, password string) (*oauth2.Config, *oauth2.Token, error)
	RefreshToken(context.Context, oauth2.Config, oauth2.Token) (*oauth2.Token, error)
}

// NewOpenIDStrategy returns OpenId auth strategy
func NewOpenIDStrategy(requestFunc DoRequestFunc) AuthenticationStrategy {
	return &OpenIDStrategy{
		ReadConfigurationFunc: requestFunc,
	}
}

// OpenIDStrategy implementation of OpenID strategy
type OpenIDStrategy struct{
	// ReadConfigurationFunc is the function used to call the token issuer. If one is not provided, http.DefaultClient.Do will be used
	ReadConfigurationFunc DoRequestFunc
}

// DoRequestFunc is an alias for any function that takes an http request and returns a response and error
type DoRequestFunc func(request *http.Request) (*http.Response, error)

// Authenticate is used to perform authentication action for OpenID strategy
func (s *OpenIDStrategy) Authenticate(ctx context.Context, issuerURL, user, password string) (*oauth2.Config, *oauth2.Token, error) {
	endpoints, err := s.fetchOpenidConfiguration(issuerURL)
	if err != nil {
		return nil, nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     defaultClientID,
		ClientSecret: defaultClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  endpoints.AuthorizationEndpoint,
			TokenURL: endpoints.TokenEndpoint,
		},
	}

	token, err := oauth2Config.PasswordCredentialsToken(ctx, user, password)
	return oauth2Config, token, err
}

// RefreshToken tries to refresh the access token if it has expired and refresh token is provided
func (s *OpenIDStrategy) RefreshToken(ctx context.Context, config oauth2.Config, token oauth2.Token) (*oauth2.Token, error) {
	refresher := config.TokenSource(ctx, &token)
	return refresher.Token()
}

func (s *OpenIDStrategy) fetchOpenidConfiguration(issuerURL string) (*openIDConfiguration, error) {
	req, err := http.NewRequest(http.MethodGet, issuerURL+"/.well-known/openid-configuration", nil)
	if err != nil {
		return nil, err
	}

	response, err := s.ReadConfigurationFunc(req)
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
