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

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
)

var (
	// <term> <reference-type> <elementId> <tStart> <valueStart> <tEnd> <valueEnd>
	distLoadDefinitionRegex = regexp.MustCompile(
		"^" + loadTermExpr + distributedLoadRefExpr +
			loadElementID +
			floatGroupExpr("t_start") + spaceExpr +
			floatGroupExpr("val_start") + spaceExpr +
			floatGroupExpr("t_end") + spaceExpr +
			floatGroupExpr("val_end") + optionalSpaceExpr + "$",
	)

	// <term> <reference> <elementId> <t> <value>
	concLoadDefinitionRegex = regexp.MustCompile(
		"^" + loadTermExpr + concentratedLoadRefExpr +
			loadElementID +
			floatGroupExpr("t") + spaceExpr +
			floatGroupExpr("val") + optionalSpaceExpr + "$",
	)
)

func readLoads(scanner *bufio.Scanner, count int) map[contracts.StrID][]load.Load {
	lines := definitionLines(scanner, count)
	return deserializeLoadsByElementID(lines)
}

func deserializeLoadsByElementID(lines []string) map[contracts.StrID][]load.Load {
	var (
		elementID contracts.StrID
		_load     load.Load
		loads     = make(map[contracts.StrID][]load.Load)
	)

	for _, line := range lines {
		var (
			matchesDistributed  = distLoadDefinitionRegex.MatchString(line)
			matchesConcentrated = concLoadDefinitionRegex.MatchString(line)
		)

		if !(matchesDistributed || matchesConcentrated) {
			panic(fmt.Sprintf("Found load with wrong format: '%s'", line))
		}

		switch {
		case matchesDistributed:
			elementID, _load = deserializeDistributedLoad(line)

		case matchesConcentrated:
			elementID, _load = deserializeConcentratedLoad(line)

		default:
			panic(fmt.Sprintf("Unknown type of load: '%s'", line))
		}

		loads[elementID] = append(loads[elementID], _load)
	}

	return loads
}

func deserializeDistributedLoad(line string) (contracts.StrID, load.Load) {
	groups := distLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementID := groups[3]
	tStart := ensureParseFloat(groups[4], "distributed load start T")
	valStart := ensureParseFloat(groups[5], "distributed load start value")
	tEnd := ensureParseFloat(groups[6], "distributed load end T")
	valEnd := ensureParseFloat(groups[7], "distributed load end value")

	return elementID,
		load.MakeDistributed(
			term,
			isInLocalCoords,
			inkgeom.MakeTParam(tStart),
			valStart,
			inkgeom.MakeTParam(tEnd),
			valEnd,
		)
}

func deserializeConcentratedLoad(line string) (contracts.StrID, load.Load) {
	groups := concLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementID := groups[3]
	t := ensureParseFloat(groups[4], "concentrated load T")
	val := ensureParseFloat(groups[5], "concentrated load value")

	return elementID,
		load.MakeConcentrated(term, isInLocalCoords, inkgeom.MakeTParam(t), val)
}
