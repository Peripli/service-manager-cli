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
	"fmt"

	"github.com/Peripli/service-manager-cli/pkg/query"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
)

// GetBrokerCmd wraps the smctl list-brokers command
type GetBrokerCmd struct {
	*cmd.Context

	name         string
	prepare      cmd.PrepareFunc
	outputFormat output.Format
}

// NewGetBrokerCmd returns new get status command with context
func NewGetBrokerCmd(context *cmd.Context) *GetBrokerCmd {
	return &GetBrokerCmd{Context: context}
}

// Run runs the command's logic
func (gb *GetBrokerCmd) Run() error {
	brokers, err := gb.Client.ListBrokers(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", gb.name),
		},
	})
	if err != nil {
		return err
	}
	if len(brokers.Brokers) < 1 {
		output.PrintMessage(gb.Output, "No broker found with name: %s", gb.name)
		return nil
	}

	id := brokers.Brokers[0].ID
	broker, err := gb.Client.GetBrokerByID(id, &query.Parameters{
		GeneralParams: []string{
			"last_op=true",
		},
	})
	if err != nil {
		return err
	}
	output.PrintServiceManagerObject(gb.Output, gb.outputFormat, broker)
	output.Println(gb.Output)

	return nil
}

// Validate validates command's arguments
func (gb *GetBrokerCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("name is required")
	}

	gb.name = args[0]

	return nil
}

// SetOutputFormat set output format
func (gb *GetBrokerCmd) SetOutputFormat(format output.Format) {
	gb.outputFormat = format
}

// HideUsage hide command's usage
func (gb *GetBrokerCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (gb *GetBrokerCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	gb.prepare = prepare
	result := &cobra.Command{
		Use:     "get-broker [name]",
		Aliases: []string{"gb"},
		Short:   "Get single broker",
		Long:    `Get single broker by its name`,
		PreRunE: gb.prepare(gb, gb.Context),
		RunE:    cmd.RunE(gb),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &gb.Parameters)

	return result
}
