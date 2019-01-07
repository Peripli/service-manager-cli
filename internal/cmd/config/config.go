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

package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
)

// Cmd wraps the smctl curl command
type Cmd struct {
	*cmd.Context

	property string
	value    string
}

// NewConfigCmd returns new curl command with context
func NewConfigCmd(context *cmd.Context) *Cmd {
	return &Cmd{Context: context}
}

// Prepare returns the cobra command
func (c *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "config [property] [value]",
		Aliases: []string{"config"},
		Short:   "set SM config property to a certain value",
		Long:    `set SM config property to a certain value`,

		PreRunE: prepare(c, c.Context),
		RunE:    cmd.RunE(c),
	}

	return result
}

// Validate validates command's arguments
func (c *Cmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[property] is required")
	}
	c.property = args[0]

	if len(args) == 2 && len(args[1]) > 0 {
		c.value = args[1]
	}

	return nil
}

// Run runs the command's logic
func (c *Cmd) Run() error {
	if c.value == "" {
		val, err := c.Context.Configuration.Get(c.property)
		if err != nil {
			return err
		}
		if val != nil {
			output.PrintMessage(c.Output, "%v\n", val)
		}
	} else {
		return c.Context.Configuration.Set(c.property, c.value)
	}

	return nil
}

// HideUsage hide command's usage
func (c *Cmd) HideUsage() bool {
	return true
}
