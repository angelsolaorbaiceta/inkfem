package io

import (
	"bufio"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var versionRegex = regexp.MustCompile(`(?:inkfem\s+v)(\d+)(?:[.])(\d+)`)

// parseVersionNumbers expectes the passed in first line to follow the format "inkfem vM.m"
// where "M" and "m" are the major and minor versions of the application that created
// the parsed file. It returns these two version numbers or panics if the line couldn't
// be matched.
func ParseVersionNumbers(firstLine string) (majorVersion, minorVersion int) {
	if foundMatch := versionRegex.MatchString(firstLine); !foundMatch {
		panic(
			"Could not parse major and minor version numbers." +
				"Are you missing 'inkfem vM.m' in your file's first line?",
		)
	}

	versions := versionRegex.FindStringSubmatch(firstLine)
	majorVersion, _ = strconv.Atoi(versions[1])
	minorVersion, _ = strconv.Atoi(versions[2])

	return
}

// ParseMetadata reads the structure metadata from the structure files first line: "inkfem vM.m".
// Panics if the first line doesn't follow the expected format.
func ParseMetadata(scanner *bufio.Scanner) structure.StrMetadata {
	// First line must be "inkfem vM.m"
	scanner.Scan()
	majorVersion, minorVersion := ParseVersionNumbers(scanner.Text())
	return structure.StrMetadata{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
	}
}
