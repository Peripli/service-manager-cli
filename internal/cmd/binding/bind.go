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
	"encoding/json"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"fmt"
	"github.com/spf13/cobra"
)

// BindCmd wraps the smctl bind command
type BindCmd struct {
	*cmd.Context

	binding        types.ServiceBinding
	instanceName   string
	parametersJSON string

	outputFormat output.Format
}

// NewBindCmd returns new bind command with context
func NewBindCmd(context *cmd.Context) *BindCmd {
	return &BindCmd{Context: context, binding: types.ServiceBinding{}}
}

// Prepare returns cobra command
func (bc *BindCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "bind [instance-name] [binding-name]",
		Short: "Creates binding in SM",
		Long:  `Creates binding in SM`,

		PreRunE: prepare(bc, bc.Context),
		RunE:    cmd.RunE(bc),
	}

	result.Flags().StringVarP(&bc.binding.ServiceInstanceID, "id", "", "", "ID of the service instance. Required when name is ambiguous")
	result.Flags().StringVarP(&bc.parametersJSON, "parameters", "c", "", "Valid JSON object containing binding parameters")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &bc.Parameters)
	cmd.AddModeFlag(result.Flags(), "async")

	return result
}

// Validate validates command's arguments
func (bc *BindCmd) Validate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("instance and binding names are required")
	}

	bc.instanceName = args[0]
	bc.binding.Name = args[1]
	return nil
}

// Run runs the command's logic
func (bc *BindCmd) Run() error {
	if bc.binding.ServiceInstanceID == "" {
		instanceToBind, err := bc.Client.ListInstances(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", bc.instanceName),
			},
			GeneralParams: bc.Parameters.GeneralParams,
		})
		if err != nil {
			return err
		}
		if len(instanceToBind.ServiceInstances) < 1 {
			return fmt.Errorf("service instance with name %s not found", bc.instanceName)
		}
		if len(instanceToBind.ServiceInstances) > 1 {
			return fmt.Errorf("more than one service instance with name %s found. Use --id flag to specify id of the instance to bind", bc.instanceName)
		}
		bc.binding.ServiceInstanceID = instanceToBind.ServiceInstances[0].ID
	}

	bc.binding.Parameters = json.RawMessage(bc.parametersJSON)
	resultBinding, location, err := bc.Client.Bind(&bc.binding, &bc.Parameters)
	if err != nil {
		return err
	}

	if len(location) != 0 {
		cmd.CommonHandleAsyncExecution(bc.Context, location, fmt.Sprintf("Service Binding %s successfully scheduled. To see status of the operation use:\n", bc.binding.Name))
		return nil
	}

	resultBinding.ServiceInstanceName = bc.instanceName
	output.PrintServiceManagerObject(bc.Output, bc.outputFormat, resultBinding)
	output.Println(bc.Output)
	return nil
}

// SetOutputFormat set output format
func (bc *BindCmd) SetOutputFormat(format output.Format) {
	bc.outputFormat = format
}

// HideUsage hide command's usage
func (bc *BindCmd) HideUsage() bool {
	return true
}
