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
	"github.com/Peripli/service-manager-cli/internal/print"
	"github.com/Peripli/service-manager-cli/pkg/types"

	"github.com/Peripli/service-manager-cli/pkg/smclient"

	"github.com/spf13/cobra"
)

// Wraps the smctl register-platform command
type RegisterPlatformCmd struct {
	*cmd.Context

	platform types.Platform

	outputFormat int
}

func NewRegisterPlatformCmd(context *cmd.Context) *RegisterPlatformCmd {
	return &RegisterPlatformCmd{Context: context, platform: types.Platform{}}
}

func (rpc *RegisterPlatformCmd) SetSMClient(client smclient.Client) {
	rpc.Client = client
}

func (rpc *RegisterPlatformCmd) SetOutputFormat(format int) {
	rpc.outputFormat = format
}

func (rpc *RegisterPlatformCmd) HideUsage() bool {
	return true
}

func (rpc *RegisterPlatformCmd) Command() *cobra.Command {
	result := rpc.buildCommand()
	result = rpc.addFlags(result)

	return result
}

func (rpc *RegisterPlatformCmd) buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "register-platform [name] [type] <description>",
		Aliases: []string{"rp"},
		Short:   "Registers a platform",
		Long:    `Registers a platform`,

		PreRunE: cmd.PreRunE(rpc, rpc.Context),
		RunE:    cmd.RunE(rpc),
	}
}

func (rpc *RegisterPlatformCmd) addFlags(command *cobra.Command) *cobra.Command {
	command.Flags().StringVarP(&rpc.platform.ID, "id", "i", "", "external platform ID")
	cmd.AddFormatFlag(command.Flags())

	return command
}

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

func (rpc *RegisterPlatformCmd) Run() error {
	resultPlatform, err := rpc.Client.RegisterPlatform(&rpc.platform)
	if err != nil {
		return err
	}
	print.PrintServiceManagerObject(rpc.Output, rpc.outputFormat, resultPlatform)
	print.Println(rpc.Output)
	return nil
}