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
	"encoding/json"

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
	if _, err := fmt.Fprintf(wr, "Error: %s\n", err); err != nil {
		panic(err)
	}
}

// PrintMessage prints a message.
func PrintMessage(wr io.Writer, format string, a ...interface{}) {
	if _, err := fmt.Fprintf(wr, format, a...); err != nil {
		panic(err)
	}
}

// Println prints a new line.
func Println(wr io.Writer) {
	if _, err := fmt.Fprintln(wr); err != nil {
		panic(err)
	}
}

// PrintTable prints in table format
func PrintTable(wr io.Writer, data *types.TableData) {
	if _, err := fmt.Fprint(wr, data); err != nil {
		panic(err)
	}
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

type converterFunc func([]byte) (interface{}, error)

// PrintFormat prints the object in the provided format if possible
func PrintFormat(wr io.Writer, outputFormat Format, encodedObject []byte, converter converterFunc) error {
	object, err := converter(encodedObject)
	if err != nil {
		return err
	}
	printer, found := printers[outputFormat]
	if !found {
		PrintMessage(wr, string(encodedObject))
		return nil
	}
	printer.Print(wr, object)
	return nil
}

func PrintParameters(parameters map[string]interface{}) string{
	jsonParameters,_ := json.MarshalIndent(parameters, "", "   ")
	stringParams := string(jsonParameters)
	return stringParams
}