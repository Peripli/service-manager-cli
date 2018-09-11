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
	"reflect"
	"strings"

	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
	"github.com/Peripli/service-manager/pkg/env"
	"github.com/fatih/structs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Settings struct {
	SMClient *smclient.ClientConfig
	HTTP     *httputil.HTTPConfig
}

func DefaultSettings() *Settings {
	return &Settings{
		SMClient: smclient.DefaultSettings(),
		HTTP:     httputil.DefaultHTTPConfig(),
	}
}

// Configuration should be implemented for load and save of SM client config
// go:generate counterfeiter . Configuration
type Configuration interface {
	env.Environment
	Save(interface{}) error
}

type smConfiguration struct {
	env.Environment
	viperEnv *viper.Viper
}

func DefaultConfigFile() env.File {
	fileDir, err := defaultFilePath()
	if err != nil {
		panic(fmt.Sprintf("Could not find home dir %s", err))
	}
	return env.File{
		Location: fileDir,
		Name:     "config",
		Format:   "json",
	}
}

func AddPFlags(set *pflag.FlagSet) {
	env.CreatePFlags(set, struct{ File env.File }{File: DefaultConfigFile()})
	env.CreatePFlags(set, DefaultSettings())
}

func New(env env.Environment) (*Settings, error) {
	config := DefaultSettings()
	if err := env.Unmarshal(config); err != nil {
		return nil, err
	}
	return config, nil
}

// NewSMConfiguration returns implementation of Configuration interface
func NewEnv(flags *pflag.FlagSet) (Configuration, error) {
	cfg := struct{ File File }{File: File{}}
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("could not find configuration cfg: %s", err)
	}
	ensureDirExists(flags)

	environment, err := env.New(flags)
	if err != nil {
		return nil, fmt.Errorf("Could not create environment: %s", err)
	}

	return &smConfiguration{
		viperEnv:    environment.Viper,
		Environment: environment,
	}, nil
}

// Save implements configuration save
func (smCfg *smConfiguration) Save(value interface{}) error {
	properties := make(map[string]interface{})
	traverseFields(value, "", properties)
	for key, value := range properties {
		smCfg.viperEnv.Set(key, value)
	}

	return smCfg.viperEnv.WriteConfig()
}

// traverseFields traverses the provided structure and prepares a slice of strings that contains
// the paths to the structure fields (nested paths in the provided structure use dot as a separator)
func traverseFields(value interface{}, buffer string, result map[string]interface{}) {
	if !structs.IsStruct(value) {
		index := strings.LastIndex(buffer, ".")
		if index == -1 {
			index = 0
		}
		key := strings.ToLower(buffer[0:index])
		result[key] = value
		return
	}

	s := structs.New(value)
	for _, field := range s.Fields() {
		if field.IsExported() && field.Kind() != reflect.Interface && field.Kind() != reflect.Func {
			var name string
			if field.Tag("mapstructure") != "" {
				name = field.Tag("mapstructure")
			} else {
				name = field.Name()
			}
			buffer += name + "."
			traverseFields(field.Value(), buffer, result)
			buffer = buffer[0:strings.LastIndex(buffer, name)]
		}
	}
}
