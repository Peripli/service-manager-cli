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

package binding

import (
	"fmt"
	"strings"

	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/spf13/cobra"
)

// GetBindingCmd wraps the smctl get-binding command
type GetBindingCmd struct {
	*cmd.Context

	bindingName  string
	outputFormat output.Format
	bindingParams bool
}

// NewGetBindingCmd returns new get status command with context
func NewGetBindingCmd(context *cmd.Context) *GetBindingCmd {
	return &GetBindingCmd{Context: context}
}

// Run runs the command's logic
func (gb *GetBindingCmd) Run() error {
	bindings, err := gb.Client.ListBindings(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", gb.bindingName),
		},
	})
	if err != nil {
		return err
	}
	if len(bindings.ServiceBindings) < 1 {
		output.PrintMessage(gb.Output, "No binding found with name: %s", gb.bindingName)
		return nil
	}
	if gb.bindingParams {
		return gb.printParameters(bindings)
	}

	resultBindings := &types.ServiceBindings{Vertical: true}
	for _, binding := range bindings.ServiceBindings {
		bd, err := gb.Client.GetBindingByID(binding.ID, &gb.Parameters)
		if err != nil {
			// The binding could be deleted after List and before Get
			if strings.Contains(err.Error(), "StatusCode: 404") {
				continue
			}
			return err
		}
		instance, err := gb.Client.GetInstanceByID(bd.ServiceInstanceID, &query.Parameters{})
		if err != nil {
			return err
		}
		bd.ServiceInstanceName = instance.Name
		resultBindings.ServiceBindings = append(resultBindings.ServiceBindings, *bd)
	}

	if len(resultBindings.ServiceBindings) < 1 {
		output.PrintMessage(gb.Output, "No binding found with name: %s", gb.bindingName)
		return nil
	}

	output.PrintServiceManagerObject(gb.Output, gb.outputFormat, resultBindings)
	output.Println(gb.Output)

	return nil
}


func (gb *GetBindingCmd) printParameters(bindings *types.ServiceBindings) error {
	for _, binding := range bindings.ServiceBindings {
		parameters, err := gb.Client.GetBindingParameters(binding.ID, &gb.Parameters)
		if err != nil {
			// The binding could be deleted after List and before Get
			if strings.Contains(err.Error(), "StatusCode: 404") {
				continue
			}
			output.PrintMessage(gb.Output, "Unable to show parameters for service binding id: %s", binding.ID)
			output.PrintMessage(gb.Output, "The error is: %s", err)
			continue
		}
		if len(parameters) == 0 {
			output.PrintMessage(gb.Output, "No parameters are set for service binding id: %s", binding.ID)
			continue
		}
		output.PrintMessage(gb.Output, "Showing parameters for service binding id: %s", binding.ID)
		output.PrintMessage(gb.Output, "The parameters are: %s", parameters)
	}

	output.Println(gb.Output)
	return nil
}


// Validate validates command's arguments
func (gb *GetBindingCmd) Validate(args []string) error {
	if len(args) < 1 || len(args[0]) < 1 {
		return fmt.Errorf("binding name is required")
	}

	gb.bindingName = args[0]

	return nil
}

// SetOutputFormat set output format
func (gb *GetBindingCmd) SetOutputFormat(format output.Format) {
	gb.outputFormat = format
}

// HideUsage hide command's usage
func (gb *GetBindingCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (gb *GetBindingCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "get-binding [name]",
		Aliases: []string{"gsb"},
		Short:   "Get single binding",
		Long:    `Get single binding by its name`,
		PreRunE: prepare(gb, gb.Context),
		RunE:    cmd.RunE(gb),
	}

	result.Flags().BoolVar(&gb.bindingParams, "binding-params", false, "Get service binding params")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &gb.Parameters)

	return result
}
