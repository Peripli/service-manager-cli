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

package offering

import (
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/spf13/cobra"
)

// ListOfferingsCmd wraps the smctl list-offerings command
type ListOfferingsCmd struct {
	*cmd.Context

	prepare      cmd.PrepareFunc
	outputFormat output.Format

	offering string
}

// NewListOfferingsCmd returns new list-offerings command with context
func NewListOfferingsCmd(context *cmd.Context) *ListOfferingsCmd {
	return &ListOfferingsCmd{Context: context}
}

// Run runs the command's logic
func (lo *ListOfferingsCmd) Run() error {
	offerings, err := lo.Client.ListOfferingsWithQuery(lo.Parameters.Copy())
	if err != nil {
		return err
	}
	if lo.offering == "" {
		output.PrintServiceManagerObject(lo.Output, lo.outputFormat, offerings)
	} else {
		plans := &types.ServicePlans{}
		for _, v := range offerings.ServiceOfferings {
			if v.Name == lo.offering {
				plans.ServicePlans = append(plans.ServicePlans, v.Plans...)
			}
		}
		output.PrintServiceManagerObject(lo.Output, lo.outputFormat, plans)
	}
	output.Println(lo.Output)
	return nil
}

// SetOutputFormat set output format
func (lo *ListOfferingsCmd) SetOutputFormat(format output.Format) {
	lo.outputFormat = format
}

// HideUsage hide command's usage
func (lo *ListOfferingsCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (lo *ListOfferingsCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	lo.prepare = prepare
	result := &cobra.Command{
		Use:     "list-offerings",
		Aliases: []string{"lo"},
		Short:   "List service offerings",
		Long:    `List all service offerings.`,
		PreRunE: lo.prepare(lo, lo.Context),
		RunE:    cmd.RunE(lo),
	}

	cmd.AddFormatFlag(result.Flags())
	result.Flags().StringVarP(&lo.offering, "service", "s", "", "Plan details for a single service offering")
	cmd.AddQueryingFlags(result.Flags(), lo.Parameters)

	return result
}
