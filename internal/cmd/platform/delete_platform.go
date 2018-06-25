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

package platform

import (
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/internal/util"

	"github.com/spf13/cobra"

	"github.com/Peripli/service-manager-cli/internal/cmd"
)

// DeletePlatformCmd wraps the smctl list-brokers command
type DeletePlatformCmd struct {
	*cmd.Context

	names []string
}

// NewDeletePlatformCmd returns new list-brokers command with context
func NewDeletePlatformCmd(context *cmd.Context) *DeletePlatformCmd {
	return &DeletePlatformCmd{Context: context}
}

// Validate validates command's arguments
func (dpc *DeletePlatformCmd) Validate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("[name] is required")
	}

	dpc.names = args

	return nil
}

// Run runs the command's logic
func (dpc *DeletePlatformCmd) Run() error {
	allPlatforms, err := dpc.Client.ListPlatforms()
	if err != nil {
		return err
	}

	toDeletePlatforms := util.GetPlatformsByName(allPlatforms, dpc.names)
	if len(toDeletePlatforms) < 1 {
		output.PrintMessage(dpc.Output, "Platform(s) not found\n")
		return nil
	}

	deletedPlatforms := make(map[string]bool)

	for _, toDelete := range toDeletePlatforms {
		err := dpc.Client.DeletePlatform(toDelete.ID)
		if err != nil {
			output.PrintMessage(dpc.Output, "Could not delete platform %s\n", toDelete.Name, err)
		} else {
			output.PrintMessage(dpc.Output, "Platform with name: %s successfully deleted\n", toDelete.Name)
			deletedPlatforms[toDelete.Name] = true
		}
	}

	for _, platformName := range dpc.names {
		if _, deleted := deletedPlatforms[platformName]; !deleted {
			output.PrintError(dpc.Output, fmt.Errorf("Platform with name: %s was not found", platformName))
		}
	}

	return nil
}

// HideUsage hide command's usage
func (dpc *DeletePlatformCmd) HideUsage() bool {
	return true
}

// Prepare returns cobra command
func (dpc *DeletePlatformCmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:     "delete-platform [name] <name2 <name3> ... <nameN>>",
		Aliases: []string{"dp"},
		Short:   "Deletes platforms",
		Long:    `Delete one or more platforms with name.`,
		PreRunE: prepare(dpc, dpc.Context),
		RunE:    cmd.RunE(dpc),
	}

	return result
}
