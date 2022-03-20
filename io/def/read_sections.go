package def

import (
	"fmt"
	"regexp"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// '<name>' -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
var sectionDefinitionRegex = regexp.MustCompile(
	"^" + inkio.NameGrpExpr + inkio.ArrowExpr +
		inkio.FloatGroupExpr("area") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("istrong") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("iweak") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("sstrong") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("sweak") + inkio.OptionalSpaceExpr + "$")

func ReadSections(linesReader *inkio.LinesReader, count int) map[string]*structure.Section {
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
	area := inkio.EnsureParseFloat(groups[2], "section area")
	iStrong := inkio.EnsureParseFloat(groups[3], "section iStrong")
	iWeak := inkio.EnsureParseFloat(groups[4], "section iWeak")
	sStrong := inkio.EnsureParseFloat(groups[5], "section sStrong")
	sWeak := inkio.EnsureParseFloat(groups[6], "section sWeak")

	return &structure.Section{
		Name:    name,
		Area:    area,
		IStrong: iStrong,
		IWeak:   iWeak,
		SStrong: sStrong,
		SWeak:   sWeak,
	}
}
