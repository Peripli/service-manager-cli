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
	"fmt"
	"net/http"

	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type openIDConfiguration struct {
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

// OpenIDStrategy implementation of OpenID strategy
type OpenIDStrategy struct {
	oauth2Config *oauth2.Config
	ccConfig     *clientcredentials.Config
	httpClient   *http.Client
}

// NewOpenIDStrategy returns OpenId auth strategy
func NewOpenIDStrategy(options *auth.Options) (*OpenIDStrategy, *auth.Options, error) {
	httpClient := util.BuildHTTPClient(options.SSLDisabled)
	httpClient.Timeout = options.Timeout

	openIDConfig, err := fetchOpenidConfiguration(options.IssuerURL, httpClient.Do)
	if err != nil {
		return nil, nil, fmt.Errorf("Error occurred while fetching openid configuration: %s", err)
	}
	options.AuthorizationEndpoint = openIDConfig.AuthorizationEndpoint
	options.TokenEndpoint = openIDConfig.TokenEndpoint

	oauthConfig := &oauth2.Config{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  options.AuthorizationEndpoint,
			TokenURL: options.TokenEndpoint,
		},
	}

	ccConfig := &clientcredentials.Config{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		TokenURL:     options.TokenEndpoint,
	}

	return &OpenIDStrategy{
		oauth2Config: oauthConfig,
		ccConfig:     ccConfig,
		httpClient:   httpClient,
	}, options, nil
}

// ClientCredentials is used to perform client credentials grant type flow
func (s *OpenIDStrategy) ClientCredentials() (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)
	token, err := s.ccConfig.Token(ctx)

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

// PasswordCredentials is used to perform password grant type flow
func (s *OpenIDStrategy) PasswordCredentials(user, password string) (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)
	token, err := s.oauth2Config.PasswordCredentialsToken(ctx, user, password)
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
