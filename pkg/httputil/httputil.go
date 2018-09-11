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

package httputil

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
)

type HTTPConfig struct {
	SSLDisabled bool
	Timeout     time.Duration
	KeepAlive   time.Duration
}

func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		SSLDisabled: false,
		Timeout:     time.Second * 10,
		KeepAlive:   time.Second * 30,
	}
}

// BuildHTTPClient builds custom http client with configured ssl validation
func BuildHTTPClient(config *HTTPConfig) *http.Client {
	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   config.Timeout,
				KeepAlive: config.KeepAlive,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	if config.SSLDisabled {
		client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return client
}

// UnmarshalResponse reads the response body and tries to parse it.
func UnmarshalResponse(response *http.Response, jsonResult interface{}) error {
	defer func() {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	return json.NewDecoder(response.Body).Decode(&jsonResult)
}

// NormalizeURL removes trailing slashesh in url
func NormalizeURL(url string) string {
	for strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return url
}
