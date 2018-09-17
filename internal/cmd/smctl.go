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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/pkg/httputil"
)

// Execute executes the root command
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// BuildRootCommand builds a new SM root command with context
func BuildRootCommand(ctx *Context) *cobra.Command {
	var cfgFile string
	httpConfig := httputil.DefaultHTTPConfig()

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
				config, err := configuration.New(viper.New(), cfgFile)
				if err != nil {
					return fmt.Errorf("Could not create configuration: %s", err)
				}
				ctx.Configuration = config
				config.Set(configuration.HTTPConfigKey, httpConfig)
			}

			cmd.SilenceUsage = false
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sm/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&ctx.Verbose, "verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().DurationVarP(&httpConfig.Timeout, "timeout", "t", httpConfig.Timeout, "Timeout for http requests")
	rootCmd.PersistentFlags().DurationVarP(&httpConfig.KeepAlive, "keepalive", "k", httpConfig.KeepAlive, "Timeout for keepalive http connections")
	rootCmd.PersistentFlags().BoolVar(&httpConfig.SSLDisabled, "skip-ssl-validation", httpConfig.SSLDisabled, "Skip ssl validation for http requests")

	return rootCmd
}
