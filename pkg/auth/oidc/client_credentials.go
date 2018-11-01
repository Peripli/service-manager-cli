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

// ClientCredentialsStrategy implements client credentials flow authentication strategy
type ClientCredentialsStrategy struct {
	*clientcredentials.Config
	httpClient *http.Client
}

// NewClientCredentialsStrategy returns client credentials auth strategy
func NewClientCredentialsStrategy(options *auth.Options) (auth.AuthenticationStrategy, *auth.Options, error) {
	httpClient := util.BuildHTTPClient(options.SSLDisabled)
	httpClient.Timeout = options.Timeout

	openIDConfig, err := fetchOpenidConfiguration(options.IssuerURL, httpClient.Do)
	if err != nil {
		return nil, nil, fmt.Errorf("Error occurred while fetching openid configuration: %s", err)
	}
	options.AuthorizationEndpoint = openIDConfig.AuthorizationEndpoint
	options.TokenEndpoint = openIDConfig.TokenEndpoint

	oauthConfig := &clientcredentials.Config{
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		TokenURL:     options.TokenEndpoint,
	}

	return &ClientCredentialsStrategy{
		Config:     oauthConfig,
		httpClient: httpClient,
	}, options, nil
}

// Authenticate is used to perform authentication action for client credentials strategy
func (s *ClientCredentialsStrategy) Authenticate(user, password string) (*auth.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, s.httpClient)
	token, err := s.Token(ctx)

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
