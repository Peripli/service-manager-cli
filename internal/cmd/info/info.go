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

package info

import (
	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/print"
)

// Wraps the smctl info command
type InfoCmd struct {
	*cmd.Context
}

func NewInfoCmd(context *cmd.Context) *InfoCmd {
	return &InfoCmd{context}
}

func (ic *InfoCmd) Command() *cobra.Command {
	result := ic.buildCommand()

	return result
}

func (ic *InfoCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Short:   "Prints information for logged user",
		Long:    `Prints information for logged user`,

		PreRunE: cmd.PreRunE(ic, ic.Context),
		RunE:    cmd.RunE(ic),
	}
}

func (ic *InfoCmd) Run() error {
	clientConfig, err := ic.Configuration.Load()
	if err != nil {
		print.PrintMessage(ic.Output, "There is no logged user. Use \"smctl login\" to log in.\n")
	} else {
		print.PrintMessage(ic.Output, "Service Manager URL: %s\n", clientConfig.URL)
		print.PrintMessage(ic.Output, "Logged user: %s\n", clientConfig.User)
	}

	return nil
}
