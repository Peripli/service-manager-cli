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
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Peripli/service-manager-cli/internal/configuration"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/auth/oidc"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

var supportedFormats = map[string]output.Format{
	"text": output.FormatText,
	"json": output.FormatJSON,
	"yaml": output.FormatYAML,
}

// CommandPreparator used to wrap CLI commands
type CommandPreparator interface {
	Prepare(PrepareFunc) *cobra.Command
}

// Command provides common logic for SM commands
type Command interface {
	Run() error
}

// ValidatedCommand should be implemented if the command will have validation
type ValidatedCommand interface {
	// Validate command usage in the implementation of this method
	Validate([]string) error
}

// HiddenUsageCommand should be implemented if the command should not print its usage
type HiddenUsageCommand interface {
	// HideUsage returns true when usage should be hidden and false otherwise
	HideUsage() bool
}

// FormattedCommand should be implemented if the command supports different output formatting through a --format or -f flag
type FormattedCommand interface {
	// SetOutputFormat sets the command's output format
	SetOutputFormat(output.Format)
}

// PrepareFunc is function type which executes common prepare logic for commands
type PrepareFunc func(cmd Command, ctx *Context) func(*cobra.Command, []string) error

// SmPrepare creates a SM client for SM commands
func SmPrepare(cmd Command, ctx *Context) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		if err := CommonPrepare(cmd, ctx)(c, args); err != nil {
			return err
		}

		if ctx.Client == nil {
			settings := configuration.DefaultSettings()
			// TODO: check err
			ctx.Configuration.Unmarshal(settings)
			clientConfig := settings.SMClient

			// if err != nil {
			// 	return fmt.Errorf("no logged user. Use \"smctl login\" to log in. Reason: %s", err)
			// }

			oidcClient := oidc.NewClient(&auth.Options{
				AuthorizationEndpoint: clientConfig.AuthorizationEndpoint,
				TokenEndpoint:         clientConfig.TokenEndpoint,
				ClientID:              clientConfig.ClientID,
				ClientSecret:          clientConfig.ClientSecret,
				IssuerURL:             clientConfig.IssuerURL,
				HTTP:                  settings.HTTP,
			}, &clientConfig.Token)

			refresher, isRefresher := oidcClient.(auth.Refresher)
			if isRefresher {
				token, err := refresher.Token()
				if err != nil {
					return fmt.Errorf("Error refreshing token. Reason: %s", err)
				}
				if clientConfig.AccessToken != token.AccessToken {
					clientConfig.Token = *token
					if saveErr := ctx.Configuration.Save(settings); saveErr != nil {
						return fmt.Errorf("Error saving config file. Reason: %s", saveErr)
					}
				}
			}

			ctx.Client = smclient.NewClient(oidcClient, clientConfig)
		}

		return nil
	}
}

// CommonPrepare provides common pre-run logic for SM commands
func CommonPrepare(cmd Command, ctx *Context) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		if valCmd, ok := cmd.(ValidatedCommand); ok {
			if err := valCmd.Validate(args); err != nil {
				return err
			}
		}

		if fmtCmd, ok := cmd.(FormattedCommand); ok {
			outputFormat, err := getOutputFormat(c.Flags())
			if err != nil {
				return err
			}
			fmtCmd.SetOutputFormat(outputFormat)
		}

		if huCmd, ok := cmd.(HiddenUsageCommand); ok {
			c.SilenceUsage = huCmd.HideUsage()
		}

		return nil
	}
}

// RunE provides common run logic for SM commands
func RunE(cmd Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		return cmd.Run()
	}
}

// AddFormatFlag adds the --format (-f) flag.
func AddFormatFlag(flags *pflag.FlagSet) {
	flags.StringP("output", "o", "", "output format")
}

func getOutputFormat(flags *pflag.FlagSet) (output.Format, error) {
	outputFormat, _ := flags.GetString("output")
	outputFormat = strings.ToLower(outputFormat)

	if outputFormat == "" {
		return output.FormatText, nil
	}
	format, exists := supportedFormats[outputFormat]
	if !exists {
		return output.FormatUnknown, errors.New("unknown output: " + outputFormat)
	}
	return format, nil
}
