/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const commentDeclaration = "#"

func lineIsComment(line string) bool {
	return strings.HasPrefix(line, commentDeclaration)
}

func lineIsEmpty(line string) bool {
	return len(line) < 1
}

func definitionLines(scanner *bufio.Scanner, count int) []string {
	var (
		line  string
		lines = make([]string, count)
	)

	for i := 0; i < count; {
		if !scanner.Scan() {
			panic("Couldn't read all expected lines")
		}

		line = scanner.Text()
		if lineIsComment(line) {
			continue
		}

		lines[i] = line
		i++
	}

	return lines
}

func ensureParseFloat(stringValue string, context string) float64 {
	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		panic(
			fmt.Sprintf(
				"Error reading %s: can't parse floating point number from %s",
				context,
				stringValue,
			),
		)
	}

	return value
}

func ensureParseInt(stringValue string, context string) int {
	value, err := strconv.Atoi(stringValue)
	if err != nil {
		panic(
			fmt.Sprintf(
				"Error reading %s: can't parse integr number from %s",
				context,
				stringValue,
			),
		)
	}

	return value
}
