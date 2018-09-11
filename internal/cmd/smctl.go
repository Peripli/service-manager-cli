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

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/configuration"
)

// Execute executes the root command
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// BuildRootCommand builds a new SM root command with context
func BuildRootCommand(ctx *Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "smctl",
		Short: "Service Manager CLI",
		Long:  `smctl controls a Service Manager instance.`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if ctx.Output == nil {
				ctx.Output = cmd.OutOrStdout()
			}

			cmd.SilenceUsage = false
			return nil
		},
	}

	flags := rootCmd.PersistentFlags()
	configuration.AddPFlags(flags)
	environment, err := configuration.NewEnv(flags)
	if err != nil {
		panic(err)
	}
	if ctx.Configuration == nil {
		ctx.Configuration = environment
	}

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sm/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&ctx.Verbose, "verbose", "v", false, "verbose")

	return rootCmd
}
