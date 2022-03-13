package io

import (
	"regexp"
	"strconv"
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
