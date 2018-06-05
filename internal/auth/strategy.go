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

package auth

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	resterror "github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
)

// Token contains the structure of a typical UAA response token
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	JTI          string `json:"jti"`
}

type openIDConfiguration struct {
	TokenEndpoint string `json:"token_endpoint"`
}

// AuthenticationStrategy should be implemented for different authentication strategies
//go:generate counterfeiter . AuthenticationStrategy
type AuthenticationStrategy interface {
	Authenticate(issuerURL, user, password string) (*Token, error)
}

// NewOpenIDStrategy returns OpenId auth strategy
func NewOpenIDStrategy() AuthenticationStrategy {
	return &OpenIDStrategy{}
}

// OpenIDStrategy implementation of OpenID strategy
type OpenIDStrategy struct{}

// Authenticate is used to perform authentication action for OpenID strategy
func (s *OpenIDStrategy) Authenticate(issuerURL, user, password string) (*Token, error) {
	tokenEndpoint, err := s.getTokenEndpoint(issuerURL)
	if err != nil {
		return nil, err
	}

	token, err := s.getAccessToken(tokenEndpoint, user, password)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *OpenIDStrategy) getAccessToken(tokenEndpoint, user, password string) (*Token, error) {
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("response_type", "token")
	form.Add("username", user)
	form.Add("password", password)

	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("cf", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpClient := getInsecureClient()
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		respContent := make(map[string]interface{})
		if err := httputil.UnmarshalResponse(response, &respContent); err != nil {
			return nil, err
		}
		responseString, err := json.Marshal(respContent)
		if err != nil {
			return nil, err
		}

		return nil, resterror.ResponseError{StatusCode: response.StatusCode, URL: req.URL.String(), Description: fmt.Sprintf("Could not get access token.\nReason:\n%s", responseString)}
	}

	var token *Token
	return token, httputil.UnmarshalResponse(response, &token)
}

func (s *OpenIDStrategy) getTokenEndpoint(issuerURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, issuerURL+"/.well-known/openid-configuration", nil)
	if err != nil {
		return "", err
	}

	httpClient := getInsecureClient()
	response, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New("Error getting OpenID configuration")
	}

	var configuration *openIDConfiguration
	if err = httputil.UnmarshalResponse(response, &configuration); err != nil {
		return "", err
	}

	return configuration.TokenEndpoint, nil
}

func getInsecureClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}
