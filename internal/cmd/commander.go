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

	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/smclient"
)

// CommandWrapper used to wrap CLI commands
type CommandWrapper interface {
	Command() *cobra.Command
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
	SetOutputFormat(int)
}

// SvcManagerCommand should be implemented if the command needs to use the SM Client
type SvcManagerCommand interface {
	// SetSMClient sets the command's Service Manager client
	SetSMClient(smclient.Client)
}

// PreRunE provides common pre-run logic for SM commands
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

// RunE provides common run logic for SM commands
func RunE(cmd Command) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		return cmd.Run()
	}
}

// AddFormatFlag adds the --format (-f) flag.
func AddFormatFlag(flags *pflag.FlagSet) {
	flags.StringP("format", "f", "", "output format")
}

func getOutputFormat(flags *pflag.FlagSet) (int, error) {
	outputFormat, _ := flags.GetString("format")
	outputFormat = strings.ToLower(outputFormat)

	if outputFormat == "" || outputFormat == "text" {
		return output.FormatText, nil
	} else if outputFormat == "json" {
		return output.FormatJSON, nil
	} else if outputFormat == "yaml" {
		return output.FormatYAML, nil
	} else if outputFormat == "raw" {
		return output.FormatRaw, nil
	} else {
		return 0, errors.New("unknown format: " + outputFormat)
	}
}
