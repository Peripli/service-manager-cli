package label

import (
	"fmt"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/Peripli/service-manager/pkg/query"
	"github.com/Peripli/service-manager/pkg/web"
	"github.com/spf13/cobra"
)

// Cmd wraps smctl label command
type Cmd struct {
	*cmd.Context

	resourcePath string
	id           string
	values       []string
	labelChanges types.LabelChanges
}

// NewLabelCmd returns new label command with context
func NewLabelCmd(context *cmd.Context) *Cmd {
	return &Cmd{Context: context, labelChanges: types.LabelChanges{}}
}

// Prepare returns cobra command
func (c *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "label [resource] [id] [operation] [key] [--val value1 --val value2 ...]",
		Short: "Label resource",
		Long:  "Label resource",

		PreRunE: prepare(c, c.Context),
		RunE:    cmd.RunE(c),
	}

	result.Flags().StringArrayVar(&c.values, "val", []string{}, "Label value to be used")
	cmd.AddCommonQueryFlag(result.Flags(), &c.Parameters)

	return result
}

// Validate validates command's arguments
func (c *Cmd) Validate(args []string) error {

	resources := map[string]string{
		"platform": web.PlatformsURL,
		"broker":   web.ServiceBrokersURL,
	}

	operations := map[string]string{
		"add":           "add",
		"remove":        "remove",
		"add-values":    "add_values",
		"remove-values": "remove_values",
	}

	if len(args) < 4 {
		return fmt.Errorf("resource type, id, operation, key and values are required")
	}

	if len(args) > 4 {
		return fmt.Errorf("too much arguments, in case you have whitespaces in some of the arguments consider enclosig it with single quotes")
	}

	if v, ok := resources[args[0]]; ok {
		c.resourcePath = v
		c.id = args[1]
	} else {
		return fmt.Errorf("unknown resource")
	}

	labelChange := &query.LabelChange{}

	if v, ok := operations[args[2]]; ok {
		labelChange.Operation = query.LabelOperation(v)
	} else {
		return fmt.Errorf("unknown operation")
	}

	labelChange.Key = args[3]

	if len(c.values) < 1 {
		return fmt.Errorf("at least one value is required")
	}

	labelChange.Values = c.values
	c.labelChanges.LabelChanges = append(c.labelChanges.LabelChanges, labelChange)
	return nil
}

// Run runs the command's logic
func (c *Cmd) Run() error {
	if err := c.Client.Label(c.resourcePath, c.id, &c.labelChanges, &c.Parameters); err != nil {
		return err
	}

	output.PrintMessage(c.Output, "Resource labeled successfully!")
	return nil
}

// HideUsage hide command's usage
func (c *Cmd) HideUsage() bool {
	return true
}
