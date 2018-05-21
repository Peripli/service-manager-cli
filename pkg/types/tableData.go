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

type TableData struct {
	Headers []string
	Data    [][]string
}

func (data *TableData) String() string {
	output := ""
	if len(data.Data) == 0 {
		return output
	}

	// get fields lengths
	fieldLen := make([]int, len(data.Headers))
	for i, header := range data.Headers {
		if fieldLen[i] < len(header)+2 {
			fieldLen[i] = len(header) + 2
		}
	}

	for _, row := range data.Data {
		for i, cell := range row {
			if fieldLen[i] < len(cell)+2 {
				fieldLen[i] = len(cell) + 2
			}
		}
	}

	for i, header := range data.Headers {
		output += pad(header, fieldLen[i])
	}
	output += "\n"

	for i := range data.Headers {
		output += line(fieldLen[i]-2) + "  "
	}
	output += "\n"

	for _, row := range data.Data {
		for i, cell := range row {
			output += pad(cell, fieldLen[i])
		}
		output += "\n"
	}

	return output
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
