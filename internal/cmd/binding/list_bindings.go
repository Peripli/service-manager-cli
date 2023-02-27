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
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
)

// ListBindingsCmd wraps the smctl list-bindings command
type ListBindingsCmd struct {
	*cmd.Context

	outputFormat output.Format
}

// NewListBindingsCmd returns new list-bindings command with context
func NewListBindingsCmd(context *cmd.Context) *ListBindingsCmd {
	return &ListBindingsCmd{Context: context}
}

// Run runs the command's logic
func (li *ListBindingsCmd) Run() error {
	bindings, err := li.Client.ListBindings(&li.Parameters)
	if err != nil {
		return err
	}

	for i := range bindings.ServiceBindings {
		instance, err := li.Client.GetInstanceByID(bindings.ServiceBindings[i].ServiceInstanceID, &li.Parameters)
		if err != nil {
			return err
		}
		bindings.ServiceBindings[i].ServiceInstanceName = instance.Name
	}

	output.PrintServiceManagerObject(li.Output, li.outputFormat, bindings)
	output.Println(li.Output)

	return nil
}

// SetOutputFormat sets output format
func (li *ListBindingsCmd) SetOutputFormat(format output.Format) {
	li.outputFormat = format
}

// HideUsage hides command's usage
func (li *ListBindingsCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (li *ListBindingsCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "list-bindings",
		Aliases: []string{"lsb"},
		Short:   "List service bindings",
		Long:    `Show all service bindings associated with a Service Manager instance.`,
		PreRunE: prepare(li, li.Context),
		RunE:    cmd.RunE(li),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddQueryingFlags(result.Flags(), &li.Parameters)
	cmd.AddCommonQueryFlag(result.Flags(), &li.Parameters)

	return result
}
