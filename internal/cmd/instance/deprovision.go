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

package instance

import (
	"fmt"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager/pkg/web"
	"io"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// DeprovisionCmd wraps the smctl deprovision command
type DeprovisionCmd struct {
	*cmd.Context

	input io.Reader
	force bool
	forceDelete bool

	name string
	id   string
}

// NewDeprovisionCmd returns new deprovision command with context
func NewDeprovisionCmd(context *cmd.Context, input io.Reader) *DeprovisionCmd {
	return &DeprovisionCmd{Context: context, input: input}
}

// Validate validates command's arguments
func (dbc *DeprovisionCmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("single [name] is required")
	}

	dbc.name = args[0]

	return nil
}

// Run runs the command's logic
func (dbc *DeprovisionCmd) Run() error {
	if dbc.id == "" {
		toDeprovision, err := dbc.Client.ListInstances(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", dbc.name),
			},
			GeneralParams: dbc.Parameters.GeneralParams,
		})
		if err != nil {
			return err
		}
		if len(toDeprovision.ServiceInstances) < 1 {
			output.PrintMessage(dbc.Output, "Service Instance not found.\n")
			return nil
		}
		if len(toDeprovision.ServiceInstances) > 1 {
			return fmt.Errorf("more than one service instance with name %s found. Use --id flag to specify id of the instance to be deleted", dbc.name)
		}
		dbc.id = toDeprovision.ServiceInstances[0].ID
	}


	if dbc.forceDelete {
		dbc.Parameters.GeneralParams = append(dbc.Parameters.GeneralParams, fmt.Sprintf("%s=%s", web.QueryParamCascade, "true"))
		dbc.Parameters.GeneralParams = append(dbc.Parameters.GeneralParams, fmt.Sprintf("%s=%s", web.QueryParamForce, "true"))
	}
	location, err := dbc.Client.Deprovision(dbc.id, &dbc.Parameters)
	if err != nil {
		output.PrintMessage(dbc.Output, "Could not delete service instance. Reason: ")
		return err
	}
	if len(location) != 0 {
		cmd.CommonHandleAsyncExecution(dbc.Context, location, fmt.Sprintf("Service Instance %s successfully scheduled for deletion. To see status of the operation use:\n", dbc.name))
		return nil
	}
	output.PrintMessage(dbc.Output, "Service Instance successfully deleted.\n")
	return nil
}

// HideUsage hide command's usage
func (dbc *DeprovisionCmd) HideUsage() bool {
	return true
}

// AskForConfirmation asks the user to confirm deletion
func (dbc *DeprovisionCmd) AskForConfirmation() (bool, error) {
	if !dbc.force {
		message := fmt.Sprintf("Do you really want to delete instance with name [%s] (Y/n): ", dbc.name)
		return cmd.CommonConfirmationPrompt(message, dbc.Context, dbc.input)
	}
	return true, nil
}

// PrintDeclineMessage prints confirmation decline message to the user
func (dbc *DeprovisionCmd) PrintDeclineMessage() {
	cmd.CommonPrintDeclineMessage(dbc.Output)
}

// Prepare returns cobra command
func (dbc *DeprovisionCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "deprovision [name]",
		Short:   "Deletes service instance",
		Long:    `Deletes service instance by name.`,
		PreRunE: prepare(dbc, dbc.Context),
		RunE:    cmd.RunE(dbc),
	}

	forceUsage := "Use this parameter to delete a resource without raising a confirmation message."
	forceDeleteUsage := "Delete the service instance and all of its associated resources from the database, including all service bindings. Use this parameter if the service instance cannot be properly deleted. This parameter can only be used by operators with technical access."
	result.Flags().BoolVarP(&dbc.force, "force", "f", false, forceUsage)
	result.Flags().BoolVarP(&dbc.forceDelete, "force-delete", "", false, forceDeleteUsage)
	result.Flags().StringVarP(&dbc.id, "id", "", "", "ID of the service instance. Required when name is ambiguous.")
	cmd.AddCommonQueryFlag(result.Flags(), &dbc.Parameters)
	cmd.AddModeFlag(result.Flags(), "async")

	return result
}
