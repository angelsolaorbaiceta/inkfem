package io

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var versionRegex = regexp.MustCompile(`(?:inkfem\s+v)(\d+)(?:[.])(\d+)`)

// ParseVersionNumbers expectes the passed in string to follow the format "inkfem vM.m"
// where "M" and "m" are the major and minor versions of the application.
// It returns these two version numbers or panics if the line couldn't be matched.
func ParseVersionNumbers(versionString string) (majorVersion, minorVersion int) {
	if foundMatch := versionRegex.MatchString(versionString); !foundMatch {
		panic(
			fmt.Sprintf("Could not parse version string '%s' (expected inkfem vM.m)", versionString),
		)
	}

	versions := versionRegex.FindStringSubmatch(versionString)
	majorVersion, _ = strconv.Atoi(versions[1])
	minorVersion, _ = strconv.Atoi(versions[2])

	return
}

// ParseMetadata reads the structure metadata from the structure files first line: "inkfem vM.m".
// Panics if the first line doesn't follow the expected format.
func ParseMetadata(linesReader *LinesReader) structure.StrMetadata {
	// First line must be "inkfem vM.m"
	if !linesReader.ReadNext() {
		panic("The first line should be 'inkfem vM.m'")
	}

	majorVersion, minorVersion := ParseVersionNumbers(linesReader.GetNextLine())

	return structure.StrMetadata{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
	}
}
