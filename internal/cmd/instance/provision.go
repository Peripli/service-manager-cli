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
	"encoding/json"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/query"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"fmt"
	"github.com/spf13/cobra"
)

// ProvisionCmd wraps the smctl provision command
type ProvisionCmd struct {
	*cmd.Context

	instance       types.ServiceInstance
	offeringName   string
	planName       string
	brokerName     string
	parametersJSON string

	outputFormat output.Format
}

// NewProvisionCmd returns new provision command with context
func NewProvisionCmd(context *cmd.Context) *ProvisionCmd {
	return &ProvisionCmd{Context: context, instance: types.ServiceInstance{}}
}

// Prepare returns cobra command
func (pi *ProvisionCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "provision [name] [offering] [plan]",
		Short: "Provisions an instance in SM",
		Long:  `Provisions an instance in SM`,

		PreRunE: prepare(pi, pi.Context),
		RunE:    cmd.RunE(pi),
	}

	result.Flags().StringVarP(&pi.brokerName, "broker-name", "b", "", "Name of the broker which provides the service offering. Required when offering name is ambiguous")
	result.Flags().StringVarP(&pi.parametersJSON, "parameters", "c", "", "Valid JSON object containing instance parameters")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &pi.Parameters)
	cmd.AddSyncFlag(result.Flags())

	return result
}

// Validate validates command's arguments
func (pi *ProvisionCmd) Validate(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("name, offering and plan are required")
	}

	pi.instance.Name = args[0]
	pi.offeringName = args[1]
	pi.planName = args[2]

	return nil
}

// Run runs the command's logic
func (pi *ProvisionCmd) Run() error {
	offerings, err := pi.Client.ListOfferings(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", pi.offeringName),
		},
	})
	if err != nil {
		return err
	}
	if len(offerings.ServiceOfferings) == 0 {
		return fmt.Errorf("service offering with name %s not found", pi.offeringName)
	}

	pi.instance.ServiceID = offerings.ServiceOfferings[0].ID

	if len(offerings.ServiceOfferings) > 1 {
		if len(pi.brokerName) == 0 {
			return fmt.Errorf("more than one service offering with name %s found. Use -b flag to specify broker name", pi.offeringName)
		}

		brokers, err := pi.Client.ListBrokers(&query.Parameters{
			FieldQuery: []string{
				fmt.Sprintf("name eq '%s'", pi.brokerName),
			},
		})
		if err != nil {
			return err
		}
		if len(brokers.Brokers) != 1 {
			return fmt.Errorf("exactly one broker with name %s expected, found %d", pi.brokerName, len(brokers.Brokers))
		}
		for _, offering := range offerings.ServiceOfferings {
			if offering.BrokerID == brokers.Brokers[0].ID {
				pi.instance.ServiceID = offering.ID
				break
			}
		}
	}

	plans, err := pi.Client.ListPlans(&query.Parameters{
		FieldQuery: []string{
			fmt.Sprintf("name eq '%s'", pi.planName),
			fmt.Sprintf("service_offering_id eq '%s'", pi.instance.ServiceID),
		},
	})
	if err != nil {
		return err
	}
	if len(plans.ServicePlans) != 1 {
		return fmt.Errorf("exactly one service plan with name %s for offering with id %s expected", pi.planName, pi.instance.ServiceID)
	}

	pi.instance.ServicePlanID = plans.ServicePlans[0].ID
	pi.instance.Parameters = json.RawMessage(pi.parametersJSON)

	resultInstance, location, err := pi.Client.Provision(&pi.instance, &pi.Parameters)
	if err != nil {
		return err
	}

	if len(location) != 0 {
		output.PrintMessage(pi.Output, "Service Instance %s successfully scheduled for provisioning. To see status of the operation use:\n", pi.instance.Name)
		output.PrintMessage(pi.Output, "smctl poll %s\n", location)
		return nil
	}
	output.PrintServiceManagerObject(pi.Output, pi.outputFormat, resultInstance)
	output.Println(pi.Output)
	return nil
}

// SetOutputFormat set output format
func (pi *ProvisionCmd) SetOutputFormat(format output.Format) {
	pi.outputFormat = format
}

// HideUsage hide command's usage
func (pi *ProvisionCmd) HideUsage() bool {
	return true
}
