package io

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/math"
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

// EnsureParseFloat attempts to parse a floating point number from the given string and panics
// if the operation fails. The "context" is used as part of the panic message and it refers
// to the name of the number being parsed.
func EnsureParseFloat(stringValue string, context string) float64 {
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

// EnsureParseInt attempts to parse a floating point number from the given string and panics
// if the operation fails. The "context" is used as part of the panic message and it refers
// to the name of the number being parsed.
func EnsureParseInt(stringValue string, context string) int {
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

// EnsureParseTorsor attempts to parse a torsor given it's string form: {%f %f %f}.
// Panics if the operation fails. The "context" is used as part of the panic message and it
// refers to the name of the number being parsed.
func EnsureParseTorsor(torsorString string, context string) *math.Torsor {
	nums := strings.Fields(strings.Trim(torsorString, " {}"))
	if len(nums) != 3 {
		panic(
			fmt.Sprintf("Error reading %s: can't parse torsor from %s", context, torsorString),
		)
	}

	return math.MakeTorsor(
		EnsureParseFloat(nums[0], context+" (Fx)"),
		EnsureParseFloat(nums[1], context+" (Fy)"),
		EnsureParseFloat(nums[2], context+" (Mz)"),
	)
}

// EnsureParseDOF attempts to parse three degrees of freedom given the format: [%d %d %d].
// Panics if the operation fails. The "context" is used as part of the panic message and it
// refers to the name of the number being parsed.
func EnsureParseDOF(dofString string, context string) (int, int, int) {
	dofs := strings.Fields(strings.Trim(dofString, " []"))
	if len(dofs) != 3 {
		panic(
			fmt.Sprintf("Error reading %s: can't parse DOF from %s", context, dofString),
		)
	}

	var (
		dof1 = EnsureParseInt(dofs[0], context+" (dx DOF)")
		dof2 = EnsureParseInt(dofs[1], context+" (dy DOF)")
		dof3 = EnsureParseInt(dofs[2], context+" (rz DOF)")
	)

	return dof1, dof2, dof3
}
