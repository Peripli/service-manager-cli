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
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

// ListBrokersCmd wraps the smctl list-brokers command
type ListBrokersCmd struct {
	*cmd.Context

	outputFormat int
	names        []string
}

// NewListBrokersCmd returns new list-brokers command with context
func NewListBrokersCmd(context *cmd.Context) *ListBrokersCmd {
	return &ListBrokersCmd{Context: context}
}

func (lb *ListBrokersCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list-brokers <name1 <name2> ... <nameN>>",
		Aliases: []string{"lb"},
		Short:   "List brokers",
		Long:    `List all brokers.`,
		PreRunE: cmd.PreRunE(lb, lb.Context),
		RunE:    cmd.RunE(lb),
	}
}

// Validate validates command's arguments
func (lb *ListBrokersCmd) Validate(args []string) error {
	if len(args) > 0 {
		lb.names = args
	}

	return nil
}

// Run runs the command's logic
func (lb *ListBrokersCmd) Run() error {
	brokers, err := lb.Client.ListBrokers(lb.names)
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(lb.Output, lb.outputFormat, brokers)
	output.Println(lb.Output)

	return nil
}

func (lb *ListBrokersCmd) addFlags(command *cobra.Command) *cobra.Command {
	cmd.AddFormatFlag(command.Flags())
	return command
}

// SetSMClient set the SM client
func (lb *ListBrokersCmd) SetSMClient(client smclient.Client) {
	lb.Client = client
}

// SetOutputFormat set output format
func (lb *ListBrokersCmd) SetOutputFormat(format int) {
	lb.outputFormat = format
}

// HideUsage hide command's usage
func (lb *ListBrokersCmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (lb *ListBrokersCmd) Command() *cobra.Command {
	result := lb.buildCommand()
	result = lb.addFlags(result)

	return result
}
