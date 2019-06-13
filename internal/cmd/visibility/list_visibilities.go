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

package visibility

import (
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/spf13/cobra"
)

// ListVisibilitiesCmd wraps the smctl list-visibilities command
type ListVisibilitiesCmd struct {
	*cmd.Context

	outputFormat output.Format
}

// NewListVisibilitiesCmd returns new list-visibilities command with context
func NewListVisibilitiesCmd(context *cmd.Context) *ListVisibilitiesCmd {
	return &ListVisibilitiesCmd{Context: context}
}

//Run runs the command's logic
func (lv *ListVisibilitiesCmd) Run() error {
	visibilities, err := lv.Client.ListVisibilitiesWithQuery(lv.Parameters.Copy())
	if err != nil {
		return err
	}

	output.PrintServiceManagerObject(lv.Output, lv.outputFormat, visibilities)
	output.Println(lv.Output)
	return nil
}

// SetOutputFormat sets output format
func (lv *ListVisibilitiesCmd) SetOutputFormat(format output.Format) {
	lv.outputFormat = format
}

// HideUsage hides command's usage
func (lv *ListVisibilitiesCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (lv *ListVisibilitiesCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "list-visibilities",
		Aliases: []string{"lv"},
		Short:   "List visibilities",
		Long:    "List all visibilities.",
		PreRunE: prepare(lv, lv.Context),
		RunE:    cmd.RunE(lv),
	}

	cmd.AddFormatFlag(result.Flags())
	cmd.AddQueryingFlags(result.Flags(), lv.Parameters)

	return result
}
