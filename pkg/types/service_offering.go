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
	"strings"

	"github.com/Peripli/service-manager/pkg/types"
)

// ServiceOffering defines the data of a service offering.
type ServiceOffering struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`

	Bindable             bool   `json:"bindable,omitempty" yaml:"bindable,omitempty"`
	InstancesRetrievable bool   `json:"instances_retrievable,omitempty" yaml:"instances_retrievable,omitempty"`
	BindingsRetrievable  bool   `json:"bindings_retrievable,omitempty" yaml:"bindings_retrievable,omitempty"`
	AllowContextUpdates  bool   `json:"allow_context_updates,omitempty" yaml:"allow_context_updates,omitempty"`
	PlanUpdatable        bool   `json:"plan_updateable,omitempty" yaml:"plan_updateable,omitempty"`
	CatalogID            string `json:"catalog_id,omitempty" yaml:"catalog_id,omitempty"`
	CatalogName          string `json:"catalog_name,omitempty" yaml:"catalog_name,omitempty"`

	Tags     json.RawMessage `json:"tags,omitempty" yaml:"-"`
	Requires json.RawMessage `json:"requires,omitempty" yaml:"-"`
	Metadata json.RawMessage `json:"metadata,omitempty" yaml:"-"`

	BrokerID   string        `json:"broker_id,omitempty" yaml:"broker_id,omitempty"`
	BrokerName string        `json:"broker_name,omitempty" yaml:"broker_name,omitempty"`
	Plans      []ServicePlan `json:"plans,omitempty" yaml:"plans,omitempty"`
	Labels     types.Labels  `json:"labels,omitempty" yaml:"labels,omitempty"`
	Ready      bool          `json:"ready,omitempty" yaml:"ready,omitempty"`
}

const BROKER_ID = "BROKER ID"

// Message title of the table
func (so *ServiceOffering) Message() string {
	return ""
}

// IsEmpty whether the structure is empty
func (so *ServiceOffering) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (so *ServiceOffering) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"Name", "Plans", "Description", BROKER_ID}

	plans := make([]string, len(so.Plans))
	for i, v := range so.Plans {
		plans[i] = v.Name
	}

	row := []string{so.Name, strings.Join(plans, ", "), so.Description, so.BrokerID}
	result.Data = append(result.Data, row)

	return result
}

// Marketplace wraps an array of service offerings with marketplace related table data
type Marketplace struct {
	ServiceOfferings []ServiceOffering `json:"items" yaml:"items"`
}

// Message title of the table
func (m *Marketplace) Message() string {
	var msg string

	if len(m.ServiceOfferings) == 0 {
		msg = "There are no service offerings."
	} else if len(m.ServiceOfferings) == 1 {
		msg = "One service offering."
	} else {
		msg = fmt.Sprintf("%d service offerings.", len(m.ServiceOfferings))
	}

	return msg
}

// IsEmpty whether the structure is empty
func (m *Marketplace) IsEmpty() bool {
	return len(m.ServiceOfferings) == 0
}

// TableData returns the data to populate a table
func (m *Marketplace) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"Name", "Plans", "Description", BROKER_ID}

	for _, v := range m.ServiceOfferings {
		plans := make([]string, len(v.Plans))
		for i, v := range v.Plans {
			plans[i] = v.Name
		}

		row := []string{v.Name, strings.Join(plans, ", "), v.Description, v.BrokerID}
		result.Data = append(result.Data, row)
	}

	return result
}

// ServiceOfferings wraps an array of service offerings
type ServiceOfferings struct {
	ServiceOfferings []ServiceOffering `json:"items" yaml:"items"`
}

// Message title of the table
func (so *ServiceOfferings) Message() string {
	var msg string

	if len(so.ServiceOfferings) == 0 {
		msg = "There are no service offerings."
	} else if len(so.ServiceOfferings) == 1 {
		msg = "One service offering."
	} else {
		msg = fmt.Sprintf("%d service offerings.", len(so.ServiceOfferings))
	}

	return msg
}

// IsEmpty whether the structure is empty
func (so *ServiceOfferings) IsEmpty() bool {
	return len(so.ServiceOfferings) == 0
}

// TableData returns the data to populate a table
func (so *ServiceOfferings) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "Description", BROKER_ID, "Ready", "Labels"}

	for _, v := range so.ServiceOfferings {
		row := []string{v.ID, v.Name, v.Description, v.BrokerID, strconv.FormatBool(v.Ready), formatLabels(v.Labels)}
		result.Data = append(result.Data, row)
	}

	return result
}
