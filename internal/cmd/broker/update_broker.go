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

package broker

import (
	"encoding/json"
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

// UpdateBrokerCmd wraps the smctl update-broker command
type UpdateBrokerCmd struct {
	*cmd.Context

	outputFormat  output.Format
	name          string
	updatedBroker *types.Broker
}

// NewUpdateBrokerCmd returns new update-broker command with context
func NewUpdateBrokerCmd(context *cmd.Context) *UpdateBrokerCmd {
	return &UpdateBrokerCmd{Context: context}
}

// Validate validates command's arguments
func (ubc *UpdateBrokerCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[name] is required")
	}

	ubc.name = args[0]

	if len(args) < 2 {
		return fmt.Errorf("nothing to update. Broker JSON is not provided")
	}

	if err := json.Unmarshal([]byte(args[1]), &ubc.updatedBroker); err != nil {
		return fmt.Errorf("broker JSON is invalid. Reason: %s", err.Error())
	}

	return nil
}

// Run runs the command's logic
func (ubc *UpdateBrokerCmd) Run() error {
	ubc.Parameters.FieldQuery = append(ubc.Parameters.FieldQuery, fmt.Sprintf("name eq '%s'", ubc.name))
	toUpdateBrokers, err := ubc.Client.ListBrokers(&ubc.Parameters)
	if err != nil {
		return err
	}
	if len(toUpdateBrokers.Brokers) < 1 {
		return fmt.Errorf("broker with name %s not found", ubc.name)
	}
	toUpdateBroker := toUpdateBrokers.Brokers[0]
	result, err := ubc.Client.UpdateBroker(toUpdateBroker.ID, ubc.updatedBroker, &ubc.Parameters)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(ubc.Output, ubc.outputFormat, result)
	output.Println(ubc.Output)

	return nil
}

// HideUsage hide command's usage
func (ubc *UpdateBrokerCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (ubc *UpdateBrokerCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "update-broker [name] <json_broker>",
		Aliases: []string{"ub"},
		Short:   "Updates broker",
		Long: `Update broker with name.
Example:
smctl update-broker broker '{"name": "new-name", "description": "new-description", "broker-url": "http://broker.com", "credentials": { "basic": { "username": "admin", "password": "admin" } }}'`,
		PreRunE: prepare(ubc, ubc.Context),
		RunE:    cmd.RunE(ubc),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &ubc.Parameters)

	return result
}

// SetOutputFormat set output format
func (ubc *UpdateBrokerCmd) SetOutputFormat(format output.Format) {
	ubc.outputFormat = format
}
