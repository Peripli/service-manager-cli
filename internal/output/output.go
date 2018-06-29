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

package output

import (
	"fmt"
	"io"

	"github.com/Peripli/service-manager-cli/pkg/types"
)

// Format is predefined type for output format
type Format int

const (
	// FormatText const for text output format
	FormatText = iota
	// FormatJSON const for json output format
	FormatJSON
	// FormatYAML const for yaml output format
	FormatYAML
	// FormatUnknown const for unknown output format
	FormatUnknown
)

var (
	printers = map[Format]Printer{
		FormatJSON: &JSONPrinter{},
		FormatYAML: &YAMLPrinter{},
	}
)

// PrintError prints an error.
func PrintError(wr io.Writer, err error) {
	fmt.Fprintf(wr, "Error: %s\n", err)
}

// PrintMessage prints a message.
func PrintMessage(wr io.Writer, format string, a ...interface{}) {
	fmt.Fprintf(wr, format, a...)
}

// Println prints a new line.
func Println(wr io.Writer) {
	fmt.Fprintln(wr)
}

// PrintTable prints in table format
func PrintTable(wr io.Writer, data *types.TableData) {
	fmt.Fprint(wr, data)
}

// PrintServiceManagerObject should be used for printing SM objects in different formats
func PrintServiceManagerObject(wr io.Writer, outputFormat Format, object types.ServiceManagerObject) {
	tableDataPrinter, isTableDataPrinter := object.(types.TableDataPrinter)
	if outputFormat == FormatText && isTableDataPrinter {
		PrintMessage(wr, object.Message())
		Println(wr)
		if !object.IsEmpty() {
			PrintTable(wr, tableDataPrinter.TableData())
			Println(wr)
		}
	} else {
		printers[outputFormat].Print(wr, object)
	}
}
