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

import (
	"fmt"
)

// Error wraps command errors
type Error struct {
	Description string
	Cause       error
}

// New is used to create command errors
func New(description string, cause error) *Error {
	return &Error{
		description,
		cause,
	}
}

// Error prints the error description and the reason for the error
func (e *Error) Error() string {
	return fmt.Sprintf("%s\nReason: %s", e.Description, e.Cause.Error())
}
