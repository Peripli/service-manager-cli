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
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
)

// ListInstancesCmd wraps the smctl list-instances command
type ListInstancesCmd struct {
	*cmd.Context

	outputFormat output.Format
}

// NewListInstancesCmd returns new list-instances command with context
func NewListInstancesCmd(context *cmd.Context) *ListInstancesCmd {
	return &ListInstancesCmd{Context: context}
}

func (li *ListInstancesCmd) Run() error {
	instances, err := li.Client.ListInstances(&li.Parameters)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(li.Output, li.outputFormat, instances)
	output.Println(li.Output)

	return nil
}

// SetOutputFormat sets output format
func (li *ListInstancesCmd) SetOutputFormat(format output.Format) {
	li.outputFormat = format
}

// HideUsage hides command's usage
func (li *ListInstancesCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (li *ListInstancesCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "list-instances",
		Aliases: []string{"li"},
		Short:   "List service-instances",
		Long:    `List all service-instances.`,
		PreRunE: prepare(li, li.Context),
		RunE:    cmd.RunE(li),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddQueryingFlags(result.Flags(), &li.Parameters)
	cmd.AddCommonQueryFlag(result.Flags(), &li.Parameters)

	return result
}
