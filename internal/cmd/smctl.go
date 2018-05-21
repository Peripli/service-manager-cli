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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func BuildRootCommand(ctx *Context) *cobra.Command {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use:   "smctl",
		Short: "Service Manager CLI",
		Long:  `smctl controls a Service Manager instance.`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if ctx.Output == nil {
				ctx.Output = cmd.OutOrStdout()
			}
			if ctx.Configuration == nil {
				configuration, err := configuration.NewSMConfiguration(cfgFile)
				if err != nil {
					return err
				}
				ctx.Configuration = configuration
			}

			cmd.SilenceUsage = false
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.servicemanager)")
	rootCmd.PersistentFlags().BoolVarP(&ctx.Verbose, "verbose", "v", false, "verbose")

	return rootCmd
}
