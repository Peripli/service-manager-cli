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

// ServiceInstance defines the data of a service instance.
type ServiceInstance struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`

	Labels types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`

	ServiceID     string `json:"service_id,omitempty" yaml:"service_id,omitempty"`
	ServicePlanID string `json:"service_plan_id,omitempty" yaml:"service_plan_id,omitempty"`
	PlatformID    string `json:"platform_id,omitempty" yaml:"platform_id,omitempty"`

	Parameters      json.RawMessage `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	MaintenanceInfo json.RawMessage `json:"maintenance_info,omitempty" yaml:"-"`
	Context         json.RawMessage `json:"context,omitempty" yaml:"context,omitempty"`
	PreviousValues  json.RawMessage `json:"-" yaml:"-"`

	Ready         bool             `json:"ready,omitempty" yaml:"ready,omitempty"`
	Usable        bool             `json:"usable,omitempty" yaml:"usable,omitempty"`
	Shared        bool             `json:"shared" yaml:"shared"`
	LastOperation *types.Operation `json:"last_operation,omitempty" yaml:"last_operation,omitempty"`
}

// Message title of the table
func (si *ServiceInstance) Message() string {
	return ""
}

// IsEmpty whether the structure is empty
func (si *ServiceInstance) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (si *ServiceInstance) TableData() *TableData {
	result := &TableData{Vertical: true}
	result.Headers = []string{"ID", "Name", "Service Plan ID", "Platform ID", "Shared", "Created", "Updated", "Ready", "Usable", "Labels", "Last Op"}

	lastState := "-"
	if si.LastOperation != nil {
		lastState = formatLastOp(si.LastOperation)
	}
	row := []string{si.ID, si.Name, si.ServicePlanID, si.PlatformID, strconv.FormatBool(si.Shared), si.CreatedAt, si.UpdatedAt, strconv.FormatBool(si.Ready), strconv.FormatBool(si.Usable), formatLabels(si.Labels), lastState}
	result.Data = append(result.Data, row)

	return result
}

// ServiceInstances wraps an array of service instances
type ServiceInstances struct {
	ServiceInstances []ServiceInstance `json:"items" yaml:"items"`
	Vertical         bool              `json:"-" yaml:"-"`
}

// Message title of the table
func (si *ServiceInstances) Message() string {
	var msg string

	if len(si.ServiceInstances) == 0 {
		msg = "There are no service instances."
	} else if len(si.ServiceInstances) == 1 {
		msg = "One service instance."
	} else {
		msg = fmt.Sprintf("%d service instances.", len(si.ServiceInstances))
	}

	return msg
}

// IsEmpty whether the structure is empty
func (si *ServiceInstances) IsEmpty() bool {
	return len(si.ServiceInstances) == 0
}

// TableData returns the data to populate a table
func (si *ServiceInstances) TableData() *TableData {
	result := &TableData{Vertical: si.Vertical}
	result.Headers = []string{"ID", "Name", "Service Plan ID", "Platform ID", "Shared", "Created", "Updated", "Ready", "Usable", "Labels"}

	addLastOpColumn := false
	for _, instance := range si.ServiceInstances {
		lastState := "-"
		if instance.LastOperation != nil {
			lastState = formatLastOp(instance.LastOperation)
			addLastOpColumn = true
		}
		row := []string{instance.ID, instance.Name, instance.ServicePlanID, instance.PlatformID, strconv.FormatBool(instance.Shared), instance.CreatedAt, instance.UpdatedAt, strconv.FormatBool(instance.Ready), strconv.FormatBool(instance.Usable), formatLabels(instance.Labels), lastState}
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
