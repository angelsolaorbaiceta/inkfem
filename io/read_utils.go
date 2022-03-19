package io

import (
	"bufio"
	"fmt"
	"regexp"
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

// ExtractNamedGroups returns a map of matches by group id.
func ExtractNamedGroups(re *regexp.Regexp, str string) map[string]string {
	var (
		matches = re.FindStringSubmatch(str)
		result  = make(map[string]string)
	)

	for i, name := range re.SubexpNames() {
		if i == 0 || len(matches[i]) == 0 {
			continue
		}

		result[name] = matches[i]
	}

	return result
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
				"Error reading %s: can't parse integer number from %s",
				context,
				stringValue,
			),
		)
	}

	return value
}
