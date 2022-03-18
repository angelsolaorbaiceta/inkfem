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
		`(?:\| \[(?P<dof>\d+ \d+ \d+)\])?` + optionalSpaceExpr +
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
		groups = ExtractNamedGroups(nodeDefinitionRegex, definition) //nodeDefinitionRegex.FindStringSubmatch(definition)

		id                 = groups["id"]
		x                  = ensureParseFloat(groups["x"], "node x position")
		y                  = ensureParseFloat(groups["y"], "node y position")
		externalConstraint = groups["constraints"]

		node = structure.MakeNodeAtPosition(
			id,
			x, y,
			constraintFromString(externalConstraint),
		)
	)

	if dofString, hasDof := groups["dof"]; hasDof {
		var (
			dofs = strings.Fields(dofString)
			dof1 = ensureParseInt(dofs[0], "node dx DOF")
			dof2 = ensureParseInt(dofs[1], "node dy DOF")
			dof3 = ensureParseInt(dofs[2], "node rz DOF")
		)

		node.SetDegreesOfFreedomNum(dof1, dof2, dof3)
	}

	return node
}
