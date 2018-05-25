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

package platform

import (
	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

// ListPlatformsCmd wraps the smctl list-brokers command
type ListPlatformsCmd struct {
	*cmd.Context

	outputFormat int
}

// NewListPlatformsCmd returns new list-brokers command with context
func NewListPlatformsCmd(context *cmd.Context) *ListPlatformsCmd {
	return &ListPlatformsCmd{Context: context}
}

func (lb *ListPlatformsCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list-platforms",
		Aliases: []string{"lp"},
		Short:   "List platforms",
		Long:    `List all platforms.`,
		PreRunE: cmd.PreRunE(lb, lb.Context),
		RunE:    cmd.RunE(lb),
	}
}

// Run runs the command's logic
func (lb *ListPlatformsCmd) Run() error {
	platforms, err := lb.Client.ListPlatforms()
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(lb.Output, lb.outputFormat, platforms)
	output.Println(lb.Output)

	return nil
}

func (lb *ListPlatformsCmd) addFlags(command *cobra.Command) *cobra.Command {
	cmd.AddFormatFlag(command.Flags())
	return command
}

// SetSMClient set the SM client
func (lb *ListPlatformsCmd) SetSMClient(client smclient.Client) {
	lb.Client = client
}

// SetOutputFormat set output format
func (lb *ListPlatformsCmd) SetOutputFormat(format int) {
	lb.outputFormat = format
}

// HideUsage hide command's usage
func (lb *ListPlatformsCmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (lb *ListPlatformsCmd) Command() *cobra.Command {
	result := lb.buildCommand()
	result = lb.addFlags(result)

	return result
}
