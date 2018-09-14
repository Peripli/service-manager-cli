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
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/httputil"
)

// Options is used to configure new authenticators and clients
type Options struct {
	ClientID              string `mapstructure:"client_id"`
	ClientSecret          string `mapstructure:"client_secret"`
	AuthorizationEndpoint string `mapstructure:"authorization_endpoint"`
	TokenEndpoint         string `mapstructure:"token_endpoint"`
	IssuerURL             string `mapstructure:"issuer_url"`

	HTTP *httputil.HTTPConfig
}

// Token contains the structure of a typical UAA response token
type Token struct {
	AccessToken  string `mapstructure:"access_token" json:"access_token"`
	TokenType    string `mapstructure:"token_type" json:"token_type"`
	RefreshToken string `mapstructure:"refresh_token" json:"refresh_token"`
	ExpiresIn    string `mapstructure:"expires_in" json:"expires_in"`
	Scope        string `mapstructure:"scope" json:"scope"`
}

// AuthenticationStrategy should be implemented for different authentication strategies
//go:generate counterfeiter . AuthenticationStrategy
type AuthenticationStrategy interface {
	Authenticate(user, password string) (*Token, error)
}

// Refresher should be implemented for refreshing access tokens with refresh token flow
//go:generate counterfeiter . Refresher
type Refresher interface {
	Token() (*Token, error)
}

// Client should be implemented for http like clients which do automatic authentication
//go:generate counterfeiter . Client
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
