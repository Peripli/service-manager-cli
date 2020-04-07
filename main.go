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
	"context"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/cmd/binding"
	"github.com/Peripli/service-manager-cli/internal/cmd/broker"
	"github.com/Peripli/service-manager-cli/internal/cmd/curl"
	"github.com/Peripli/service-manager-cli/internal/cmd/info"
	"github.com/Peripli/service-manager-cli/internal/cmd/instance"
	"github.com/Peripli/service-manager-cli/internal/cmd/label"
	"github.com/Peripli/service-manager-cli/internal/cmd/login"
	"github.com/Peripli/service-manager-cli/internal/cmd/offering"
	"github.com/Peripli/service-manager-cli/internal/cmd/plan"
	"github.com/Peripli/service-manager-cli/internal/cmd/platform"
	"github.com/Peripli/service-manager-cli/internal/cmd/status"
	"github.com/Peripli/service-manager-cli/internal/cmd/version"
	"github.com/Peripli/service-manager-cli/internal/cmd/visibility"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/auth/oidc"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"os"
)

func oidcAuthBuilder(options *auth.Options) (auth.Authenticator, *auth.Options, error) {
	return oidc.NewOpenIDStrategy(options)
}

func main() {
	cmdContext := &cmd.Context{
		Ctx: context.Background(),
	}
	rootCmd := cmd.BuildRootCommand(cmdContext)
	fs := afero.NewOsFs()

	normalCommandsGroup := cmd.Group{
		Commands: []cmd.CommandPreparator{
			login.NewLoginCmd(cmdContext, os.Stdin, oidcAuthBuilder),
			version.NewVersionCmd(cmdContext),
			info.NewInfoCmd(cmdContext),
		},
		PrepareFn: cmd.CommonPrepare,
	}

	smCommandsGroup := cmd.Group{
		Commands: []cmd.CommandPreparator{
			curl.NewCurlCmd(cmdContext, fs),
			binding.NewListBindingsCmd(cmdContext),
			binding.NewGetBindingCmd(cmdContext),
			binding.NewBindCmd(cmdContext),
			binding.NewUnbindCmd(cmdContext, os.Stdin),
			broker.NewRegisterBrokerCmd(cmdContext),
			broker.NewGetBrokerCmd(cmdContext),
			broker.NewListBrokersCmd(cmdContext),
			broker.NewDeleteBrokerCmd(cmdContext, os.Stdin),
			broker.NewUpdateBrokerCmd(cmdContext),
			platform.NewRegisterPlatformCmd(cmdContext),
			platform.NewListPlatformsCmd(cmdContext),
			platform.NewDeletePlatformCmd(cmdContext, os.Stdin),
			platform.NewUpdatePlatformCmd(cmdContext),
			visibility.NewRegisterVisibilityCmd(cmdContext),
			visibility.NewListVisibilitiesCmd(cmdContext),
			visibility.NewUpdateVisibilityCmd(cmdContext),
			visibility.NewDeleteVisibilityCmd(cmdContext, os.Stdin),
			offering.NewListOfferingsCmd(cmdContext),
			offering.NewMarketplaceCmd(cmdContext),
			plan.NewListPlansCmd(cmdContext),
			label.NewLabelCmd(cmdContext),
			status.NewStatusCmd(cmdContext),
			instance.NewListInstancesCmd(cmdContext),
			instance.NewGetInstanceCmd(cmdContext),
			instance.NewProvisionCmd(cmdContext),
			instance.NewDeprovisionCmd(cmdContext, os.Stdin),
			instance.NewTransferCmd(cmdContext, os.Stdin),
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
