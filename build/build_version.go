package build

import (
	_ "embed"
	"regexp"
	"strconv"
)

//go:embed VERSION
var versionString string

var versionRegex = regexp.MustCompile(`v(\d+)(?:[.])(\d+)`)

type BuildInfo struct {
	MajorVersion, MinorVersion int
}

// Info is the binary's build information that's globally available to the application.
// This is set by the ReadBuildInfo function and will be nil before the function is called.
var Info *BuildInfo = nil

// ReadBuildInfo reads the "VERSION" file and parses the version numbers from it.
// This function should be called when the binary is executed to make the build information
// available to the rest of the application.
func ReadBuildInfo() {
	if Info != nil {
		return
	}

	if matches := versionRegex.MatchString(versionString); !matches {
		panic(
			"Could not parse version string '" + versionString + "' from the VERSION file (expected vM.m)",
		)
	}

	versions := versionRegex.FindStringSubmatch(versionString)
	majorVersion, err := strconv.Atoi(versions[1])
	if err != nil {
		panic("Could not parse major version number '" + versions[1] + "' from the VERSION file")
	}
	minorVersion, err := strconv.Atoi(versions[2])
	if err != nil {
		panic("Could not parse minor version number '" + versions[2] + "' from the VERSION file")
	}

	Info = &BuildInfo{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
	}
}
