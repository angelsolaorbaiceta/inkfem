package io

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// <id> -> <xCoord> <yCoord> {[dx dy rz]} <| DOF: [0 1 2]>
var nodeDefinitionRegex = regexp.MustCompile(
	"^" + idGrpExpr + arrowExpr +
		floatGroupExpr("x") + spaceExpr +
		floatGroupExpr("y") + spaceExpr +
		constraintGroupExpr("constraints") + optionalSpaceExpr +
		`(?:\| \[(\d+ \d+ \d+)\])?` + optionalSpaceExpr +
		"$",
)

// ReadNodes reads and parses "count" nodes from the lines in "scanner".
func ReadNodes(scanner *bufio.Scanner, count int) map[contracts.StrID]*structure.Node {
	lines := ExtractDefinitionLines(scanner, count)
	return deserializeNodesByID(lines)
}

func deserializeNodesByID(lines []string) map[contracts.StrID]*structure.Node {
	var (
		node  *structure.Node
		nodes = make(map[contracts.StrID]*structure.Node)
	)

	for _, line := range lines {
		node = deserializeNode(line)
		nodes[node.GetID()] = node
	}

	return nodes
}

func deserializeNode(definition string) *structure.Node {
	if !nodeDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found node with wrong format: '%s'", definition))
	}

	var (
		groups = nodeDefinitionRegex.FindStringSubmatch(definition)

		id                 = groups[1]
		x                  = ensureParseFloat(groups[2], "node x position")
		y                  = ensureParseFloat(groups[3], "node y position")
		externalConstraint = groups[4]

		node = structure.MakeNodeAtPosition(
			id,
			x, y,
			constraintFromString(externalConstraint),
		)
	)

	if len(groups) > 5 {
		var (
			dofsAsStrings = strings.Fields(groups[5])
			dof1          = ensureParseInt(dofsAsStrings[0], "node dx DOF")
			dof2          = ensureParseInt(dofsAsStrings[1], "node dy DOF")
			dof3          = ensureParseInt(dofsAsStrings[2], "node rz DOF")
		)

		node.SetDegreesOfFreedomNum(dof1, dof2, dof3)
	}

	return node
}
