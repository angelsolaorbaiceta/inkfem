package io

import (
	"regexp"
)

var (
	genericSectionHeaderRegex = regexp.MustCompile(`^\|([\w-_]+)\|(\s*\d+)?$`)
	NodesHeader               = "nodes"
	MaterialsHeader           = "materials"
	SectionsHeader            = "sections"
	LoadsHeader               = "loads"
	BarsHeader                = "bars"
)

func IsSectionHeaderLine(line string) bool {
	return genericSectionHeaderRegex.MatchString(line)
}

func ParseSectionHeader(line string) string {
	return genericSectionHeaderRegex.FindStringSubmatch(line)[1]
}
