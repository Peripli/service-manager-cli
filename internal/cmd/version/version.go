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

package version

import (
	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/print"
)

// Wraps the smctl version command
type VersionCmd struct {
	*cmd.Context

	clientVersion string
}

func NewVersionCmd(context *cmd.Context, clientVersion string) *VersionCmd {
	return &VersionCmd{context, clientVersion}
}

func (vc *VersionCmd) Command() *cobra.Command {
	result := vc.buildCommand()

	return result
}

func (vc *VersionCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Prints smctl version",
		Long:    `Prints smctl version.`,

		PreRunE: cmd.PreRunE(vc, vc.Context),
		RunE:    cmd.RunE(vc),
	}
}

func (vc *VersionCmd) Run() error {
	print.PrintMessage(vc.Output, "Service Manager Client %s\n", vc.clientVersion)

	return nil
}
