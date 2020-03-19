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

package types

import (
	"fmt"
	"github.com/Peripli/service-manager/pkg/types"
	"strings"
)

// TableData holds data for table header and content
type TableData struct {
	Headers  []string
	Data     [][]string
	Vertical bool
}

// String implements Stringer interface
func (table *TableData) String() string {
	if table.Vertical {
		return table.verticalTable()
	}
	return table.horizontalTable()
}

func (table *TableData) horizontalTable() string {
	output := ""
	if len(table.Data) == 0 {
		return output
	}

	// get fields lengths
	fieldLen := table.fieldsLen()

	for i, header := range table.Headers {
		output += pad(header, fieldLen[i])
	}
	output += "\n"

	for i := range table.Headers {
		output += line(fieldLen[i]-2) + "  "
	}
	output += "\n"

	for _, row := range table.Data {
		for i, cell := range row {
			output += pad(cell, fieldLen[i])
		}
		output += "\n"
	}

	return output
}

func (table *TableData) verticalTable() string {
	output := ""
	if len(table.Data) == 0 {
		return output
	}

	headerLen := table.headerLen()
	maxHeaderLen := max(headerLen...)

	dataLen := table.dataLen()
	maxDataLen := max(dataLen...)

	for i, header := range table.Headers {
		output += "| "
		output += pad(header, maxHeaderLen)
		output += "| "
		for _, row := range table.Data {
			output += pad(row[i], maxDataLen)
			output += "| "
		}
		output += "\n"
	}

	return output
}

func (table *TableData) headerLen() []int {
	fieldLen := make([]int, len(table.Headers))
	for i, header := range table.Headers {
		if fieldLen[i] < len(header)+2 {
			fieldLen[i] = len(header) + 2
		}
	}
	return fieldLen
}

func (table *TableData) dataLen() []int {
	fieldLen := make([]int, len(table.Headers))
	for _, row := range table.Data {
		for i, cell := range row {
			if fieldLen[i] < len(cell)+2 {
				fieldLen[i] = len(cell) + 2
			}
		}
	}

	return fieldLen
}

func (table *TableData) fieldsLen() []int {
	fieldLen := make([]int, len(table.Headers))
	headerLen := table.headerLen()
	dataLen := table.dataLen()

	for i := range fieldLen {
		fieldLen[i] = max(headerLen[i], dataLen[i])
	}
	return fieldLen
}

func pad(s string, p int) string {
	result := s

	for len(result) < p {
		result += " "
	}

	return result
}

func line(p int) string {
	result := ""

	for len(result) < p {
		result += "-"
	}

	return result
}

func max(arr ...int) int {
	tmp := arr[0]
	for _, i := range arr {
		if i > tmp {
			tmp = i
		}
	}
	return tmp
}

func formatLabels(labels types.Labels) string {
	formattedLabels := make([]string, 0, len(labels))
	for i, v := range labels {
		formattedLabels = append(formattedLabels, i+"="+strings.Join(v, ","))
	}
	return strings.Join(formattedLabels, " ")
}

func formatLastOp(operation *types.Operation) string {
	if operation.State != types.FAILED {
		return fmt.Sprintf("%s %s", operation.Type, operation.State)
	}
	return fmt.Sprintf("%s %s %s", operation.Type, operation.State, operation.Errors)
}
