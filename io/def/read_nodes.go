package def

import (
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	xPosGroupName        = "x"
	yPosGroupName        = "y"
	constraintsGroupName = "constraints"
)

// <id> -> <xCoord> <yCoord> {[dx dy rz]} [| DOF: [0 1 2]]
var nodeDefinitionRegex = regexp.MustCompile(
	"^" + inkio.IdGrpExpr + inkio.ArrowExpr +
		inkio.FloatGroupExpr(xPosGroupName) + inkio.SpaceExpr +
		inkio.FloatGroupExpr(yPosGroupName) + inkio.SpaceExpr +
		inkio.ConstraintGroupExpr(constraintsGroupName) + inkio.OptionalSpaceExpr +
		`(?:\|` + inkio.OptionalSpaceExpr + inkio.DofGrpExpr + inkio.OptionalSpaceExpr + `)?` + "$",
)

// ReadNodes reads and parses "count" nodes from the lines in the lines reader.
func ReadNodes(linesReader *inkio.LinesReader, count int) map[contracts.StrID]*structure.Node {
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
		groups = inkio.ExtractNamedGroups(nodeDefinitionRegex, definition)

		id                 = groups["id"]
		x                  = inkio.EnsureParseFloat(groups[xPosGroupName], "node x position")
		y                  = inkio.EnsureParseFloat(groups[yPosGroupName], "node y position")
		externalConstraint = groups[constraintsGroupName]

		node = structure.MakeNodeAtPosition(
			id,
			x, y,
			constraintFromString(externalConstraint),
		)
	)

	if dofString, hasDof := groups[inkio.DofGrpName]; hasDof {
		dof1, dof2, dof3 := inkio.EnsureParseDOF(dofString, "node")
		node.SetDegreesOfFreedomNum(dof1, dof2, dof3)
	}

	return node
}
