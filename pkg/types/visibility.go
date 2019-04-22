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

// Visibility defines the data of a visibility
type Visibility struct {
	ID            string       `json:"id,omitempty" yaml:"id,omitempty"`
	PlatformID    string       `json:"platform_id,omitempty" yaml:"platform_id,omitempty"`
	ServicePlanID string       `json:"service_plan_id,omitempty" yaml:"service_plan_id,omitempty"`
	CreatedAt     string       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt     string       `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Labels        types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Message title of the table
func (v *Visibility) Message() string {
	return ""
}

// IsEmpty whether the struct is empty
func (v *Visibility) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (v *Visibility) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Platform ID", "Service Plan ID", "Labels"}
	labels := fmt.Sprintf("%v", v.Labels)

	row := []string{v.ID, v.PlatformID, v.ServicePlanID, labels}
	result.Data = append(result.Data, row)

	return result
}

// Visibilities wraps an array of Visibilities
type Visibilities struct {
	Visibilities []Visibility `json:"visibilities,omitempty" yaml:"visibilities,omitempty"`
}

// IsEmpty whether the structure is empty
func (v *Visibilities) IsEmpty() bool {
	return len(v.Visibilities) == 0
}

// Message title of the table
func (v *Visibilities) Message() string {
	var msg string

	if len(v.Visibilities) == 0 {
		msg = "No visibilities registered."
	} else if len(v.Visibilities) == 1 {
		msg = "One visibility registered."
	} else {
		msg = fmt.Sprintf("%d visibilities registered.", len(v.Visibilities))
	}

	return msg
}

// TableData returns the data to populate a table
func (v *Visibilities) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Platform ID", "Service Plan ID", "Labels"}

	for _, visibility := range v.Visibilities {
		labels := fmt.Sprintf("%v", visibility.Labels)
		row := []string{visibility.ID, visibility.PlatformID, visibility.ServicePlanID, labels}
		result.Data = append(result.Data, row)
	}

	return result
}
