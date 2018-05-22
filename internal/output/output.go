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
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/Peripli/service-manager-cli/pkg/types"
)

const (
	// FormatText const for text format
	FormatText = iota
	// FormatJSON const for json format
	FormatJSON
	// FormatYAML const for yaml format
	FormatYAML
	// FormatRaw const for raw format
	FormatRaw
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

// PrintJSON prints in json format
func PrintJSON(wr io.Writer, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		PrintError(wr, err)
	} else {
		var out bytes.Buffer
		err := json.Indent(&out, b, "", "  ")
		if err != nil {
			PrintError(wr, err)
		}
		PrintMessage(wr, out.String())
	}
}

// PrintYAML prints in yaml format
func PrintYAML(wr io.Writer, v interface{}) {
	b, err := yaml.Marshal(v)
	if err != nil {
		PrintError(wr, err)
	} else {
		PrintMessage(wr, string(b))
	}
}

// PrintTable prints in table format
func PrintTable(wr io.Writer, data *types.TableData) {
	fmt.Fprint(wr, data)
}

// PrintServiceManagerObject should be used for printing SM objects in different formats
func PrintServiceManagerObject(wr io.Writer, outputFormat int, object types.ServiceManagerObject) {
	tableDataPrinter, isTableDataPrinter := object.(types.TableDataPrinter)
	if isTableDataPrinter && outputFormat == FormatText {
		PrintMessage(wr, object.Message())
		Println(wr)
		if !object.IsEmpty() {
			PrintTable(wr, tableDataPrinter.TableData())
			Println(wr)
		}
	} else if outputFormat == FormatJSON {
		PrintJSON(wr, object)
	} else if outputFormat == FormatYAML {
		PrintYAML(wr, object)
	}
}
