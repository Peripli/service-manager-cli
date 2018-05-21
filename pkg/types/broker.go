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

import "fmt"

// Broker defines the data of a service broker.
type Broker struct {
	ID          string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string       `json:"name" yaml:"name"`
	URL         string       `json:"broker_url" yaml:"broker_url"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Created     string       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Updated     string       `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
}

func (b *Broker) Message() string {
	return ""
}

func (b *Broker) IsEmpty() bool {
	return false
}

func (b *Broker) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "URL", "Description", "Created", "Updated"}

	row := []string{b.ID, b.Name, b.URL, b.Description, b.Created, b.Updated}
	result.Data = append(result.Data, row)

	return result
}

type Brokers struct {
	Brokers []Broker `json:"brokers"`
}

func (b *Brokers) IsEmpty() bool {
	return len(b.Brokers) == 0
}

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

func (b *Brokers) TableData() *TableData {
	result := &TableData{}
	result.Headers = []string{"ID", "Name", "URL", "Description", "Created", "Updated"}

	for _, broker := range b.Brokers {
		row := []string{broker.ID, broker.Name, broker.URL, broker.Description, broker.Created, broker.Updated}
		result.Data = append(result.Data, row)
	}

	return result
}
