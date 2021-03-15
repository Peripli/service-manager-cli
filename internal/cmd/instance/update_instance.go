/*
 * Copyright 2021 The Service Manager Authors
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
	"encoding/json"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"io"

	"fmt"

	"github.com/spf13/cobra"
)

// TransferCmd wraps the smctl provision command
type UpdateCmd struct {
	*cmd.Context

	input io.Reader

	instanceName   string
	instanceID     string
	planName       string
	planId         string
	parametersJSON string
	outputFormat   output.Format
}

// NewTransferCmd returns new transfer instance command with context
func NewUpdateInstanceCmd(context *cmd.Context, input io.Reader) *UpdateCmd {
	return &UpdateCmd{Context: context, input: input}
}

// Prepare returns cobra command
func (uc *UpdateCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "update-instance [name] --plan-name new-plan-name  --parameters new-configuration-parameters",
		Short: "Update a service instance",
		Long:  `Update an existing service instance. You can update its name, associated service plan, or configuration parameters`,

		PreRunE: prepare(uc, uc.Context),
		RunE:    cmd.RunE(uc),
	}

	result.Flags().StringVarP(&uc.instanceID, "id", "", "", "Id of the instance. Required in case when there are instances with same name")
	result.Flags().StringVarP(&uc.instanceName, "name", "", "", "The name of the service instance to update")
	result.Flags().StringVarP(&uc.planName, "plan-name", "", "", "The name of the new service plan to use for the instance")
	result.Flags().StringVarP(&uc.planId, "plan-id", "", "", "The name of the new service plan to use for the instance.")
	result.Flags().StringVarP(&uc.parametersJSON, "parameters", "c", "", "Valid JSON object containing instance configuration parameters")
	cmd.AddFormatFlag(result.Flags())

	return result
}

func (uc *UpdateCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("name is required")
	}
	uc.instanceName = args[0]

	return nil
}

func (uc *UpdateCmd) Run() error {
	var serviceOfferingId string = ""
	var servicePlanId string = ""
	if uc.instanceID == "" {
		instances, err := uc.Client.ListInstances(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", uc.instanceName),
			},
			GeneralParams: uc.Parameters.GeneralParams,
		})
		if err != nil {
			return err
		}
		if len(instances.ServiceInstances) == 0 {
			return fmt.Errorf("no instances found with name %s", uc.instanceName)
		}

		if len(instances.ServiceInstances) > 1 {
			return fmt.Errorf("more than 1 instance found with name %s. Use --id flag to specify one", uc.instanceName)
		}

		uc.instanceID = instances.ServiceInstances[0].ID
		serviceOfferingId = instances.ServiceInstances[0].ServiceID
	}
	plans, err := uc.Client.ListPlans(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", uc.planName),
			fmt.Sprintf("service_instance_id eq '%s'", serviceOfferingId),
		},
		GeneralParams: uc.Parameters.GeneralParams,
	})
	if err != nil {
		return err
	}
	if len(plans.ServicePlans) != 1 {
		return fmt.Errorf("exactly one service plan with name %s for offering with id %s expected", uc.planName, serviceOfferingId)
	}

	servicePlanId = plans.ServicePlans[0].ID
	resultInstance, location, err := uc.Client.UpdateInstance(uc.instanceID, &types.ServiceInstance{
		ID:            uc.instanceID,
		ServicePlanID: servicePlanId,
		Parameters:    json.RawMessage(uc.parametersJSON),
	}, nil)
	if err != nil {
		output.PrintMessage(uc.Output, "Could not update service instance. Reason: ")
		return err
	}

	if len(location) != 0 {
		cmd.CommonHandleAsyncExecution(uc.Context, location, fmt.Sprintf("Service Instance %s successfully scheduled for update. To see status of the operation use:\n", uc.instanceName, location))
		return nil
	}
	output.PrintServiceManagerObject(uc.Output, uc.outputFormat, resultInstance)
	output.Println(uc.Output)
	return nil
}

// SetOutputFormat set output format
func (uc *UpdateCmd) SetOutputFormat(format output.Format) {
	uc.outputFormat = format
}

// HideUsage hide command's usage
func (uc *UpdateCmd) HideUsage() bool {
	return true
}
