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
	"encoding/json"
	"net/http"
	"strings"
)

// UnmarshalResponse reads the response body and tries to parse it.
func UnmarshalResponse(response *http.Response, jsonResult interface{}) error {

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&jsonResult); err != nil {
		return err
	}

	return nil
}

func NormalizeURL(url string) string {
	for strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return url
}
