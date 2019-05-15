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

// Broker defines the data of a service broker.
type Broker struct {
	ID          string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string       `json:"name,omitempty" yaml:"name,omitempty"`
	URL         string       `json:"broker_url,omitempty" yaml:"broker_url,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Created     string       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Updated     string       `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
	Labels      types.Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// Message title of the table
func (b *Broker) Message() string {
	return ""
}

// IsEmpty whether the structure is empty
func (b *Broker) IsEmpty() bool {
	return false
}

// TableData returns the data to populate a table
func (b *Broker) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "URL", "Description", "Created", "Updated", "Labels"}

	row := []string{b.ID, b.Name, b.URL, b.Description, b.Created, b.Updated, formatLabels(b.Labels)}
	result.Data = append(result.Data, row)

	return result
}

// Brokers wraps an array of brokers
type Brokers struct {
	Brokers []Broker `json:"service_brokers"`
}

// IsEmpty whether the structure is empty
func (b *Brokers) IsEmpty() bool {
	return len(b.Brokers) == 0
}

// Message title of the table
func (b *Brokers) Message() string {
	var msg string

	if len(b.Brokers) == 0 {
		msg = "No brokers registered."
	} else if len(b.Brokers) == 1 {
		msg = "One broker registered."
	} else {
		msg = fmt.Sprintf("%d brokers registered.", len(b.Brokers))
	}

	return msg
}

// TableData returns the data to populate a table
func (b *Brokers) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "URL", "Description", "Created", "Updated", "Labels"}

	for _, broker := range b.Brokers {
		row := []string{broker.ID, broker.Name, broker.URL, broker.Description, broker.Created, broker.Updated, formatLabels(broker.Labels)}
		result.Data = append(result.Data, row)
	}

	return result
}
