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
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// SMConfigKey is the configuration key to load/save Service Manager config
	SMConfigKey = "smclient"

	// HTTPConfigKey is the configuration key to load/save HTTP client config
	HTTPConfigKey = "httpconfig"
)

// Configuration should be implemented for load and save of SM client config
//go:generate counterfeiter . Configuration
type Configuration interface {
	UnmarshalKey(string, interface{}) error

	Set(string, interface{})
	Save(string, interface{}) error
}

type smConfiguration struct {
	*viper.Viper
}

// New returns implementation of Configuration interface
func New(viperEnv *viper.Viper, cfgFile string) (Configuration, error) {
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

	config := &smConfiguration{
		Viper: viperEnv,
	}

	if err := viperEnv.ReadInConfig(); err != nil {
		if err, ok := err.(*os.PathError); ok {
			// TODO: print somewhere else
			fmt.Println("Config File was not found: ", err)
			return config, nil
		}
		return nil, fmt.Errorf("could not read configuration cfg: %s", err)
	}

	return config, nil
}

// Save implements configuration save
func (smCfg *smConfiguration) Save(key string, value interface{}) error {
	smCfg.Set(key, value)

	return smCfg.WriteConfig()
}
