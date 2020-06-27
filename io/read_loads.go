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
			elementNumber, _load = distributedLoadFromString(line)

		case concLoadDefinitionRegex.MatchString(line):
			elementNumber, _load = concentratedLoadFromString(line)

		default:
			// shouldn't happen
			panic("Unknown type of load?")
		}

		loads[elementNumber] = append(loads[elementNumber], _load)
	}

	return loads
}

func distributedLoadFromString(line string) (int, load.Load) {
	groups := distLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementNumber, _ := strconv.Atoi(groups[3])
	tStart, _ := strconv.ParseFloat(groups[4], 64)
	valStart, _ := strconv.ParseFloat(groups[5], 64)
	tEnd, _ := strconv.ParseFloat(groups[6], 64)
	valEnd, _ := strconv.ParseFloat(groups[7], 64)

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

func concentratedLoadFromString(line string) (int, load.Load) {
	groups := concLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementNumber, _ := strconv.Atoi(groups[3])
	t, _ := strconv.ParseFloat(groups[4], 64)
	val, _ := strconv.ParseFloat(groups[5], 64)

	return elementNumber, load.MakeConcentrated(term, isInLocalCoords, inkgeom.MakeTParam(t), val)
}
