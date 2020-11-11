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
	"fmt"
	"github.com/Peripli/service-manager/pkg/log"
	"io"
	"net/http"

	"github.com/Peripli/service-manager/pkg/util"

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

	RegisterPlatform(*types.Platform, *query.Parameters) (*types.Platform, error)
	ListPlatforms(*query.Parameters) (*types.Platforms, error)
	UpdatePlatform(string, *types.Platform, *query.Parameters) (*types.Platform, error)
	DeletePlatforms(*query.Parameters) error

	RegisterBroker(*types.Broker, *query.Parameters) (*types.Broker, string, error)
	GetBrokerByID(string, *query.Parameters) (*types.Broker, error)
	ListBrokers(*query.Parameters) (*types.Brokers, error)
	UpdateBroker(string, *types.Broker, *query.Parameters) (*types.Broker, string, error)
	DeleteBroker(string, *query.Parameters) (string, error)

	RegisterVisibility(*types.Visibility, *query.Parameters) (*types.Visibility, error)
	ListVisibilities(*query.Parameters) (*types.Visibilities, error)
	UpdateVisibility(string, *types.Visibility, *query.Parameters) (*types.Visibility, error)
	DeleteVisibilities(*query.Parameters) error

	ListOfferings(*query.Parameters) (*types.ServiceOfferings, error)
	ListPlans(*query.Parameters) (*types.ServicePlans, error)

	ListInstances(*query.Parameters) (*types.ServiceInstances, error)
	GetInstanceByID(string, *query.Parameters) (*types.ServiceInstance, error)
	GetInstanceParameters(string, *query.Parameters) (map[string]interface{}, error)
	UpdateInstance(string, *types.ServiceInstance, *query.Parameters) (*types.ServiceInstance, string, error)
	Provision(*types.ServiceInstance, *query.Parameters) (*types.ServiceInstance, string, error)
	Deprovision(string, *query.Parameters) (string, error)

	ListBindings(*query.Parameters) (*types.ServiceBindings, error)
	GetBindingByID(string, *query.Parameters) (*types.ServiceBinding, error)
	GetBindingParameters(string, *query.Parameters) (map[string]interface{}, error)
	Bind(*types.ServiceBinding, *query.Parameters) (*types.ServiceBinding, string, error)
	Unbind(string, *query.Parameters) (string, error)

	Label(string, string, *types.LabelChanges, *query.Parameters) error

	Status(string, *query.Parameters) (*types.Operation, error)

	Marketplace(*query.Parameters) (*types.Marketplace, error)

	// Call makes HTTP request to the Service Manager server with authentication.
	// It should be used only in case there is no already implemented method for such an operation
	Call(method string, smpath string, body io.Reader, q *query.Parameters) (*http.Response, error)
}

type ServiceManagerClient struct {
	Context    context.Context
	Config     *ClientConfig
	HttpClient auth.Client
}

// NewClientWithAuth returns new SM Client configured with the provided configuration
func NewClientWithAuth(httpClient auth.Client, config *ClientConfig) (Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	client := &ServiceManagerClient{Context: context.Background(), Config: config, HttpClient: httpClient}
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
	client = &ServiceManagerClient{Context: context.Background(), Config: config, HttpClient: authClient}

	return client, nil
}

// NewClient returns new SM client which will use the http client provided to make calls
func NewClient(ctx context.Context, httpClient auth.Client, URL string) Client {
	return &ServiceManagerClient{Context: ctx, Config: &ClientConfig{URL: URL}, HttpClient: httpClient}
}

func (client *ServiceManagerClient) GetInfo(q *query.Parameters) (*types.Info, error) {
	response, err := client.Call(http.MethodGet, web.InfoURL, nil, q)
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
func (client *ServiceManagerClient) RegisterPlatform(platform *types.Platform, q *query.Parameters) (*types.Platform, error) {
	var newPlatform *types.Platform
	_, err := client.register(platform, web.PlatformsURL, q, &newPlatform)
	if err != nil {
		return nil, err
	}
	return newPlatform, nil
}

// RegisterBroker registers a broker in the service manager
func (client *ServiceManagerClient) RegisterBroker(broker *types.Broker, q *query.Parameters) (*types.Broker, string, error) {
	var newBroker *types.Broker
	location, err := client.register(broker, web.ServiceBrokersURL, q, &newBroker)
	if err != nil {
		return nil, "", err
	}
	return newBroker, location, nil
}

// RegisterVisibility registers a visibility in the service manager
func (client *ServiceManagerClient) RegisterVisibility(visibility *types.Visibility, q *query.Parameters) (*types.Visibility, error) {
	var newVisibility *types.Visibility
	_, err := client.register(visibility, web.VisibilitiesURL, q, &newVisibility)
	if err != nil {
		return nil, err
	}
	return newVisibility, nil
}

// Provision provisions a new service instance in service manager
func (client *ServiceManagerClient) Provision(instance *types.ServiceInstance, q *query.Parameters) (*types.ServiceInstance, string, error) {
	var newInstance *types.ServiceInstance
	location, err := client.register(instance, web.ServiceInstancesURL, q, &newInstance)
	if err != nil {
		return nil, "", err
	}
	return newInstance, location, nil
}

// Bind creates binding to an instance in service manager
func (client *ServiceManagerClient) Bind(binding *types.ServiceBinding, q *query.Parameters) (*types.ServiceBinding, string, error) {
	var newBinding *types.ServiceBinding
	location, err := client.register(binding, web.ServiceBindingsURL, q, &newBinding)
	if err != nil {
		return nil, "", err
	}
	return newBinding, location, nil
}

func (client *ServiceManagerClient) register(resource interface{}, url string, q *query.Parameters, result interface{}) (string, error) {
	requestBody, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer(requestBody)
	response, err := client.Call(http.MethodPost, url, buffer, q)
	if err != nil {
		return "", err
	}

	switch response.StatusCode {
	case http.StatusCreated:
		return "", httputil.UnmarshalResponse(response, &result)
	case http.StatusAccepted:
		return response.Header.Get("Location"), nil
	default:
		return "", util.HandleResponseError(response)
	}
}

// GetBrokerByID returns broker registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) GetBrokerByID(id string, q *query.Parameters) (*types.Broker, error) {
	broker := &types.Broker{}
	err := client.get(broker, web.ServiceBrokersURL+"/"+id, &query.Parameters{
		GeneralParams: q.GeneralParams,
	})

	return broker, err
}

