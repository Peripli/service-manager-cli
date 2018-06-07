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
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/spf13/viper"
)

// Configuration should be implemented for load and save of SM client config
// go:generate counterfeiter . Configuration
type Configuration interface {
	Save(*smclient.ClientConfig) error
	Load() (*smclient.ClientConfig, error)
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

	return &smConfiguration{viperEnv}, nil
}

// Save implements configuration save
func (smCfg *smConfiguration) Save(clientCfg *smclient.ClientConfig) error {
	smCfg.viperEnv.Set("url", clientCfg.URL)
	smCfg.viperEnv.Set("user", clientCfg.User)

	smCfg.viperEnv.Set("access_token", clientCfg.AccessToken)
	smCfg.viperEnv.Set("refresh_token", clientCfg.RefreshToken)
	smCfg.viperEnv.Set("expiry", clientCfg.Expiry)

	smCfg.viperEnv.Set("client_id", clientCfg.ClientID)
	smCfg.viperEnv.Set("client_secret", clientCfg.ClientSecret)
	smCfg.viperEnv.Set("token_url", clientCfg.TokenURL)
	smCfg.viperEnv.Set("auth_url", clientCfg.AuthURL)

	return smCfg.viperEnv.WriteConfig()
}

// Load implements configuration load
func (smCfg *smConfiguration) Load() (*smclient.ClientConfig, error) {
	if err := smCfg.viperEnv.ReadInConfig(); err != nil {
		return nil, err
	}

	clientConfig := &smclient.ClientConfig{}

	// fmt.Println(">>>>>", smCfg.viperEnv.AllSettings())
	if err := smCfg.viperEnv.Unmarshal(&clientConfig); err != nil {
		return nil, err
	}

	if err := clientConfig.Validate(); err != nil {
		return nil, err
	}

	return clientConfig, nil
}
