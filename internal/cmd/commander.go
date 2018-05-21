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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Peripli/service-manager-cli/internal/print"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

type CommandWrapper interface {
	Command() *cobra.Command
}

type Command interface {
	Run() error
}

// Implement if you want validation
type ValidatedCommand interface {
	// Validate command usage in the implementation of this method
	Validate([]string) error
}

// Implement if you want to omit usage printing on errors except on validation
type HiddenUsageCommand interface {
	// Return false in the implementation of this method
	HideUsage() bool
}

// Implement if you want to support different output formatting through a --format or -f flag
type FormattedCommand interface {
	// Set the command's `outputFormat` field in the implementation of this method
	SetOutputFormat(int)
}

// Implement if command needs to use the SM Client
type SvcManagerCommand interface {
	// Set the command's `Client` field in the implementation of this method
	SetSMClient(smclient.Client)
}

func PreRunE(cmd Command, ctx *Context) func(*cobra.Command, []string) error {
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

		if svcCmd, ok := cmd.(SvcManagerCommand); ok && ctx.Client == nil {
			clientConfig, err := ctx.Configuration.Load()
			if err != nil {
				return errors.New("no logged user. Use \"smctl login\" to log in")
			}
			svcCmd.SetSMClient(smclient.NewClient(clientConfig))
		}

		return nil

	}
}

func RunE(cmd Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		return cmd.Run()
	}
}

// Adds the --format (-f) flag. Use this method when building a command if it implements the FormattedCommand interface
func AddFormatFlag(flags *pflag.FlagSet) {
	flags.StringP("format", "f", "", "output format")
}

func getOutputFormat(flags *pflag.FlagSet) (int, error) {
	outputFormat, _ := flags.GetString("format")
	outputFormat = strings.ToLower(outputFormat)

	if outputFormat == "" || outputFormat == "text" {
		return print.FormatText, nil
	} else if outputFormat == "json" {
		return print.FormatJSON, nil
	} else if outputFormat == "yaml" {
		return print.FormatYAML, nil
	} else if outputFormat == "raw" {
		return print.FormatRaw, nil
	} else {
		return 0, errors.New("Unknown format: " + outputFormat)
	}
}
