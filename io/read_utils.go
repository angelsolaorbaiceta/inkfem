package io

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const commentDeclaration = "#"

// ShouldIgnoreLine decides whether a given line can be ignored.
// A line can be ignored if, after removing the surrounding white space, it's empty or starts
// with "#", the comment opener.
func ShouldIgnoreLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return lineIsComment(trimmedLine) || lineIsEmpty(trimmedLine)
}

func lineIsComment(line string) bool {
	return strings.HasPrefix(line, commentDeclaration)
}

func lineIsEmpty(line string) bool {
	return len(line) < 1
}

// ExtractDefinitionLines gets the next "count" lines from the scanner, ignoring comments
// and blank lines.
// The extracted lines are trimmed to remove the blank space around them.
func ExtractDefinitionLines(scanner *bufio.Scanner, count int) []string {
	var (
		line  string
		lines = make([]string, count)
	)

	for i := 0; i < count; {
		if !scanner.Scan() {
			panic(fmt.Sprintf("Couldn't read all expected %d lines", count))
		}

		line = strings.TrimSpace(scanner.Text())
		if ShouldIgnoreLine(line) {
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
