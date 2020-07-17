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

	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
)

var (
	// <term> <reference-type> <elementId> <tStart> <valueStart> <tEnd> <valueEnd>
	distLoadDefinitionRegex = regexp.MustCompile(
		`(?P<term>[fm]{1}[xyz]{1})(?:\s+)` +
			`(?P<ref>[lg]{1})(?:d{1})(?:\s+)` +
			`(?P<element>\d+)(?:\s+)` +
			`(?P<t_start>\d+\.*\d*)(?:\s+)` +
			`(?P<val_start>-*\d+\.*\d*)(?:\s+)` +
			`(?P<t_end>\d+\.*\d*)(?:\s+)` +
			`(?P<val_end>-*\d+\.*\d*)`)

	// <term> <reference> <elementId> <t> <value>
	concLoadDefinitionRegex = regexp.MustCompile(
		`(?P<term>[fm]{1}[xyz]{1})(?:\s+)` +
			`(?P<ref>[lg]{1})(?:c{1})(?:\s+)` +
			`(?P<element>\d+)(?:\s+)` +
			`(?P<t>\d+\.*\d*)(?:\s+)` +
			`(?P<val>-*\d+\.*\d*)`)
)

func readLoads(scanner *bufio.Scanner, count int) map[int][]load.Load {
	var (
		elementNumber int
		_load         load.Load
		loads         = make(map[int][]load.Load)
	)

	for _, line := range definitionLines(scanner, count) {
		if !(distLoadDefinitionRegex.MatchString(line) || concLoadDefinitionRegex.MatchString(line)) {
			panic(fmt.Sprintf("Found load with wrong format: '%s'", line))
		}

		switch {
		case distLoadDefinitionRegex.MatchString(line):
			elementNumber, _load = deserializeDistributedLoad(line)

		case concLoadDefinitionRegex.MatchString(line):
			elementNumber, _load = deserializeConcentratedLoad(line)

		default:
			panic(fmt.Sprintf("Unknown type of load: '%s'", line))
		}

		loads[elementNumber] = append(loads[elementNumber], _load)
	}

	return loads
}

func deserializeDistributedLoad(line string) (int, load.Load) {
	groups := distLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementNumber := ensureParseInt(groups[3], "distributed load element number")
	tStart := ensureParseFloat(groups[4], "distributed load start T")
	valStart := ensureParseFloat(groups[5], "distributed load start value")
	tEnd := ensureParseFloat(groups[6], "distributed load end T")
	valEnd := ensureParseFloat(groups[7], "distributed load end value")

	return elementNumber,
		load.MakeDistributed(
			term,
			isInLocalCoords,
			inkgeom.MakeTParam(tStart),
			valStart,
			inkgeom.MakeTParam(tEnd),
			valEnd,
		)
}

func deserializeConcentratedLoad(line string) (int, load.Load) {
	groups := concLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementNumber, _ := strconv.Atoi(groups[3])
	t, _ := strconv.ParseFloat(groups[4], 64)
	val, _ := strconv.ParseFloat(groups[5], 64)

	return elementNumber, load.MakeConcentrated(term, isInLocalCoords, inkgeom.MakeTParam(t), val)
}
