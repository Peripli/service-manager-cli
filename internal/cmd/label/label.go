package label

import (
	"fmt"
	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/internal/output"
	resperror "github.com/Peripli/service-manager-cli/pkg/errors"
	"github.com/Peripli/service-manager-cli/pkg/types"
	"github.com/Peripli/service-manager/pkg/query"
	"github.com/spf13/cobra"
	"strings"
)

// Cmd wraps smctl label command
type Cmd struct {
	*cmd.Context

	resource     string
	id           string
	labelChanges types.LabelChanges
}

// NewLabelCmd returns new label command with context
func NewLabelCmd(context *cmd.Context) *Cmd {
	return &Cmd{Context: context, labelChanges: types.LabelChanges{}}
}

// Prepare returns cobra command
func (c *Cmd) Prepare(prepare cmd.PrepareFunc) *cobra.Command {
	result := &cobra.Command{
		Use:   "label [resource] [id] [operation] [key=value1,value2,...]",
		Short: "Label resource",
		Long:  "Label resource",

		PreRunE: prepare(c, c.Context),
		RunE:    cmd.RunE(c),
	}

	return result
}

// Validate validates command's arguments
func (c *Cmd) Validate(args []string) error {

	resources := map[string]string{
		"platform": "platforms",
		"broker":   "brokers",
	}

	operations := map[string]string{
		"add":           "add",
		"remove":        "remove",
		"add-values":    "add_values",
		"remove-values": "remove_values",
	}

	if len(args) < 4 {
		return fmt.Errorf("resource type, name, operation and value are required")
	}

	if v, ok := resources[args[0]]; ok {
		c.resource = v
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
	keyValueSeparatorIndex := strings.Index(args[3], "=")
	labelChange.Key = args[3][:keyValueSeparatorIndex]
	values := args[3][keyValueSeparatorIndex+1:]
	labelChange.Values = strings.Split(values, ",")
	c.labelChanges.LabelChanges = append(c.labelChanges.LabelChanges, labelChange)
	return nil
}

// Run runs the command's logic
func (c *Cmd) Run() error {
	err := c.Client.Label(c.resource, c.id, &c.labelChanges)
	if responseErr, ok := err.(resperror.ResponseError); ok {
		return fmt.Errorf(responseErr.Description)
	} else if err != nil {
		return err
	}

	output.PrintMessage(c.Output, "Resource labeled successfully!")
	return nil
}

// HideUsage hide command's usage
func (c *Cmd) HideUsage() bool {
	return true
}
