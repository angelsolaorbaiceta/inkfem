/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// '<name>' -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
var sectionDefinitionRegex = regexp.MustCompile(
	"^" + nameGrpExpr + arrowExpr +
		floatGroupExpr("area") + spaceExpr +
		floatGroupExpr("istrong") + spaceExpr +
		floatGroupExpr("iweak") + spaceExpr +
		floatGroupExpr("sstrong") + spaceExpr +
		floatGroupExpr("sweak") + optionalSpaceExpr + "$")

func readSections(scanner *bufio.Scanner, count int) *map[string]*structure.Section {
	lines := definitionLines(scanner, count)
	return deserializeSectionsByName(lines)
}

func deserializeSectionsByName(lines []string) *map[string]*structure.Section {
	var (
		section  *structure.Section
		sections = make(map[string]*structure.Section)
	)

	for _, line := range lines {
		section = deserializeSection(line)
		sections[section.Name] = section
	}

	return &sections
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
