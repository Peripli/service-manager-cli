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

	input           io.Reader
	instance        types.ServiceInstance
	instanceName    string
	planName        string
	parametersJSON  string
	outputFormat    output.Format
}

// NewTransferCmd returns new transfer instance command with context
func NewUpdateInstanceCmd(context *cmd.Context, input io.Reader) *UpdateCmd {
	return &UpdateCmd{Context: context, input: input, instance: types.ServiceInstance{}}
}

// Prepare returns cobra command
func (uc *UpdateCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "update-instance [name] --id service-instance-id --new-name new-service-name --plan new-plan-name --instance-params new-configuration-parameters",
		Short: "Update a service instance",
		Long:  `Update the name, associated plan, and configuration parameters of an existing service instance`,

		PreRunE: prepare(uc, uc.Context),
		RunE:    cmd.RunE(uc),
	}
	result.Flags().StringVarP(&uc.instance.ID, "id", "", "", "The id of the service instance to update")
	result.Flags().StringVarP(&uc.instance.Name, "new-name", "", "", "The new name of the service instance")
	result.Flags().StringVarP(&uc.planName, "plan", "", "", "The name of the new service plan to use for the instance")
	result.Flags().StringVarP(&uc.parametersJSON, "instance-params", "c", "", "Valid JSON object containing instance configuration parameters")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddModeFlag(result.Flags(), "async")
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
	var instanceBeforeUpdate *types.ServiceInstance
	if uc.instance.ID == "" {
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
		instanceBeforeUpdate = &instances.ServiceInstances[0]

	} else {
		instance, err := uc.Client.GetInstanceByID(uc.instance.ID, nil)
		if err != nil {
			return err
		}
		instanceBeforeUpdate = instance

	}

	if uc.planName != "" {
		currentPlan, err := uc.Client.GetPlanByID(instanceBeforeUpdate.ServicePlanID, &uc.Parameters)
		if err != nil {
			return err
		}

		plans, err := uc.Client.ListPlans(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("catalog_name eq '%s'", uc.planName),
				fmt.Sprintf("service_offering_id eq '%s'", currentPlan.ServiceOfferingID),
			},
			GeneralParams: uc.Parameters.GeneralParams,
		})

		if err != nil {
			return err
		}
		if len(plans.ServicePlans) == 0 {
			return fmt.Errorf("service plan with name %s for offering with id %s not found", uc.planName, currentPlan.ServiceOfferingID)
		}
		if len(plans.ServicePlans) > 1 {
			return fmt.Errorf("exactly one service plan with name %s for offering with id %s expected", uc.planName, currentPlan.ServiceOfferingID)
		}

		uc.instance.ServicePlanID = plans.ServicePlans[0].ID
	}
	if len(uc.parametersJSON) > 0 {
		uc.instance.Parameters = json.RawMessage(uc.parametersJSON)
	}

	resultInstance, location, err := uc.Client.UpdateInstance(instanceBeforeUpdate.ID, &uc.instance, &uc.Parameters)
	if err != nil {
		output.PrintMessage(uc.Output, "Could not update service instance. Reason: ")
		return err
	}

	if len(location) != 0 {
		cmd.CommonHandleAsyncExecution(uc.Context, location, fmt.Sprintf("Service Instance %s successfully scheduled for update. To see status of the operation use:\n", uc.instance.Name))
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
