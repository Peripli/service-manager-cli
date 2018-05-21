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
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

//go:generate counterfeiter . Client
type Client interface {
	RegisterPlatform(platform *types.Platform) (*types.Platform, error)
	RegisterBroker(broker *types.Broker) (*types.Broker, error)
	ListBrokers() (*types.Brokers, error)
}

type ServiceManagerClient struct {
	config     *ClientConfig
	httpClient *http.Client
	headers    *http.Header
}

func NewClient(config *ClientConfig) Client {
	client := &ServiceManagerClient{config: config, httpClient: &http.Client{}}
	client.headers = &http.Header{}
	client.headers.Add("Content-Type", "application/json")
	if len(client.config.Token) > 0 {
		client.headers.Add("Authorization", client.config.Token)
	}

	return client
}

// Registers a platform in the service manager
func (client *ServiceManagerClient) RegisterPlatform(platform *types.Platform) (*types.Platform, error) {
	requestBody, err := json.Marshal(platform)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.call(http.MethodPost, "/v1/platforms", buffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 201 {
		return nil, errors.ResponseError{StatusCode: response.StatusCode}
	}

	var newPlatform *types.Platform
	err = httputil.UnmarshalResponse(response, &newPlatform)
	if err != nil {
		return nil, err
	}

	return newPlatform, nil
}

// Registers a broker in the service manager
func (client *ServiceManagerClient) RegisterBroker(broker *types.Broker) (*types.Broker, error) {
	requestBody, err := json.Marshal(broker)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.call(http.MethodPost, "/v1/service_brokers", buffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 201 {
		return nil, errors.ResponseError{StatusCode: response.StatusCode}
	}

	var newBroker *types.Broker
	err = httputil.UnmarshalResponse(response, &newBroker)
	if err != nil {
		return nil, err
	}

	return newBroker, nil
}

func (client *ServiceManagerClient) ListBrokers() (*types.Brokers, error) {
	resp, err := client.call(http.MethodGet, "/v1/service_brokers", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.ResponseError{StatusCode: resp.StatusCode}
	}

	var brokers *types.Brokers = &types.Brokers{}
	err = httputil.UnmarshalResponse(resp, &brokers)
	if err != nil {
		return nil, err
	}

	return brokers, nil
}

func (client *ServiceManagerClient) call(method string, smpath string, body io.Reader) (*http.Response, error) {
	fullURL := httputil.NormalizeURL(client.config.URL)
	fullURL = fullURL + smpath

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header = *client.headers

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// handle errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respErr := errors.ResponseError{
			URL:        fullURL,
			StatusCode: resp.StatusCode,
		}

		respContent := make(map[string]interface{})
		if err := httputil.UnmarshalResponse(resp, &respContent); err != nil {
			return resp, respErr
		}

		if errorMessage, ok := respContent["error"].(string); ok {
			respErr.ErrorMessage = errorMessage
		}

		if description, ok := respContent["description"].(string); ok {
			respErr.Description = description
		}

		return nil, respErr
	}

	return resp, nil
}
