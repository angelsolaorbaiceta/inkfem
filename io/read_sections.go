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
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// <name> -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
var sectionDefinitionRegex = regexp.MustCompile(
	`(?P<name>'\w+')(?:\s*->\s*)` +
		`(?P<area>\d+\.*\d*)(?:\s+)` +
		`(?P<istrong>\d+\.+\d+)(?:\s+)` +
		`(?P<iweak>\d+\.+\d+)(?:\s+)` +
		`(?P<sstrong>\d+\.+\d+)(?:\s+)` +
		`(?P<sweak>\d+\.+\d+)`)

func readSections(scanner *bufio.Scanner, count int) *map[string]*structure.Section {
	var (
		name                                 string
		area, iStrong, iWeak, sStrong, sWeak float64
		sections                             = make(map[string]*structure.Section)
	)

	for _, line := range definitionLines(scanner, count) {
		if !sectionDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found section with wrong format: '%s'", line))
		}

		groups := sectionDefinitionRegex.FindStringSubmatch(line)

		name = groups[1]
		area, _ = strconv.ParseFloat(groups[2], 64)
		iStrong, _ = strconv.ParseFloat(groups[3], 64)
		iWeak, _ = strconv.ParseFloat(groups[4], 64)
		sStrong, _ = strconv.ParseFloat(groups[5], 64)
		sWeak, _ = strconv.ParseFloat(groups[6], 64)

		sections[name] = &structure.Section{
			Name:    name,
			Area:    area,
			IStrong: iStrong,
			IWeak:   iWeak,
			SStrong: sStrong,
			SWeak:   sWeak}
	}

	return &sections
}
