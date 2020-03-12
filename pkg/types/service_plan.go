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

// ServicePlan defines the data of a service plan.
type ServicePlan struct {
	ID          string `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`

	CatalogID     string `json:"catalog_id,omitempty" yaml:"catalog_id,omitempty"`
	CatalogName   string `json:"catalog_name,omitempty" yaml:"catalog_name,omitempty"`
	Free          bool   `json:"free,omitempty" yaml:"free,omitempty"`
	Bindable      bool   `json:"bindable,omitempty" yaml:"bindable,omitempty"`
	PlanUpdatable bool   `json:"plan_updateable,omitempty" yaml:"plan_updateable,omitempty"`

	Metadata json.RawMessage `json:"metadata,omitempty" yaml:"-"`
	Schemas  json.RawMessage `json:"schemas,omitempty" yaml:"-"`

	ServiceOfferingID string       `json:"service_offering_id,omitempty" yaml:"service_offering_id,omitempty"`
	Labels            types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
	Ready             bool         `json:"ready,omitempty" yaml:"ready,omitempty"`
}

// Message title of the table
func (sp *ServicePlan) Message() string {
	return ""
}

// IsEmpty whether the structure is empty
func (sp *ServicePlan) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (sp *ServicePlan) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"Plan", "Description", "ID"}

	row := []string{sp.Name, sp.Description, sp.ID}
	result.Data = append(result.Data, row)

	return result
}

// ServicePlansForOffering wraps an array of service plans for marketplace command
type ServicePlansForOffering struct {
	ServicePlans []ServicePlan `json:"items" yaml:"items"`
}

// Message title of the table
func (sp *ServicePlansForOffering) Message() string {
	var msg string

	if len(sp.ServicePlans) == 0 {
		msg = "There are no service plans for this service offering."
	} else if len(sp.ServicePlans) == 1 {
		msg = "One service plan for this service offering."
	} else {
		msg = fmt.Sprintf("%d service plans for this service offering.", len(sp.ServicePlans))
	}

	return msg
}

// IsEmpty whether the structure is empty
func (sp *ServicePlansForOffering) IsEmpty() bool {
	return len(sp.ServicePlans) == 0
}

// TableData returns the data to populate a table
func (sp *ServicePlansForOffering) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"Plan", "Description", "ID"}

	for _, v := range sp.ServicePlans {
		row := []string{v.Name, v.Description, v.ID}
		result.Data = append(result.Data, row)
	}

	return result
}

// ServicePlans wraps an array of service plans
type ServicePlans struct {
	ServicePlans []ServicePlan `json:"items" yaml:"items"`
}

// Message title of the table
func (sp *ServicePlans) Message() string {
	var msg string

	if len(sp.ServicePlans) == 0 {
		msg = "There are no service plans."
	} else if len(sp.ServicePlans) == 1 {
		msg = "One service plan."
	} else {
		msg = fmt.Sprintf("%d service plans.", len(sp.ServicePlans))
	}

	return msg
}

// IsEmpty whether the structure is empty
func (sp *ServicePlans) IsEmpty() bool {
	return len(sp.ServicePlans) == 0
}

// TableData returns the data to populate a table
func (sp *ServicePlans) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "Description", "Offering ID", "Ready"}

	for _, v := range sp.ServicePlans {
		row := []string{v.ID, v.Name, v.Description, v.ServiceOfferingID, strconv.FormatBool(v.Ready)}
		result.Data = append(result.Data, row)
	}

	return result
}
