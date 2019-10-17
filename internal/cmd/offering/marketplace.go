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

// MarketplaceCmd wraps the smctl marketplace command
type MarketplaceCmd struct {
	*cmd.Context

	prepare      cmd.PrepareFunc
	outputFormat output.Format

	offering string
}

// NewMarketplaceCmd returns new list-offerings command with context
func NewMarketplaceCmd(context *cmd.Context) *MarketplaceCmd {
	return &MarketplaceCmd{Context: context}
}

// Run runs the command's logic
func (m *MarketplaceCmd) Run() error {
	marketplace, err := m.Client.Marketplace(&m.Parameters)
	if err != nil {
		return err
	}
	if m.offering == "" {
		output.PrintServiceManagerObject(m.Output, m.outputFormat, marketplace)
	} else {
		plans := &types.ServicePlansForOffering{}
		for _, v := range marketplace.ServiceOfferings {
			if v.Name == m.offering {
				plans.ServicePlans = append(plans.ServicePlans, v.Plans...)
			}
		}
		output.PrintServiceManagerObject(m.Output, m.outputFormat, plans)
	}
	output.Println(m.Output)
	return nil
}

// SetOutputFormat set output format
func (m *MarketplaceCmd) SetOutputFormat(format output.Format) {
	m.outputFormat = format
}

// HideUsage hide command's usage
func (m *MarketplaceCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (m *MarketplaceCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	m.prepare = prepare
	result := &cobra.Command{
		Use:     "marketplace",
		Aliases: []string{"m"},
		Short:   "Shows marketplace for all the service-offerings",
		Long:    `Shows marketplace for all the service-offerings`,
		PreRunE: m.prepare(m, m.Context),
		RunE:    cmd.RunE(m),
	}

	cmd.AddFormatFlag(result.Flags())
	result.Flags().StringVarP(&m.offering, "service", "s", "", "Plan details for a single service offering")
	cmd.AddCommonQueryFlag(result.Flags(), &m.Parameters)

	return result
}
