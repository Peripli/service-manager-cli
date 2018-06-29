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
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Printer should be implemented for different output formats
type Printer interface {
	// Print is executed with writer and data to be printed
	Print(io.Writer, interface{})
}

// JSONPrinter implements Printer interface and outputs in JSON format
type JSONPrinter struct{}

// Print prints in json format
func (p *JSONPrinter) Print(wr io.Writer, data interface{}) {
	b, err := json.Marshal(data)
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

// YAMLPrinter implements Printer interface and outputs in YAML format
type YAMLPrinter struct{}

// Print prints in yaml format
func (p *YAMLPrinter) Print(wr io.Writer, data interface{}) {
	b, err := yaml.Marshal(data)
	if err != nil {
		PrintError(wr, err)
	} else {
		PrintMessage(wr, string(b))
	}
}
