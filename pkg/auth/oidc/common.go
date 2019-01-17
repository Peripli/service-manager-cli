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
	"net/http"

	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// NewClient builds configured HTTP client.
//
// If token is provided will try to refresh the token if it has expired,
// otherwise if token is not provided will do client_credentials flow and fetch token
func NewClient(options *auth.Options, token *auth.Token) auth.Client {
	if !options.UseBasicAuth {
		oauth2.RegisterBrokenAuthHeaderProvider(options.IssuerURL)
	}

	httpClient := util.BuildHTTPClient(options.SSLDisabled)
	httpClient.Timeout = options.Timeout

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	var oauthClient *http.Client
	var tokenSource oauth2.TokenSource

	if token == nil || token.RefreshToken == "" {
		oauthConfig := &clientcredentials.Config{
			ClientID:     options.ClientID,
			ClientSecret: options.ClientSecret,
			TokenURL:     options.TokenEndpoint,
		}
		clientCredentialsSource := oauthConfig.TokenSource(ctx)
		// The double wrapping of TokenSource objects is needed, because there is no other way
		// to pass the existing access token and the client will try to fetch a token for each request
		tokenSource = oauth2.ReuseTokenSource(&oauth2.Token{
			AccessToken: token.AccessToken,
			Expiry:      token.ExpiresIn,
			TokenType:   token.TokenType,
		}, clientCredentialsSource)
	} else {
		oauthConfig := &oauth2.Config{
			ClientID:     options.ClientID,
			ClientSecret: options.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  options.AuthorizationEndpoint,
				TokenURL: options.TokenEndpoint,
			},
		}
		tokenSource = oauthConfig.TokenSource(ctx, &oauth2.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.ExpiresIn,
			TokenType:    token.TokenType,
		})
	}

	oauthClient = oauth2.NewClient(ctx, tokenSource)
	oauthClient.Timeout = options.Timeout

	return &Client{
		tokenSource: tokenSource,
		httpClient:  oauthClient,
	}
}

// Client is used to make http requests including bearer token automatically and refreshing it
// if necessary
type Client struct {
	tokenSource oauth2.TokenSource
	httpClient  *http.Client
}

// Do makes a http request with the underlying HTTP client which includes an access token in the request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// Token returns the token, refreshing it if necessary
func (c *Client) Token() (*auth.Token, error) {
	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return &auth.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry,
		TokenType:    token.TokenType,
	}, nil
}

// DoRequestFunc is an alias for any function that takes an http request and returns a response and error
type DoRequestFunc func(request *http.Request) (*http.Response, error)

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
		return nil, errors.New("unexpected status code")
	}

	var configuration *openIDConfiguration
	if err = httputil.UnmarshalResponse(response, &configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}
