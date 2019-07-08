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
	"github.com/Peripli/service-manager-cli/internal/output"
)

// Cmd wraps the smctl version command
type Cmd struct {
	*cmd.Context
}

// Version is the tool version, injected by the build
var Version = "local.build"

// GitCommit is the git commit id, injected by the build
var GitCommit string

// NewVersionCmd returns new version command
func NewVersionCmd(context *cmd.Context) *Cmd {
	return &Cmd{context}
}

// Prepare returns cobra command
func (vc *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Prints smctl version",
		Long:    `Prints smctl version.`,

		PreRunE: prepare(vc, vc.Context),
		RunE:    cmd.RunE(vc),
	}

	return result
}

// Run runs command's logic
func (vc *Cmd) Run() error {
	output.PrintMessage(vc.Output, "Service Manager Client %s (%s)\n", Version, GitCommit)

	return nil
}
