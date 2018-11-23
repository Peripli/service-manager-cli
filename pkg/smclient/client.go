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

	"github.com/Peripli/service-manager-cli/pkg/auth/oidc"

	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

// Client should be implemented by SM clients
//go:generate counterfeiter . Client
type Client interface {
	GetInfo() (*types.Info, error)
	RegisterPlatform(*types.Platform) (*types.Platform, error)
	RegisterBroker(*types.Broker) (*types.Broker, error)
	ListBrokers() (*types.Brokers, error)
	ListPlatforms() (*types.Platforms, error)
	DeleteBroker(string) error
	DeletePlatform(string) error
	UpdateBroker(string, *types.Broker) (*types.Broker, error)
	UpdatePlatform(string, *types.Platform) (*types.Platform, error)

	// Call makes HTTP request to the Service Manager server with authentication.
	// It should be used only in case there is no already implemented method for such an operation
	Call(method string, smpath string, body io.Reader) (*http.Response, error)
}

type serviceManagerClient struct {
	config     *ClientConfig
	httpClient auth.Client
}

// NewClientWithAuth returns new SM Client configured with the provided configuration
func NewClientWithAuth(httpClient auth.Client, config *ClientConfig) (Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := &serviceManagerClient{config: config, httpClient: httpClient}
	info, err := client.GetInfo()
	if err != nil {
		return nil, err
	}

	authOptions := &auth.Options{
		IssuerURL:    info.TokenIssuerURL,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		SSLDisabled:  config.SSLDisabled,
	}
	var authStrategy auth.Authenticator
	authStrategy, authOptions, err = oidc.NewOpenIDStrategy(authOptions)

	if err != nil {
		return nil, err
	}

	token, err := auth.GetToken(authOptions, authStrategy)
	if err != nil {
		return nil, err
	}
	authClient := oidc.NewClient(authOptions, token)
	client = &serviceManagerClient{config: config, httpClient: authClient}

	return client, nil
}

// NewClient returns new SM client which will use the http client provided to make calls
func NewClient(httpClient auth.Client, URL string) Client {
	return &serviceManagerClient{config: &ClientConfig{URL: URL}, httpClient: httpClient}
}

func (client *serviceManagerClient) GetInfo() (*types.Info, error) {
	response, err := client.Call(http.MethodGet, "/v1/info", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.ResponseError{StatusCode: response.StatusCode}
	}

	var result *types.Info

	err = httputil.UnmarshalResponse(response, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RegisterPlatform registers a platform in the service manager
func (client *serviceManagerClient) RegisterPlatform(platform *types.Platform) (*types.Platform, error) {
	requestBody, err := json.Marshal(platform)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.Call(http.MethodPost, "/v1/platforms", buffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.ResponseError{StatusCode: response.StatusCode}
	}

	var newPlatform *types.Platform
	err = httputil.UnmarshalResponse(response, &newPlatform)
	if err != nil {
		return nil, err
	}

	return newPlatform, nil
}

// RegisterBroker registers a broker in the service manager
func (client *serviceManagerClient) RegisterBroker(broker *types.Broker) (*types.Broker, error) {
	requestBody, err := json.Marshal(broker)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.Call(http.MethodPost, "/v1/service_brokers", buffer)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.ResponseError{StatusCode: response.StatusCode}
	}

	var newBroker *types.Broker
	err = httputil.UnmarshalResponse(response, &newBroker)
	if err != nil {
		return nil, err
	}

	return newBroker, nil
}

// ListBrokers returns brokers registered in the Service Manager
func (client *serviceManagerClient) ListBrokers() (*types.Brokers, error) {
	brokers := &types.Brokers{}
	err := client.list(brokers, "/v1/service_brokers")

	return brokers, err
}

// ListPlatforms returns platforms registered in the Service Manager
func (client *serviceManagerClient) ListPlatforms() (*types.Platforms, error) {
	platforms := &types.Platforms{}
	err := client.list(platforms, "/v1/platforms")

	return platforms, err
}

func (client *serviceManagerClient) list(result interface{}, path string) error {
	resp, err := client.Call(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.ResponseError{StatusCode: resp.StatusCode}
	}

	return httputil.UnmarshalResponse(resp, &result)
}

func (client *serviceManagerClient) DeleteBroker(id string) error {
	return client.delete(id, "/v1/service_brokers/")
}

func (client *serviceManagerClient) DeletePlatform(id string) error {
	return client.delete(id, "/v1/platforms/")
}

func (client *serviceManagerClient) delete(id, path string) error {
	resp, err := client.Call(http.MethodDelete, path+id, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.ResponseError{StatusCode: resp.StatusCode}
	}

	return nil
}

func (client *serviceManagerClient) UpdateBroker(id string, updatedBroker *types.Broker) (*types.Broker, error) {
	requestBody, err := json.Marshal(updatedBroker)
	if err != nil {
		return nil, err
	}

	result := &types.Broker{}
	return result, client.update(result, requestBody, id, "/v1/service_brokers/")
}

func (client *serviceManagerClient) UpdatePlatform(id string, updatedPlatform *types.Platform) (*types.Platform, error) {
	requestBody, err := json.Marshal(updatedPlatform)
	if err != nil {
		return nil, err
	}
	result := &types.Platform{}
	return result, client.update(result, requestBody, id, "/v1/platforms/")
}

func (client *serviceManagerClient) update(result interface{}, body []byte, id, path string) error {
	buffer := bytes.NewBuffer(body)
	resp, err := client.Call(http.MethodPatch, path+id, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.ResponseError{StatusCode: resp.StatusCode}
	}

	return httputil.UnmarshalResponse(resp, &result)
}

func (client *serviceManagerClient) Call(method string, smpath string, body io.Reader) (*http.Response, error) {
	fullURL := httputil.NormalizeURL(client.config.URL)
	fullURL = fullURL + smpath

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

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
