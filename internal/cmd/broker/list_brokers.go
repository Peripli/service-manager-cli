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

package broker

import (
	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// ListBrokersCmd wraps the smctl list-brokers command
type ListBrokersCmd struct {
	*cmd.Context

	prepare      cmd.PrepareFunc
	outputFormat output.Format
}

// NewListBrokersCmd returns new list-brokers command with context
func NewListBrokersCmd(context *cmd.Context) *ListBrokersCmd {
	return &ListBrokersCmd{Context: context}
}

// Run runs the command's logic
func (lb *ListBrokersCmd) Run() error {
	brokers, err := lb.Client.ListBrokers()
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(lb.Output, lb.outputFormat, brokers)
	output.Println(lb.Output)

	return nil
}

// SetOutputFormat set output format
func (lb *ListBrokersCmd) SetOutputFormat(format output.Format) {
	lb.outputFormat = format
}

// HideUsage hide command's usage
func (lb *ListBrokersCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (lb *ListBrokersCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	lb.prepare = prepare
	result := &cobra.Command{
		Use:     "list-brokers",
		Aliases: []string{"lb"},
		Short:   "List brokers",
		Long:    `List all brokers.`,
		PreRunE: lb.prepare(lb, lb.Context),
		RunE:    cmd.RunE(lb),
	}

	cmd.AddFormatFlag(result.Flags())

	return result
}
