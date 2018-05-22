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

package errors

import "fmt"

// ResponseError custom error
type ResponseError struct {
	URL          string
	StatusCode   int
	ErrorMessage string
	Description  string
}

// Error implementation of Error interface
func (e ResponseError) Error() string {
	errorMessage := "<nil>"
	description := "<nil>"

	if e.ErrorMessage != "" {
		errorMessage = e.ErrorMessage
	}
	if e.Description != "" {
		description = e.Description
	}

	return fmt.Sprintf("URL: %s, Status: %v; ErrorMessage: %v; Description: %v", e.URL, e.StatusCode, errorMessage, description)
}
