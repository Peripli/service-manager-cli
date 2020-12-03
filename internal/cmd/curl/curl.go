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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/Peripli/service-manager-cli/internal/output"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// Cmd wraps the smctl curl command
type Cmd struct {
	*cmd.Context
	Fs afero.Fs

	outputFormat output.Format

	path   string
	method string
	body   string
}

// NewCurlCmd returns new curl command with context
func NewCurlCmd(context *cmd.Context, fs afero.Fs) *Cmd {
	return &Cmd{Context: context, Fs: fs}
}

// Prepare returns the cobra command
func (c *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "curl [path]",
		Aliases: []string{"c"},
		Short:   "Call arbitrary SM endpoint",
		Long:    `Call arbitrary SM endpoint`,

		PreRunE: prepare(c, c.Context),
		RunE:    cmd.RunE(c),
	}

	// TODO: Add flag for headers
	result.Flags().StringVarP(&c.method, "X", "X", "GET", "HTTP method (GET,POST,PUT,DELETE,etc)")
	result.Flags().StringVarP(&c.body, "d", "d", "", "HTTP data to include in the request body, or '@' followed by a file name to read the data from")
	cmd.AddFormatFlagDefault(result.Flags(), "json")
	cmd.AddCommonQueryFlag(result.Flags(), &c.Parameters)

	return result
}

// Validate validates command's arguments
func (c *Cmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[path] is required")
	}
	c.path = args[0]

	if strings.HasPrefix(c.body, "@") {
		filename, err := filepath.Abs(c.body[1:])
		if err != nil {
			return err
		}

		f, err := c.Fs.Open(filename)
		if err != nil {
			return err
		}
		fileContents, err := ioutil.ReadAll(f)
		if err != nil {
			return fmt.Errorf("error reading file %s. Reason: %v", filename, err)
		}
		c.body = string(fileContents)
	}

	return nil
}

// Run runs the command's logic
func (c *Cmd) Run() error {
	var reader io.Reader
	if c.method != http.MethodGet {
		reader = bytes.NewReader([]byte(c.body))
	}

	resp, err := c.Client.Call(c.method, c.path, reader, &c.Parameters)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return output.PrintFormat(c.Output, c.outputFormat, data, toMap)
	}
	output.PrintMessage(c.Output, string(data))

	return nil
}

func toMap(data []byte) (interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SetOutputFormat set output format
func (c *Cmd) SetOutputFormat(format output.Format) {
	c.outputFormat = format
}

// HideUsage hide command's usage
func (c *Cmd) HideUsage() bool {
	return true
}
