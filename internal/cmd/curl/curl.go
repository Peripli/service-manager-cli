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

package curl

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// Cmd wraps the smctl curl command
type Cmd struct {
	*cmd.Context

	path   string
	method string
	body   string
}

// NewCurlCmd returns new curl command with context
func NewCurlCmd(context *cmd.Context) *Cmd {
	return &Cmd{Context: context}
}

// Prepare returns the cobra command
func (c *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "curl [path]",
		Aliases: []string{"c"},
		Short:   "call arbitrary SM endpoint",
		Long:    `call arbitrary SM endpoint`,

		PreRunE: prepare(c, c.Context),
		RunE:    cmd.RunE(c),
	}

	// TODO: Add flag for headers
	result.Flags().StringVarP(&c.method, "X", "X", "GET", "HTTP method (GET,POST,PUT,DELETE,etc)")
	result.Flags().StringVarP(&c.body, "d", "d", "", "HTTP data to include in the request body, or '@' followed by a file name to read the data from")

	return result
}

// Validate validates command's arguments
func (c *Cmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[path] is required")
	}
	c.path = args[0]

	return nil
}

// Run runs the command's logic
func (c *Cmd) Run() error {
	var reader io.Reader
	if c.method != http.MethodGet {
		reader = bytes.NewReader([]byte(c.body))
	}

	resp, err := c.Client.Call(c.method, c.path, reader)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	output.PrintMessage(c.Output, string(data))

	return nil
}
