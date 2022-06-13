package def

import (
	"fmt"
	"regexp"

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

func DeserializeNode(definition string) *structure.Node {
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
