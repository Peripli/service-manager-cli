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
	"path/filepath"

	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const defaultConfigFileName = ".servicemanager.json"

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
func NewSMConfiguration(cfgFile string) (Configuration, error) {
	viperEnv := viper.New()

	absCfgFilePath, err := getConfigFileAbsPath(cfgFile)
	if err != nil {
		return nil, err
	}
	viperEnv.SetConfigFile(absCfgFilePath)

	return &smConfiguration{viperEnv}, nil
}

// Save implements configuration save
func (smCfg *smConfiguration) Save(clientCfg *smclient.ClientConfig) error {
	smCfg.viperEnv.Set("url", clientCfg.URL)
	smCfg.viperEnv.Set("user", clientCfg.User)
	smCfg.viperEnv.Set("token", clientCfg.Token)

	return smCfg.viperEnv.WriteConfig()
}

// Load implements configuration load
func (smCfg *smConfiguration) Load() (*smclient.ClientConfig, error) {
	if err := smCfg.viperEnv.ReadInConfig(); err != nil {
		return nil, err
	}

	clientConfig := &smclient.ClientConfig{}

	if err := smCfg.viperEnv.Unmarshal(clientConfig); err != nil {
		return nil, err
	}

	if err := clientConfig.Validate(); err != nil {
		return nil, err
	}

	return clientConfig, nil
}

func getConfigFileAbsPath(cfgFile string) (string, error) {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		cfgFile = filepath.Join(home, defaultConfigFileName)
	}

	filename, err := filepath.Abs(cfgFile)
	if err != nil {
		return "", err
	}

	return filename, nil
}
