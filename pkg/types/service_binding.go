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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Peripli/service-manager/pkg/types"
)

// ServiceBinding defines the data of a service instance.
type ServiceBinding struct {
	ID             string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name           string       `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt      string       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt      string       `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Labels         types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
	PagingSequence int64        `json:"-" yaml:"-"`

	Credentials       json.RawMessage `json:"credentials,omitempty" yaml:"credentials,omitempty"`
	ServiceInstanceID string          `json:"service_instance_id" yaml:"service_instance_id,omitempty"`
	SyslogDrainURL    string          `json:"syslog_drain_url,omitempty" yaml:"syslog_drain_url,omitempty"`
	RouteServiceURL   string          `json:"route_service_url,omitempty"`
	VolumeMounts      json.RawMessage `json:"-" yaml:"-"`
	Endpoints         json.RawMessage `json:"-" yaml:"-"`
	Context           json.RawMessage `json:"-" yaml:"-"`
	BindResource      json.RawMessage `json:"-" yaml:"-"`

	Ready bool `json:"ready,omitempty" yaml:"ready,omitempty"`

	LastOperation *types.Operation `json:"last_operation,omitempty" yaml:"last_operation,omitempty"`
}

// Message title of the table
func (sb *ServiceBinding) Message() string {
	return ""
}

// IsEmpty whether the structure is empty
func (sb *ServiceBinding) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (sb *ServiceBinding) TableData() *TableData {
	result := &TableData{Vertical: true}
	result.Headers = []string{"ID", "Name", "Service Instance ID", "Credentials", "Created", "Updated", "Ready", "Labels", "Last Op"}

	lastState := "-"
	if sb.LastOperation != nil {
		lastState = formatLastOp(sb.LastOperation)
	}
	row := []string{sb.ID, sb.Name, sb.ServiceInstanceID, string(sb.Credentials), sb.CreatedAt, sb.UpdatedAt, strconv.FormatBool(sb.Ready), formatLabels(sb.Labels), lastState}
	result.Data = append(result.Data, row)

	return result
}

// ServiceBindings wraps an array of service bindings
type ServiceBindings struct {
	ServiceBindings []ServiceBinding `json:"items" yaml:"items"`
	Vertical        bool             `json:"-" yaml:"-"`
}

// Message title of the table
func (sb *ServiceBindings) Message() string {
	var msg string

	if len(sb.ServiceBindings) == 0 {
		msg = "There are no service bindings."
	} else if len(sb.ServiceBindings) == 1 {
		msg = "One service binding."
	} else {
		msg = fmt.Sprintf("%d service bindings.", len(sb.ServiceBindings))
	}

	return msg
}

// IsEmpty whether the structure is empty
func (sb *ServiceBindings) IsEmpty() bool {
	return len(sb.ServiceBindings) == 0
}

// TableData returns the data to populate a table
func (sb *ServiceBindings) TableData() *TableData {
	result := &TableData{Vertical: sb.Vertical}
	result.Headers = []string{"ID", "Name", "Service Instance ID", "Credentials", "Created", "Updated", "Ready", "Labels"}

	addLastOpColumn := false
	for _, binding := range sb.ServiceBindings {
		lastState := "-"
		if binding.LastOperation != nil {
			lastState = formatLastOp(binding.LastOperation)
			addLastOpColumn = true
		}
		row := []string{binding.ID, binding.Name, binding.ServiceInstanceID, string(binding.Credentials), binding.CreatedAt, binding.UpdatedAt, strconv.FormatBool(binding.Ready), formatLabels(binding.Labels), lastState}
		result.Data = append(result.Data, row)
	}

	if addLastOpColumn {
		result.Headers = append(result.Headers, "Last Op")
	} else {
		for i := range result.Data {
			result.Data[i] = result.Data[i][:len(result.Data[i])-1]
		}
	}

	return result
}
