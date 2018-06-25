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

// GetBrokersByName returns array of brokers with the searched names
func GetBrokersByName(brokers *types.Brokers, names []string) []types.Broker {
	result := make([]types.Broker, 0)
	brokersMap := make(map[string]types.Broker)
	for _, broker := range brokers.Brokers {
		brokersMap[broker.Name] = broker
	}

	for _, name := range names {
		if _, exists := brokersMap[name]; exists {
			result = append(result, brokersMap[name])
		}
	}
	return result
}

// GetPlatformsByName returns array of platforms with the searched names
func GetPlatformsByName(platforms *types.Platforms, names []string) []types.Platform {
	result := make([]types.Platform, 0)
	platformsMap := make(map[string]types.Platform)
	for _, platform := range platforms.Platforms {
		platformsMap[platform.Name] = platform
	}
	for _, name := range names {
		if _, exists := platformsMap[name]; exists {
			result = append(result, platformsMap[name])
		}
	}
	return result
}
