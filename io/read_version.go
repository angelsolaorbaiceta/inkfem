package io

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var versionRegex = regexp.MustCompile(`(?:inkfem\s+v)(\d+)(?:[.])(\d+)`)
var binaryMajorVersion int = -1
var binaryMinorVersion int = -1

// SetBinaryVersion saves the binary version numbers.
// This function should be called when the binary is executed to make the version numbers available
// to the rest of the application.
func SetBinaryVersion(verstionString string) {
	majorVersion, minorVersion := ParseVersionNumbers(verstionString)
	binaryMajorVersion = majorVersion
	binaryMinorVersion = minorVersion
}

// GetBinaryVersion returns the binary version numbers.
// This numbers should have been set by the SetBinaryVersion function. Panics if the version
// numbers weren't set.
func GetBinaryVersion() (majorVersion, minorVersion int) {
	if binaryMajorVersion == -1 || binaryMinorVersion == -1 {
		panic("inkfem binary version not set")
	}

	return binaryMajorVersion, binaryMinorVersion
}

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
