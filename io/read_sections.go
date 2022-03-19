package io

import (
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// '<name>' -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
var sectionDefinitionRegex = regexp.MustCompile(
	"^" + NameGrpExpr + ArrowExpr +
		FloatGroupExpr("area") + SpaceExpr +
		FloatGroupExpr("istrong") + SpaceExpr +
		FloatGroupExpr("iweak") + SpaceExpr +
		FloatGroupExpr("sstrong") + SpaceExpr +
		FloatGroupExpr("sweak") + OptionalSpaceExpr + "$")

func ReadSections(linesReader *LinesReader, count int) map[string]*structure.Section {
	lines := linesReader.GetNextLines(count)
	return deserializeSectionsByName(lines)
}

func deserializeSectionsByName(lines []string) map[string]*structure.Section {
	var (
		section  *structure.Section
		sections = make(map[string]*structure.Section)
	)

	for _, line := range lines {
		section = deserializeSection(line)
		sections[section.Name] = section
	}

	return sections
}

func deserializeSection(definition string) *structure.Section {
	if !sectionDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found section with wrong format: '%s'", definition))
	}

	groups := sectionDefinitionRegex.FindStringSubmatch(definition)

	name := groups[1]
	area := ensureParseFloat(groups[2], "section area")
	iStrong := ensureParseFloat(groups[3], "section iStrong")
	iWeak := ensureParseFloat(groups[4], "section iWeak")
	sStrong := ensureParseFloat(groups[5], "section sStrong")
	sWeak := ensureParseFloat(groups[6], "section sWeak")

	return &structure.Section{
		Name:    name,
		Area:    area,
		IStrong: iStrong,
		IWeak:   iWeak,
		SStrong: sStrong,
		SWeak:   sWeak,
	}
}
