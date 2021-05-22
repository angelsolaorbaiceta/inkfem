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

type ConcLoadsById = map[contracts.StrID][]*load.ConcentratedLoad
type DistLoadsById = map[contracts.StrID][]*load.DistributedLoad

func readLoads(scanner *bufio.Scanner, count int) (ConcLoadsById, DistLoadsById) {
	lines := definitionLines(scanner, count)
	return deserializeLoadsByElementID(lines)
}

func deserializeLoadsByElementID(lines []string) (ConcLoadsById, DistLoadsById) {
	var (
		elementID        contracts.StrID
		concentratedLoad *load.ConcentratedLoad
		distributedLoad  *load.DistributedLoad
		concentrated = make(ConcLoadsById)
		distributed = make(DistLoadsById)
	)

	for _, line := range lines {
		var (
			matchesDistributed  = distLoadDefinitionRegex.MatchString(line)
			matchesConcentrated = concLoadDefinitionRegex.MatchString(line)
		)

		if !(matchesDistributed || matchesConcentrated) {
			panic(fmt.Sprintf("Found load with wrong format: '%s'", line))
		}

		if matchesDistributed {
			elementID, distributedLoad = deserializeDistributedLoad(line)
			distributed[elementID] = append(distributed[elementID], distributedLoad)
		} else if matchesConcentrated {
			elementID, concentratedLoad = deserializeConcentratedLoad(line)
			concentrated[elementID] = append(concentrated[elementID], concentratedLoad)
		}
	}

	return concentrated, distributed
}

func deserializeDistributedLoad(line string) (contracts.StrID, *load.DistributedLoad) {
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

func deserializeConcentratedLoad(line string) (contracts.StrID, *load.ConcentratedLoad) {
	groups := concLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementID := groups[3]
	t := ensureParseFloat(groups[4], "concentrated load T")
	val := ensureParseFloat(groups[5], "concentrated load value")

	return elementID, load.MakeConcentrated(term, isInLocalCoords, inkgeom.MakeTParam(t), val)
}
