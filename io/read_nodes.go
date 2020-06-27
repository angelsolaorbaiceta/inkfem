package io

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
)

// <id> -> <xCoord> <yCoord> {[dx dy rz]}
var nodeDefinitionRegex = regexp.MustCompile(
	`(?P<id>\d+)(?:\s*->\s*)` +
		`(?P<x>\d+\.*\d*)(?:\s+)` +
		`(?P<y>\d+\.*\d*)(?:\s+)` +
		`(?P<constraints>{.*})`)

func readNodes(scanner *bufio.Scanner, count int) *map[int]*structure.Node {
	var (
		id                 int
		x, y               float64
		externalConstraint string
		nodes              = make(map[int]*structure.Node)
	)

	for _, line := range definitionLines(scanner, count) {
		if !nodeDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found node with wrong format: '%s'", line))
		}

		groups := nodeDefinitionRegex.FindStringSubmatch(line)

		id, _ = strconv.Atoi(groups[1])
		x, _ = strconv.ParseFloat(groups[2], 64)
		y, _ = strconv.ParseFloat(groups[3], 64)
		externalConstraint = groups[4]

		nodes[id] = structure.MakeNode(
			id,
			inkgeom.MakePoint(x, y),
			constraintFromString(externalConstraint))
	}

	return &nodes
}