// ListBrokers returns brokers registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListBrokers(q *query.Parameters) (*types.Brokers, error) {
	brokers := &types.Brokers{}
	err := client.list(&brokers.Brokers, web.ServiceBrokersURL, q)

	return brokers, err
}

// ListPlatforms returns platforms registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListPlatforms(q *query.Parameters) (*types.Platforms, error) {
	platforms := &types.Platforms{}
	err := client.list(&platforms.Platforms, web.PlatformsURL, q)

	return platforms, err
}

// ListOfferings returns offerings registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListOfferings(q *query.Parameters) (*types.ServiceOfferings, error) {
	offerings := &types.ServiceOfferings{}
	err := client.list(&offerings.ServiceOfferings, web.ServiceOfferingsURL, q)

	return offerings, err
}

// ListPlans returns plans registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListPlans(q *query.Parameters) (*types.ServicePlans, error) {
	plans := &types.ServicePlans{}
	err := client.list(&plans.ServicePlans, web.ServicePlansURL, q)

	return plans, err
}

// ListVisibilities returns visibilities registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListVisibilities(q *query.Parameters) (*types.Visibilities, error) {
	visibilities := &types.Visibilities{}
	err := client.list(&visibilities.Visibilities, web.VisibilitiesURL, q)

	return visibilities, err
}

// ListInstances returns service instances registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListInstances(q *query.Parameters) (*types.ServiceInstances, error) {
	instances := &types.ServiceInstances{}
	err := client.list(&instances.ServiceInstances, web.ServiceInstancesURL, q)

	return instances, err
}
// GetInstanceParameters returns service instance configuration parameters
func (client *ServiceManagerClient) GetInstanceParameters(id string, q *query.Parameters) (map[string]interface{}, error) {
	parameters := make(map[string]interface{})
	err := client.get(&parameters, web.ServiceInstancesURL + "/" + id + web.ParametersURL, q)

	return parameters, err
}

// GetInstanceByID returns instance registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) GetInstanceByID(id string, q *query.Parameters) (*types.ServiceInstance, error) {
	instance := &types.ServiceInstance{}
	err := client.get(instance, web.ServiceInstancesURL+"/"+id, &query.Parameters{
		GeneralParams: q.GeneralParams,
	})

	return instance, err
}

// ListBindings returns service bindings registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) ListBindings(q *query.Parameters) (*types.ServiceBindings, error) {
	bindings := &types.ServiceBindings{}
	err := client.list(&bindings.ServiceBindings, web.ServiceBindingsURL, q)

	return bindings, err
}

// GetBindingParameters returns service binding configuration parameters
func (client *ServiceManagerClient) GetBindingParameters(id string, q *query.Parameters) (map[string]interface{}, error) {
	parameters := make(map[string]interface{})
	err := client.get(&parameters, web.ServiceBindingsURL + "/" + id + web.ParametersURL, q)

	return parameters, err
}

// GetBindingByID returns binding registered in the Service Manager satisfying provided queries
func (client *ServiceManagerClient) GetBindingByID(id string, q *query.Parameters) (*types.ServiceBinding, error) {
	binding := &types.ServiceBinding{}
	err := client.get(binding, web.ServiceBindingsURL+"/"+id, &query.Parameters{
		GeneralParams: q.GeneralParams,
	})

	return binding, err
}

