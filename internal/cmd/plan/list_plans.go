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

package plan

import (
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
)

// ListPlansCmd wraps the smctl list-plans command
type ListPlansCmd struct {
	*cmd.Context
	environment string
	outputFormat output.Format
}

// NewListPlansCmd returns new list-plans command with context
func NewListPlansCmd(context *cmd.Context) *ListPlansCmd {
	return &ListPlansCmd{Context: context}
}

// Run runs the command's logic
func (lp *ListPlansCmd) Run() error {
	plans, err := lp.Client.ListPlans(&lp.Parameters)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(lp.Output, lp.outputFormat, plans)
	output.Println(lp.Output)

	return nil
}

// SetOutputFormat set output format
func (lp *ListPlansCmd) SetOutputFormat(format output.Format) {
	lp.outputFormat = format
}

// HideUsage hide command's usage
func (lp *ListPlansCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (lp *ListPlansCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "list-plans",
		Short:   "List service-plans",
		Long:    `List service-plans.`,
		PreRunE: prepare(lp, lp.Context),
		RunE:    cmd.RunE(lp),
	}
	cmd.AddSupportedEnvironmentFlag(result.Flags(), &lp.Parameters, "Filters service plans by supported environments")
	cmd.AddFormatFlag(result.Flags())
	cmd.AddQueryingFlags(result.Flags(), &lp.Parameters)
	cmd.AddCommonQueryFlag(result.Flags(), &lp.Parameters)

	return result
}
