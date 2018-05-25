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

package util

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Peripli/service-manager-cli/pkg/types"
)

// ValidateURL validates a URL
func ValidateURL(URL string) error {
	if URL == "" {
		return errors.New("url not provided")
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("url cannot be parsed: %s", err)
	}

	if !parsedURL.IsAbs() || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return fmt.Errorf("url is not an HTTP URL: %s", URL)
	}

	return nil
}

// GetBrokerByName returns array of brokers with the searched names
func GetBrokerByName(brokers *types.Brokers, names []string) []types.Broker {
	result := make([]types.Broker, 0)
	for _, broker := range brokers.Brokers {
		for _, name := range names {
			if broker.Name == name {
				result = append(result, broker)
			}
		}
	}
	return result
}

// GetPlatformByName returns array of platforms with the searched names
func GetPlatformByName(platforms *types.Platforms, names []string) []types.Platform {
	result := make([]types.Platform, 0)
	for _, platform := range platforms.Platforms {
		for _, name := range names {
			if platform.Name == name {
				result = append(result, platform)
			}
		}
	}
	return result
}
