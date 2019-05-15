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

package types

import (
	"fmt"
	"github.com/Peripli/service-manager/pkg/types"
)

// Platform defines the data of a platform.
type Platform struct {
	ID          string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string       `json:"type" yaml:"type"`
	Created     string       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Updated     string       `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
	Labels      types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Message title of the table
func (p *Platform) Message() string {
	return ""
}

// IsEmpty whether the structure is empty
func (p *Platform) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (p *Platform) TableData() *TableData {
	result := &TableData{}

	result.Headers = []string{"ID", "Name", "Type", "Description", "Created", "Updated", "Labels"}
	row := []string{p.ID, p.Name, p.Type, p.Description, p.Created, p.Updated, formatLabels(p.Labels)}

	if p.Credentials != nil {
		result.Headers = append(result.Headers, "Username", "Password")
		row = append(row, p.Credentials.Basic.User, p.Credentials.Basic.Password)
	}

	result.Data = append(result.Data, row)

	return result
}

// Platforms wraps an array of platforms
type Platforms struct {
	Platforms []Platform `json:"platforms"`
}

// IsEmpty whether the structure is empty
func (p *Platforms) IsEmpty() bool {
	return len(p.Platforms) == 0
}

// Message title of the table
func (p *Platforms) Message() string {
	var msg string

	if len(p.Platforms) == 0 {
		msg = "No platforms registered."
	} else if len(p.Platforms) == 1 {
		msg = "One platform registered."
	} else {
		msg = fmt.Sprintf("%d platforms registered.", len(p.Platforms))
	}

	return msg
}

// TableData returns the data to populate a table
func (p *Platforms) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "Type", "Description", "Created", "Updated", "Labels"}

	for _, platform := range p.Platforms {
		row := []string{platform.ID, platform.Name, platform.Type, platform.Description, platform.Created, platform.Updated, formatLabels(platform.Labels)}
		result.Data = append(result.Data, row)
	}

	return result
}
