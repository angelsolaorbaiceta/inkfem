package io

import (
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

// ExtractNamedGroups returns a map of matches by group id.
// Panics if the given string doesn't match the regular expression.
func ExtractNamedGroups(re *regexp.Regexp, str string) map[string]string {
	if !re.MatchString(str) {
		panic(
			fmt.Sprintf("'%s' doesn't match expression: %s", str, re),
		)
	}

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