// Marketplace returns service offerings satisfying provided queries
func (client *ServiceManagerClient) Marketplace(q *query.Parameters) (*types.Marketplace, error) {
	marketplace := &types.Marketplace{}
	err := client.list(&marketplace.ServiceOfferings, web.ServiceOfferingsURL, q)
	if err != nil {
		return nil, err
	}
	for i, so := range marketplace.ServiceOfferings {
		plans := &types.ServicePlansForOffering{}
		err := client.list(&plans.ServicePlans, web.ServicePlansURL, &query.Parameters{
			FieldQuery:    []string{fmt.Sprintf("service_offering_id eq '%s'", so.ID)},
			GeneralParams: q.GeneralParams,
		})
		if err != nil {
			return nil, err
		}
		marketplace.ServiceOfferings[i].Plans = plans.ServicePlans
	}
	return marketplace, nil
}

func (client *ServiceManagerClient) Status(url string, q *query.Parameters) (*types.Operation, error) {
	operation := &types.Operation{}
	err := client.get(operation, url, &query.Parameters{
		GeneralParams: q.GeneralParams,
	})

	return operation, err
}

func (client *ServiceManagerClient) list(result interface{}, url string, q *query.Parameters) error {
	fullURL := httputil.NormalizeURL(client.Config.URL) + BuildURL(url, q)
	return util.ListAll(client.Context, client.HttpClient.Do, fullURL, result)
}

func (client *ServiceManagerClient) get(result interface{}, url string, q *query.Parameters) error {
	resp, err := client.Call(http.MethodGet, url, nil, q)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return util.HandleResponseError(resp)
	}

	return httputil.UnmarshalResponse(resp, &result)
}

func (client *ServiceManagerClient) DeleteBroker(id string, q *query.Parameters) (string, error) {
	return client.delete(web.ServiceBrokersURL+"/"+id, q)
}

func (client *ServiceManagerClient) DeletePlatforms(q *query.Parameters) error {
	_, err := client.delete(web.PlatformsURL, q)
	return err
}

func (client *ServiceManagerClient) DeleteVisibilities(q *query.Parameters) error {
	_, err := client.delete(web.VisibilitiesURL, q)
	return err
}

func (client *ServiceManagerClient) Deprovision(id string, q *query.Parameters) (string, error) {
	return client.delete(web.ServiceInstancesURL+"/"+id, q)
}

func (client *ServiceManagerClient) Unbind(id string, q *query.Parameters) (string, error) {
	return client.delete(web.ServiceBindingsURL+"/"+id, q)
}

func (client *ServiceManagerClient) delete(url string, q *query.Parameters) (string, error) {
	resp, err := client.Call(http.MethodDelete, url, nil, q)
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return "", nil
	case http.StatusAccepted:
		return resp.Header.Get("Location"), nil
	default:
		return "", util.HandleResponseError(resp)
	}
}

func (client *ServiceManagerClient) UpdateBroker(id string, updatedBroker *types.Broker, q *query.Parameters) (*types.Broker, string, error) {
	var result *types.Broker
	location, err := client.update(updatedBroker, web.ServiceBrokersURL, id, q, &result)
	if err != nil {
		return nil, "", err
	}
	return result, location, nil
}

func (client *ServiceManagerClient) UpdateInstance(id string, updatedInstance *types.ServiceInstance, q *query.Parameters) (*types.ServiceInstance, string, error) {
	var result *types.ServiceInstance
	location, err := client.update(updatedInstance, web.ServiceInstancesURL, id, q, &result)
	if err != nil {
		return nil, "", err
	}
	return result, location, nil
}

func (client *ServiceManagerClient) UpdatePlatform(id string, updatedPlatform *types.Platform, q *query.Parameters) (*types.Platform, error) {
	result := &types.Platform{}
	_, err := client.update(updatedPlatform, web.PlatformsURL, id, q, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *ServiceManagerClient) UpdateVisibility(id string, updatedVisibility *types.Visibility, q *query.Parameters) (*types.Visibility, error) {
	result := &types.Visibility{}
	_, err := client.update(updatedVisibility, web.VisibilitiesURL, id, q, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *ServiceManagerClient) update(resource interface{}, url string, id string, q *query.Parameters, result interface{}) (string, error) {
	requestBody, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBuffer(requestBody)
	resp, err := client.Call(http.MethodPatch, url+"/"+id, buffer, q)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return "", httputil.UnmarshalResponse(resp, &result)
	case http.StatusAccepted:
		return resp.Header.Get("Location"), nil
	default:
		return "", util.HandleResponseError(resp)
	}
}

func (client *ServiceManagerClient) Label(url string, id string, change *types.LabelChanges, q *query.Parameters) error {
	requestBody, err := json.Marshal(change)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(requestBody)
	response, err := client.Call(http.MethodPatch, url+"/"+id, buffer, q)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return util.HandleResponseError(response)
	}

	return nil
}

func (client *ServiceManagerClient) Call(method string, smpath string, body io.Reader, q *query.Parameters) (*http.Response, error) {
	fullURL := httputil.NormalizeURL(client.Config.URL) + BuildURL(smpath, q)

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	log.C(client.Context).Debugf("Sending request %s %s", req.Method, req.URL)
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, util.HandleResponseError(resp)
	}

	return resp, nil
}

// BuildURL builds the url with provided query parameters
func BuildURL(baseURL string, q *query.Parameters) string {
	queryParams := q.Encode()
	if queryParams == "" {
		return baseURL
	}
	return baseURL + "?" + queryParams
}
