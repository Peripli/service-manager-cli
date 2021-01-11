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

package logout

import (
	"github.com/spf13/cobra"
	"time"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
)

// Cmd wraps the smctl version command
type Cmd struct {
	*cmd.Context
}

// NewLogoutCmd returns new version command
func NewLogoutCmd(context *cmd.Context) *Cmd {
	return &Cmd{context}
}

// Prepare returns cobra command
func (vc *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "logout",
		Aliases: []string{"v"},
		Short:   "Logs the user out",
		Long:    `Logs the user out and deletes the active client access token.`,

		PreRunE: prepare(vc, vc.Context),
		RunE:    cmd.RunE(vc),
	}

	return result
}

// Run runs command's logic
func (vc *Cmd) Run() error {
	config, err := vc.Configuration.Load()

	if config != nil && err != nil {
		if config.Token.AccessToken == "" {
			output.PrintMessage(vc.Output, "You are currently logged out.\n")
			return nil
		}
	}

	if err != nil {
		return err
	}

	config.RefreshToken = ""
	config.ExpiresIn = time.Time{}
	config.AccessToken = ""
	config.Scope = ""

	err = vc.Configuration.Save(config)
	if err != nil {
		return err
	}

	output.PrintMessage(vc.Output, "You have successfully logged out.\n")
	return nil
}
