package io

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	xPosGroupName        = "x"
	yPosGroupName        = "y"
	constraintsGroupName = "constraints"
)

// <id> -> <xCoord> <yCoord> {[dx dy rz]} [| DOF: [0 1 2]]
var nodeDefinitionRegex = regexp.MustCompile(
	"^" + IdGrpExpr + ArrowExpr +
		FloatGroupExpr(xPosGroupName) + SpaceExpr +
		FloatGroupExpr(yPosGroupName) + SpaceExpr +
		ConstraintGroupExpr(constraintsGroupName) + OptionalSpaceExpr +
		`(?:\|` + OptionalSpaceExpr + DofGrpExpr + OptionalSpaceExpr + `)?` + "$",
)

// ReadNodes reads and parses "count" nodes from the lines in the lines reader.
func ReadNodes(linesReader *LinesReader, count int) map[contracts.StrID]*structure.Node {
	lines := linesReader.GetNextLines(count)
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
		groups = ExtractNamedGroups(nodeDefinitionRegex, definition)

		id                 = groups["id"]
		x                  = EnsureParseFloat(groups[xPosGroupName], "node x position")
		y                  = EnsureParseFloat(groups[yPosGroupName], "node y position")
		externalConstraint = groups[constraintsGroupName]

		node = structure.MakeNodeAtPosition(
			id,
			x, y,
			constraintFromString(externalConstraint),
		)
	)

	if dofString, hasDof := groups[DofGrpName]; hasDof {
		var (
			dofs = strings.Fields(dofString)
			dof1 = EnsureParseInt(dofs[0], "node dx DOF")
			dof2 = EnsureParseInt(dofs[1], "node dy DOF")
			dof3 = EnsureParseInt(dofs[2], "node rz DOF")
		)

		node.SetDegreesOfFreedomNum(dof1, dof2, dof3)
	}

	return node
}
