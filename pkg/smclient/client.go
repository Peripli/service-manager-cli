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
	"context"
	"encoding/json"
	"github.com/Peripli/service-manager/pkg/util"
	"fmt"
	"io"
	"net/http"

	"github.com/Peripli/service-manager/pkg/web"

	"github.com/Peripli/service-manager-cli/pkg/auth/oidc"

	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

// Client should be implemented by SM clients
//go:generate counterfeiter . Client
type Client interface {
	GetInfo(*query.Parameters) (*types.Info, error)
	RegisterPlatform(*types.Platform) (*types.Platform, error)
	RegisterBroker(*types.Broker) (*types.Broker, error)
	RegisterVisibility(*types.Visibility) (*types.Visibility, error)
	ListBrokers(*query.Parameters) (*types.Brokers, error)
	ListPlatforms(*query.Parameters) (*types.Platforms, error)
	ListOfferings(*query.Parameters) (*types.ServiceOfferings, error)
	ListPlans(*query.Parameters) (*types.ServicePlans, error)
	ListVisibilities(*query.Parameters) (*types.Visibilities, error)
	DeleteBroker(string) error
	DeleteBrokersByFieldQuery(string) error
	DeletePlatform(string) error
	DeleteVisibility(string) error
	DeletePlatformsByFieldQuery(string) error
	UpdateBroker(string, *types.Broker) (*types.Broker, error)
	UpdatePlatform(string, *types.Platform) (*types.Platform, error)
	UpdateVisibility(string, *types.Visibility) (*types.Visibility, error)
	Label(string, string, *types.LabelChanges) error
	Marketplace(*query.Parameters) (*types.Marketplace, error)

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
	info, err := client.GetInfo(nil)
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

func (client *serviceManagerClient) GetInfo(q *query.Parameters) (*types.Info, error) {
	response, err := client.Call(http.MethodGet, buildURL(web.InfoURL, q), nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, util.HandleResponseError(response)
	}

	info := types.DefaultInfo
	err = httputil.UnmarshalResponse(response, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// RegisterPlatform registers a platform in the service manager
func (client *serviceManagerClient) RegisterPlatform(platform *types.Platform) (*types.Platform, error) {
	var newPlatform *types.Platform
	err := client.register(platform, web.PlatformsURL, &newPlatform)
	if err != nil {
		return nil, err
	}
	return newPlatform, nil
}

// RegisterBroker registers a broker in the service manager
func (client *serviceManagerClient) RegisterBroker(broker *types.Broker) (*types.Broker, error) {
	var newBroker *types.Broker
	err := client.register(broker, web.ServiceBrokersURL, &newBroker)
	if err != nil {
		return nil, err
	}
	return newBroker, nil
}

// RegisterVisibility registers a visibility in the service manager
func (client *serviceManagerClient) RegisterVisibility(visibility *types.Visibility) (*types.Visibility, error) {
	var newVisibility *types.Visibility
	err := client.register(visibility, web.VisibilitiesURL, &newVisibility)
	if err != nil {
		return nil, err
	}
	return newVisibility, nil
}

func (client *serviceManagerClient) register(resource interface{}, url string, result interface{}) error {
	requestBody, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.Call(http.MethodPost, url, buffer)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return util.HandleResponseError(response)
	}

	return httputil.UnmarshalResponse(response, &result)
}

// ListBrokers returns brokers registered in the Service Manager satisfying provided queries
func (client *serviceManagerClient) ListBrokers(q *query.Parameters) (*types.Brokers, error) {
	brokers := &types.Brokers{}
	err := client.list(&brokers.Brokers, buildURL(web.ServiceBrokersURL, q))

	return brokers, err
}

// ListPlatforms returns platforms registered in the Service Manager satisfying provided queries
func (client *serviceManagerClient) ListPlatforms(q *query.Parameters) (*types.Platforms, error) {
	platforms := &types.Platforms{}
	err := client.list(&platforms.Platforms, buildURL(web.PlatformsURL, q))

	return platforms, err
}

// ListOfferings returns offerings registered in the Service Manager satisfying provided queries
func (client *serviceManagerClient) ListOfferings(q *query.Parameters) (*types.ServiceOfferings, error) {
	offerings := &types.ServiceOfferings{}
	err := client.list(&offerings.ServiceOfferings, buildURL(web.ServiceOfferingsURL, q))

	return offerings, err
}

// ListPlans returns plans registered in the Service Manager satisfying provided queries
func (client *serviceManagerClient) ListPlans(q *query.Parameters) (*types.ServicePlans, error) {
	plans := &types.ServicePlans{}
	err := client.list(&plans.ServicePlans, buildURL(web.ServicePlansURL, q))

	return plans, err
}

func (client *serviceManagerClient) ListVisibilities(q *query.Parameters) (*types.Visibilities, error) {
	visibilities := &types.Visibilities{}
	err := client.list(&visibilities.Visibilities, buildURL(web.VisibilitiesURL, q))
	return visibilities, err
}

// Marketplace returns service offerings satisfying provided queries
func (client *serviceManagerClient) Marketplace(q *query.Parameters) (*types.Marketplace, error) {
	marketplace := &types.Marketplace{}
	err := client.list(&marketplace.ServiceOfferings, buildURL(web.ServiceOfferingsURL, q))
	if err != nil {
		return nil, err
	}
	for i, so := range marketplace.ServiceOfferings {
		plans := &types.ServicePlansForOffering{}
		err := client.list(&plans.ServicePlans, web.ServicePlansURL+fmt.Sprintf("?fieldQuery=service_offering_id+eq+'%s'",so.ID))
		if err != nil {
			return nil, err
		}
		marketplace.ServiceOfferings[i].Plans = plans.ServicePlans

		broker := &types.Broker{}
		err = client.get(broker, web.ServiceBrokersURL+"/"+so.BrokerID)
		if err != nil {
			return nil, err
		}

		marketplace.ServiceOfferings[i].BrokerName = broker.Name
	}
	return marketplace, nil
}

func (client *serviceManagerClient) list(result interface{}, path string) error {
	fullURL := httputil.NormalizeURL(client.config.URL) + path
	return util.ListAll(context.Background(), client.httpClient.Do, fullURL, result)
}

func (client *serviceManagerClient) get(result interface{}, path string) error {
	resp, err := client.Call(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return util.HandleResponseError(resp)
	}

	return httputil.UnmarshalResponse(resp, &result)
}

func (client *serviceManagerClient) DeleteBrokersByFieldQuery(query string) error {
	return client.delete(web.ServiceBrokersURL + "?fieldQuery=" + query)
}

// DeleteBroker deletes a broker with given id from service manager
func (client *serviceManagerClient) DeleteBroker(id string) error {
	return client.delete(web.ServiceBrokersURL + "/" + id)
}

func (client *serviceManagerClient) DeletePlatformsByFieldQuery(query string) error {
	return client.delete(web.PlatformsURL + "?fieldQuery=" + query)
}

// DeletePlatform deletes a platform with given id from service manager
func (client *serviceManagerClient) DeletePlatform(id string) error {
	return client.delete(web.PlatformsURL + "/" + id)
}

// DeleteVisibility deletes a visibility with given id from service manager
func (client *serviceManagerClient) DeleteVisibility(id string) error {
	return client.delete(web.VisibilitiesURL + "/" + id)
}

func (client *serviceManagerClient) delete(path string) error {
	resp, err := client.Call(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return util.HandleResponseError(resp)
	}

	return nil
}

func (client *serviceManagerClient) UpdateBroker(id string, updatedBroker *types.Broker) (*types.Broker, error) {
	result := &types.Broker{}
	err := client.update(updatedBroker, web.ServiceBrokersURL, id, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *serviceManagerClient) UpdatePlatform(id string, updatedPlatform *types.Platform) (*types.Platform, error) {
	result := &types.Platform{}
	err := client.update(updatedPlatform, web.PlatformsURL, id, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *serviceManagerClient) UpdateVisibility(id string, updatedVisibility *types.Visibility) (*types.Visibility, error) {
	result := &types.Visibility{}
	err := client.update(updatedVisibility, web.VisibilitiesURL, id, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *serviceManagerClient) update(resource interface{}, url string, id string, result interface{}) error {
	requestBody, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(requestBody)
	resp, err := client.Call(http.MethodPatch, url+"/"+id, buffer)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return util.HandleResponseError(resp)
	}

	return httputil.UnmarshalResponse(resp, &result)
}

func (client *serviceManagerClient) Label(resourcePath string, id string, change *types.LabelChanges) error {
	requestBody, err := json.Marshal(change)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(requestBody)
	response, err := client.Call(http.MethodPatch, resourcePath+"/"+id, buffer)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return util.HandleResponseError(response)
	}

	return nil
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
		return nil, util.HandleResponseError(resp)
	}

	return resp, nil
}

func buildURL(baseURL string, q *query.Parameters) string {
	queryParams := q.Encode()
	if queryParams == "" {
		return baseURL
	}
	return baseURL + "?" + queryParams
}
