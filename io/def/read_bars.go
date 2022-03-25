package def

import (
	"regexp"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	startNodeGroupName = "start_node"
	endNodeGroupName   = "end_node"
	startLinkGroupName = "start_link"
	endLinkGroupName   = "end_link"
	materialGroupName  = "material"
	sectionGroupName   = "section"
	numNodesGroupName  = "n_nodes"
)

// <id> -> <s_node> {[dx dy rz]} <e_node> {[dx dy rz]} <material> <section> [>> <n_pre_nodes>]
var elementDefinitionRegex = regexp.MustCompile(
	"^" + inkio.IdGrpExpr + inkio.ArrowExpr +
		inkio.IdGroupExpr(startNodeGroupName) + inkio.OptionalSpaceExpr +
		inkio.ConstraintGroupExpr(startLinkGroupName) + inkio.SpaceExpr +
		inkio.IdGroupExpr(endNodeGroupName) + inkio.OptionalSpaceExpr +
		inkio.ConstraintGroupExpr(endLinkGroupName) + inkio.SpaceExpr +
		inkio.NameGroupExpr(materialGroupName) + inkio.SpaceExpr +
		inkio.NameGroupExpr(sectionGroupName) + inkio.OptionalSpaceExpr +
		`(?:>>` + inkio.OptionalSpaceExpr +
		`(?P<` + numNodesGroupName + `>\d+))?` +
		"$",
)

func readBars(
	linesReader *inkio.LinesReader,
	count int,
	data *structure.StructureData,
	readerOptions inkio.ReaderOptions,
) []*structure.Element {
	var (
		lines = linesReader.GetNextLines(count)
		bars  = make([]*structure.Element, count)
	)

	for i, line := range lines {
		bars[i], _ = DeserializeBar(line, data, readerOptions)
	}

	return bars
}

// DeserializeBar parses a bar from the definition line and given the nodes, material, section
// and loads to use for its creation.
// Using the reader options, the bar can be added loads for its own weight.
//
// If the bar has the preprocess format, it also reads the number of nodes of the sliced bar
// and returns the number as the second argument.
func DeserializeBar(
	line string,
	data *structure.StructureData,
	readerOptions inkio.ReaderOptions,
) (*structure.Element, int) {
	var (
		groups        = inkio.ExtractNamedGroups(elementDefinitionRegex, line)
		id            = groups[inkio.IdGrpName]
		startNode     = data.Nodes[groups[startNodeGroupName]]
		startLink     = constraintFromString(groups[startLinkGroupName])
		endNode       = data.Nodes[groups[endNodeGroupName]]
		endLink       = constraintFromString(groups[endLinkGroupName])
		material      = data.Materials[groups[materialGroupName]]
		section       = data.Sections[groups[sectionGroupName]]
		numberOfNodes = 2
	)

	builder := structure.MakeElementBuilder(id).
		WithStartNode(startNode, startLink).
		WithEndNode(endNode, endLink).
		WithMaterial(material).
		WithSection(section).
		AddConcentratedLoads(data.ConcentratedLoads[id]).
		AddDistributedLoads(data.DistributedLoads[id])

	if readerOptions.ShouldIncludeOwnWeight {
		builder.IncludeOwnWeightLoad()
	}

	if nNodesString, isPreprocessed := groups[numNodesGroupName]; isPreprocessed {
		numberOfNodes = inkio.EnsureParseInt(nNodesString, "bar")
	}

	return builder.Build(), numberOfNodes
}
