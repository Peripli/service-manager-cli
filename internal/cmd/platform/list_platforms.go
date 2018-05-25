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

func (lp *ListPlatformsCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list-platforms",
		Aliases: []string{"lp"},
		Short:   "List platforms",
		Long:    `List all platforms.`,
		PreRunE: cmd.PreRunE(lp, lp.Context),
		RunE:    cmd.RunE(lp),
	}
}

// Run runs the command's logic
func (lp *ListPlatformsCmd) Run() error {
	platforms, err := lp.Client.ListPlatforms()
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(lp.Output, lp.outputFormat, platforms)
	output.Println(lp.Output)

	return nil
}

func (lp *ListPlatformsCmd) addFlags(command *cobra.Command) *cobra.Command {
	cmd.AddFormatFlag(command.Flags())
	return command
}

// SetSMClient set the SM client
func (lp *ListPlatformsCmd) SetSMClient(client smclient.Client) {
	lp.Client = client
}

// SetOutputFormat set output format
func (lp *ListPlatformsCmd) SetOutputFormat(format int) {
	lp.outputFormat = format
}

// HideUsage hide command's usage
func (lp *ListPlatformsCmd) HideUsage() bool {
	return true
}

// Command returns cobra command
func (lp *ListPlatformsCmd) Command() *cobra.Command {
	result := lp.buildCommand()
	result = lp.addFlags(result)

	return result
}
