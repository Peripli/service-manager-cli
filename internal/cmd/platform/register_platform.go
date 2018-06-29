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
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/spf13/cobra"
)

// RegisterPlatformCmd wraps the smctl register-platform command
type RegisterPlatformCmd struct {
	*cmd.Context

	platform types.Platform

	outputFormat output.Format
}

// NewRegisterPlatformCmd returns new register-platform command with context
func NewRegisterPlatformCmd(context *cmd.Context) *RegisterPlatformCmd {
	return &RegisterPlatformCmd{Context: context, platform: types.Platform{}}
}

// SetOutputFormat set command's output format
func (rpc *RegisterPlatformCmd) SetOutputFormat(format output.Format) {
	rpc.outputFormat = format
}

// HideUsage hide command's usage
func (rpc *RegisterPlatformCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (rpc *RegisterPlatformCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "register-platform [name] [type] <description>",
		Aliases: []string{"rp"},
		Short:   "Registers a platform",
		Long:    `Registers a platform`,

		PreRunE: prepare(rpc, rpc.Context),
		RunE:    cmd.RunE(rpc),
	}

	result.Flags().StringVarP(&rpc.platform.ID, "id", "i", "", "external platform ID")
	cmd.AddFormatFlag(result.Flags())

	return result
}

// Validate validates command's arguments
func (rpc *RegisterPlatformCmd) Validate(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("requires at least 2 args")
	}

	rpc.platform.Name = args[0]
	rpc.platform.Type = args[1]

	if len(args) > 2 {
		rpc.platform.Description = args[2]
	}

	return nil
}

// Run runs command's logic
func (rpc *RegisterPlatformCmd) Run() error {
	resultPlatform, err := rpc.Client.RegisterPlatform(&rpc.platform)
	if err != nil {
		return err
	}
	output.PrintServiceManagerObject(rpc.Output, rpc.outputFormat, resultPlatform)
	output.Println(rpc.Output)
	return nil
}
