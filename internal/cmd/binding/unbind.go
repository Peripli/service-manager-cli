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
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/spf13/cobra"
	"io"
)

// UnbindCmd wraps the smctl bind command
type UnbindCmd struct {
	*cmd.Context

	input io.Reader
	force bool

	instanceName string
	bindingID    string
	bindingName  string

	outputFormat output.Format
}

// NewUnbindCmd returns new bind command with context
func NewUnbindCmd(context *cmd.Context, input io.Reader) *UnbindCmd {
	return &UnbindCmd{Context: context, input: input}
}

// Prepare returns cobra command
func (ubc *UnbindCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "unbind [instance-name] [binding-name]",
		Short: "Deletes a binding from SM",
		Long:  `Deletes a binding from SM`,

		PreRunE: prepare(ubc, ubc.Context),
		RunE:    cmd.RunE(ubc),
	}

	result.Flags().StringVarP(&ubc.bindingID, "id", "", "", "ID of the service binding. Required when name is ambiguous")
	result.Flags().BoolVarP(&ubc.force, "force", "f", false, "Force delete without confirmation")
	cmd.AddCommonQueryFlag(result.Flags(), &ubc.Parameters)
	cmd.AddSyncFlag(result.Flags())

	return result
}

// Validate validates command's arguments
func (ubc *UnbindCmd) Validate(args []string) error {
	if ubc.bindingID != "" {
		return nil
	}

	if len(args) < 2 {
		return fmt.Errorf("instance and binding names are required")
	}

	ubc.instanceName = args[0]
	ubc.bindingName = args[1]
	return nil
}

// Run runs the command's logic
func (ubc *UnbindCmd) Run() error {
	if ubc.bindingID == "" {
		instanceToUnbind, err := ubc.Client.ListInstances(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", ubc.instanceName),
			},
		})
		if err != nil {
			return err
		}
		if len(instanceToUnbind.ServiceInstances) < 1 {
			return fmt.Errorf("service instance with name %s not found", ubc.instanceName)
		}
		if len(instanceToUnbind.ServiceInstances) > 1 {
			return fmt.Errorf("more than one service instance with name %s found. Use --id flag to specify id of the binding to be deleted", ubc.instanceName)
		}

		bindingsToDelte, err := ubc.Client.ListBindings(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", ubc.bindingName),
				fmt.Sprintf("service_instance_id eq '%s'", instanceToUnbind.ServiceInstances[0].ID),
			},
		})
		if err != nil {
			return err
		}
		if len(bindingsToDelte.ServiceBindings) < 1 {
			output.PrintMessage(ubc.Output, "Service Binding with name %s for instance with name %s not found", ubc.bindingName, ubc.instanceName)
			return nil
		}
		ubc.bindingID = bindingsToDelte.ServiceBindings[0].ID
	}

	location, err := ubc.Client.Unbind(ubc.bindingID, &ubc.Parameters)
	if err != nil {
		output.PrintMessage(ubc.Output, "Could not delete service binding. Reason: ")
		return err
	}
	if len(location) != 0 {
		output.PrintMessage(ubc.Output, "Service Binding %s successfully scheduled for deletion. To see status of the operation use:\n", ubc.bindingName)
		output.PrintMessage(ubc.Output, "smctl poll %s\n", location)
		return nil
	}
	output.PrintMessage(ubc.Output, "Service Binding successfully deleted.\n")
	return nil
}

// AskForConfirmation asks the user to confirm deletion
func (ubc *UnbindCmd) AskForConfirmation() (bool, error) {
	if !ubc.force {
		message := fmt.Sprintf("Do you really want to delete binding with name [%s] for instance with name %s (Y/n): ", ubc.bindingName, ubc.instanceName)
		return cmd.CommonConfirmationPrompt(message, ubc.Context, ubc.input)
	}
	return true, nil
}

// PrintDeclineMessage prints confirmation decline message to the user
func (ubc *UnbindCmd) PrintDeclineMessage() {
	cmd.CommonPrintDeclineMessage(ubc.Output)
}

// HideUsage hide command's usage
func (ubc *UnbindCmd) HideUsage() bool {
	return true
}
