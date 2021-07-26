package io

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

// <id> -> <xCoord> <yCoord> {[dx dy rz]}
var nodeDefinitionRegex = regexp.MustCompile(
	"^" + idGrpExpr + arrowExpr +
		floatGroupExpr("x") + spaceExpr +
		floatGroupExpr("y") + spaceExpr +
		constraintGroupExpr("constraints") + optionalSpaceExpr + "$")

func readNodes(scanner *bufio.Scanner, count int) *map[contracts.StrID]*structure.Node {
	lines := definitionLines(scanner, count)
	return deserializeNodesByID(lines)
}

func deserializeNodesByID(lines []string) *map[contracts.StrID]*structure.Node {
	var (
		node  *structure.Node
		nodes = make(map[contracts.StrID]*structure.Node)
	)

	for _, line := range lines {
		node = deserializeNode(line)
		nodes[node.GetID()] = node
	}

	return &nodes
}

func deserializeNode(definition string) *structure.Node {
	if !nodeDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found node with wrong format: '%s'", definition))
	}

	groups := nodeDefinitionRegex.FindStringSubmatch(definition)

	id := groups[1]
	x := ensureParseFloat(groups[2], "node x position")
	y := ensureParseFloat(groups[3], "node y position")
	externalConstraint := groups[4]

	return structure.MakeNode(
		id,
		g2d.MakePoint(x, y),
		constraintFromString(externalConstraint),
	)
}
