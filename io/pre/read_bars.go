package pre

import (
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// <id> -> <s_node> {[dx dy rz]} <e_node> {[dx dy rz]} <material> <section>
// var elementDefinitionRegex = regexp.MustCompile(
// 	"^" + idGrpExpr + arrowExpr +
// 		idGroupExpr("start_node") + optionalSpaceExpr +
// 		constraintGroupExpr("start_link") + spaceExpr +
// 		idGroupExpr("end_node") + optionalSpaceExpr +
// 		constraintGroupExpr("end_link") + spaceExpr +
// 		nameGroupExpr("material") + spaceExpr +
// 		nameGroupExpr("section") + optionalSpaceExpr + "$")

func readBars(
	linesReader *inkio.LinesReader,
	count int,
	data *structure.StructureData,
) []*preprocess.Element {
	return []*preprocess.Element{}
}
