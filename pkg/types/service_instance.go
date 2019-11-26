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
	"github.com/Peripli/service-manager/pkg/types"
	"strconv"
)

// ServiceInstance defines the data of a service instance.
type ServiceInstance struct {
	ID        string `json:"id, omitempty" yaml:"id,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`

	Labels         types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
	PagingSequence int64        `json:"-" yaml:"-"`
	ServicePlanID  string       `json:"service_plan_id,omitempty" yaml:"service_plan_id,omitempty"`
	PlatformID     string       `json:"platform_id,omitempty" yaml:"platform_id,omitempty"`

	MaintenanceInfo json.RawMessage `json:"maintenance_info,omitempty" yaml:"-"`
	Context         json.RawMessage `json:"-" yaml:"-"`
	PreviousValues  json.RawMessage `json:"-" yaml: "-"`

	Ready  bool `json:"ready,omitempty" yaml:"ready,omitempty"`
	Usable bool `json:"usable,omitempty" yaml:"usable,omitempty"`
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
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "Service Plan ID", "Platform ID", "Created", "Updated", "Ready", "Usable", "Labels"}

	row := []string{si.ID, si.Name, si.ServicePlanID, si.PlatformID, si.CreatedAt, si.UpdatedAt, strconv.FormatBool(si.Ready), strconv.FormatBool(si.Usable), formatLabels(si.Labels)}
	result.Data = append(result.Data, row)

	return result
}

// ServiceInstances wraps an array of service instances
type ServiceInstances struct {
	ServiceInstances []ServiceInstance `json:"items" yaml:"items"`
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
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "Service Plan ID", "Platform ID", "Created", "Updated", "Ready", "Usable", "Labels"}

	for _, instance := range si.ServiceInstances {
		row := []string{instance.ID, instance.Name, instance.ServicePlanID, instance.PlatformID, instance.CreatedAt, instance.UpdatedAt, strconv.FormatBool(instance.Ready), strconv.FormatBool(instance.Usable), formatLabels(instance.Labels)}
		result.Data = append(result.Data, row)
	}

	return result
}
