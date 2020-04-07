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
	context := &cmd.Context{}
	rootCmd := cmd.BuildRootCommand(context)
	fs := afero.NewOsFs()

	normalCommandsGroup := cmd.Group{
		Commands: []cmd.CommandPreparator{
			login.NewLoginCmd(context, os.Stdin, oidcAuthBuilder),
			version.NewVersionCmd(context),
			info.NewInfoCmd(context),
		},
		PrepareFn: cmd.CommonPrepare,
	}

	smCommandsGroup := cmd.Group{
		Commands: []cmd.CommandPreparator{
			curl.NewCurlCmd(context, fs),
			binding.NewListBindingsCmd(context),
			binding.NewGetBindingCmd(context),
			binding.NewBindCmd(context),
			binding.NewUnbindCmd(context, os.Stdin),
			broker.NewRegisterBrokerCmd(context),
			broker.NewGetBrokerCmd(context),
			broker.NewListBrokersCmd(context),
			broker.NewDeleteBrokerCmd(context, os.Stdin),
			broker.NewUpdateBrokerCmd(context),
			platform.NewRegisterPlatformCmd(context),
			platform.NewListPlatformsCmd(context),
			platform.NewDeletePlatformCmd(context, os.Stdin),
			platform.NewUpdatePlatformCmd(context),
			visibility.NewRegisterVisibilityCmd(context),
			visibility.NewListVisibilitiesCmd(context),
			visibility.NewUpdateVisibilityCmd(context),
			visibility.NewDeleteVisibilityCmd(context, os.Stdin),
			offering.NewListOfferingsCmd(context),
			offering.NewMarketplaceCmd(context),
			plan.NewListPlansCmd(context),
			label.NewLabelCmd(context),
			status.NewStatusCmd(context),
			instance.NewListInstancesCmd(context),
			instance.NewGetInstanceCmd(context),
			instance.NewProvisionCmd(context),
			instance.NewDeprovisionCmd(context, os.Stdin),
			instance.NewTransferCmd(context),
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
