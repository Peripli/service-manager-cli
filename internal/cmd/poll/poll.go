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

package poll

import (
	"fmt"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
	"strings"
)

// Cmd wraps smctl poll command
type Cmd struct {
	*cmd.Context

	operationURL string
	outputFormat output.Format
}

// NewPollCmd returns new label command with context
func NewPollCmd(context *cmd.Context) *Cmd {
	return &Cmd{Context: context}
}

// Prepare returns cobra command
func (c *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "poll [operation URL path]",
		Short: "Poll asynchronous operation's status",
		Long:  "Poll asynchronous operation's status",

		PreRunE: prepare(c, c.Context),
		RunE:    cmd.RunE(c),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddCommonQueryFlag(result.Flags(), &c.Parameters)

	return result
}

// Validate validates command's arguments
func (c *Cmd) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("a path to operation is required")
	}
	c.operationURL = args[0]
	return nil
}

// Run runs the command's logic
func (c *Cmd) Run() error {
	operation, err := c.Client.Poll(c.operationURL, &c.Parameters)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode: 404") {
			output.PrintMessage(c.Output, "Operation not found.\n")
			return nil
		}
		return err
	}
	output.PrintServiceManagerObject(c.Output, c.outputFormat, operation)
	output.Println(c.Output)

	return nil
}

// HideUsage hide command's usage
func (c *Cmd) HideUsage() bool {
	return true
}

// SetOutputFormat set output format
func (c *Cmd) SetOutputFormat(format output.Format) {
	c.outputFormat = format
}
