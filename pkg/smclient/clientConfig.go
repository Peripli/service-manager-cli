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

package smclient

import (
	"errors"

	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/auth"
)

// ClientConfig contains the configuration of the CLI.
type ClientConfig struct {
	auth.Token

	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string

	URL         string
	User        string
	SSLDisabled bool
}

// Validate validates client config
func (clientCfg ClientConfig) Validate() error {
	if err := util.ValidateURL(clientCfg.URL); err != nil {
		return err
	}
	if clientCfg.User == "" {
		return errors.New("User must not be empty")
	}
	if clientCfg.AccessToken == "" {
		return errors.New("Token must not be empty")
	}
	return nil
}

// GetToken returns the oauth token from the client configuration
func (clientCfg ClientConfig) GetToken() auth.Token {
	return clientCfg.Token
}
