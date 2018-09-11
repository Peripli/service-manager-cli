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
	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/internal/output"
)

// Cmd wraps the smctl info command
type Cmd struct {
	*cmd.Context
}

// NewInfoCmd returns new info command with context
func NewInfoCmd(context *cmd.Context) *Cmd {
	return &Cmd{context}
}

// Prepare returns the cobra command
func (ic *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Short:   "Prints information for logged user",
		Long:    `Prints information for logged user`,

		PreRunE: prepare(ic, ic.Context),
		RunE:    cmd.RunE(ic),
	}

	return result
}

// Run runs the command's logic
func (ic *Cmd) Run() error {
	settings := configuration.DefaultSettings()
	// TODO: check err properly
	err := ic.Configuration.Unmarshal(settings)
	clientConfig := settings.SMClient
	if err != nil {
		output.PrintMessage(ic.Output, "There is no logged user. Use \"smctl login\" to log in.\n")
	} else {
		output.PrintMessage(ic.Output, "Service Manager URL: %s\n", clientConfig.URL)
		output.PrintMessage(ic.Output, "Logged user: %s\n", clientConfig.User)
	}

	return nil
}
