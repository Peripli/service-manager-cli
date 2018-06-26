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

package main

import (
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/cmd/broker"
	"github.com/Peripli/service-manager-cli/internal/cmd/info"
	"github.com/Peripli/service-manager-cli/internal/cmd/login"
	"github.com/Peripli/service-manager-cli/internal/cmd/platform"
	"github.com/Peripli/service-manager-cli/internal/cmd/version"
	"github.com/spf13/cobra"

	"os"
)

func main() {
	clientVersion := "0.0.1"

	context := &cmd.Context{}
	rootCmd := cmd.BuildRootCommand(context)

	normalCommandsGroup := cmd.Group{
		Commands: []cmd.CommandPreparator{
			login.NewLoginCmd(context, os.Stdin),
			version.NewVersionCmd(context, clientVersion),
			info.NewInfoCmd(context),
		},
		PrepareFn: cmd.CommonPrepare,
	}

	smCommandsGroup := cmd.Group{
		Commands: []cmd.CommandPreparator{
			broker.NewRegisterBrokerCmd(context),
			broker.NewListBrokersCmd(context),
			broker.NewDeleteBrokerCmd(context),
			broker.NewUpdateBrokerCmd(context),
			platform.NewRegisterPlatformCmd(context),
			platform.NewListPlatformsCmd(context),
			platform.NewDeletePlatformCmd(context),
			platform.NewUpdatePlatformCmd(context),
		},
		PrepareFn: cmd.SmPrepare,
	}

	registerGroups(rootCmd, normalCommandsGroup, smCommandsGroup)

	cmd.Execute(rootCmd)
}

func registerGroups(rootCmd *cobra.Command, groups ...cmd.Group) {
	for _, group := range groups {
		for _, command := range group.Commands {
			cobraCmd := command.Prepare(group.PrepareFn)
			rootCmd.AddCommand(cobraCmd)
		}
	}
}
