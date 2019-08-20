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

package configuration

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Peripli/service-manager-cli/internal/util"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/spf13/viper"
)

// Settings contains the information that will be saved/loaded in the CLI config file
type Settings struct {
	auth.Token

	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	IssuerURL             string

	URL            string
	User           string
	TokenBasicAuth bool
	SSLDisabled    bool
}

// Validate validates client config
func (settings Settings) Validate() error {
	if err := util.ValidateURL(settings.URL); err != nil {
		return err
	}
	if settings.User == "" {
		return errors.New("user must not be empty")
	}
	if settings.AccessToken == "" {
		return errors.New("token must not be empty")
	}
	return nil
}

// GetToken returns the oauth token from the client configuration
func (settings Settings) GetToken() auth.Token {
	return settings.Token
}

// Configuration should be implemented for load and save of SM client config
//go:generate counterfeiter . Configuration
type Configuration interface {
	Save(*Settings) error
	Load() (*Settings, error)
}

type smConfiguration struct {
	viperEnv *viper.Viper
}

// NewSMConfiguration returns implementation of Configuration interface
func NewSMConfiguration(viperEnv *viper.Viper, cfgFile string) (Configuration, error) {
	if cfgFile == "" {
		var err error
		cfgFile, err = defaultFilePath()
		if err != nil {
			return nil, err
		}
	}
	if err := ensureDirExists(cfgFile); err != nil {
		return nil, err
	}

	viperEnv.SetConfigFile(cfgFile)
	viperEnv.SetDefault("token_basic_auth", true) // RFC 6749 section 2.3.1

	return &smConfiguration{viperEnv}, nil
}

// Save implements configuration save
func (smCfg *smConfiguration) Save(settings *Settings) error {
	smCfg.viperEnv.Set("url", settings.URL)
	smCfg.viperEnv.Set("user", settings.User)
	smCfg.viperEnv.Set("ssl_disabled", settings.SSLDisabled)
	smCfg.viperEnv.Set("token_basic_auth", settings.TokenBasicAuth)

	smCfg.viperEnv.Set("access_token", settings.AccessToken)
	smCfg.viperEnv.Set("refresh_token", settings.RefreshToken)
	smCfg.viperEnv.Set("expiry", settings.ExpiresIn.Format(time.RFC1123Z))

	smCfg.viperEnv.Set("client_id", settings.ClientID)
	smCfg.viperEnv.Set("client_secret", settings.ClientSecret)
	smCfg.viperEnv.Set("issuer_url", settings.IssuerURL)
	smCfg.viperEnv.Set("token_url", settings.TokenEndpoint)
	smCfg.viperEnv.Set("auth_url", settings.AuthorizationEndpoint)

	cfgFile := smCfg.viperEnv.ConfigFileUsed()
	if err := smCfg.viperEnv.WriteConfig(); err != nil {
		return fmt.Errorf("could not save config file %s: %s", cfgFile, err)
	}
	const ownerAccessOnly = 0600
	if err := os.Chmod(cfgFile, ownerAccessOnly); err != nil {
		return fmt.Errorf("could not set access rights of config file %s: %s", cfgFile, err)
	}
	return nil
}

// Load implements configuration load
func (smCfg *smConfiguration) Load() (*Settings, error) {
	if err := smCfg.viperEnv.ReadInConfig(); err != nil {
		return nil, err
	}

	settings := &Settings{}

	if err := smCfg.viperEnv.Unmarshal(&settings); err != nil {
		return nil, err
	}

	settings.SSLDisabled = smCfg.viperEnv.Get("ssl_disabled").(bool)
	settings.TokenBasicAuth = smCfg.viperEnv.Get("token_basic_auth").(bool)
	settings.AccessToken = smCfg.viperEnv.Get("access_token").(string)
	settings.RefreshToken = smCfg.viperEnv.Get("refresh_token").(string)
	settings.ExpiresIn, _ = time.Parse(time.RFC1123Z, smCfg.viperEnv.Get("expiry").(string))
	settings.TokenEndpoint = smCfg.viperEnv.Get("token_url").(string)
	settings.AuthorizationEndpoint = smCfg.viperEnv.Get("auth_url").(string)
	settings.IssuerURL = smCfg.viperEnv.Get("issuer_url").(string)
	settings.ClientID = smCfg.viperEnv.Get("client_id").(string)
	settings.ClientSecret = smCfg.viperEnv.Get("client_secret").(string)

	if err := settings.Validate(); err != nil {
		return nil, err
	}

	return settings, nil
}
