package io

import (
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// TODO: use group names instead
const (
	idIndex = iota + 1
	startNodeIDIndex
	startLinkIndex
	endNodeIDIndex
	endLinkIndex
	materialNameIndex
	sectionNameIndex
)

const (
	startNodeGroupName = "start_node"
	endNodeGroupName   = "end_node"
	startLinkGroupName = "start_link"
	endLinkGroupName   = "end_link"
	materialGroupName  = "material"
	sectionGroupName   = "section"
)

// <id> -> <s_node> {[dx dy rz]} <e_node> {[dx dy rz]} <material> <section>
var elementDefinitionRegex = regexp.MustCompile(
	"^" + IdGrpExpr + ArrowExpr +
		IdGroupExpr(startNodeGroupName) + OptionalSpaceExpr +
		ConstraintGroupExpr(startLinkGroupName) + SpaceExpr +
		IdGroupExpr(endNodeGroupName) + OptionalSpaceExpr +
		ConstraintGroupExpr(endLinkGroupName) + SpaceExpr +
		NameGroupExpr(materialGroupName) + SpaceExpr +
		NameGroupExpr(sectionGroupName) + OptionalSpaceExpr + "$",
)

func readBars(
	linesReader *LinesReader,
	count int,
	data *structure.StructureData,
	readerOptions ReaderOptions,
) []*structure.Element {
	var (
		lines = linesReader.GetNextLines(count)
		bars  = make([]*structure.Element, count)
	)

	for i, line := range lines {
		bars[i] = DeserializeBar(line, data, readerOptions)
	}

	return bars
}

// DeserializeBar parses a bar from the definition line and given the nodes, material, section
// and loads to use for its creation.
// Using the reader options, the bar can be added loads for its own weight.
func DeserializeBar(
	line string,
	data *structure.StructureData,
	readerOptions ReaderOptions,
) *structure.Element {
	var (
		groups    = ExtractNamedGroups(elementDefinitionRegex, line)
		id        = groups[IdGrpName]
		startNode = data.Nodes[groups[startNodeGroupName]]
		startLink = constraintFromString(groups[startLinkGroupName])
		endNode   = data.Nodes[groups[endNodeGroupName]]
		endLink   = constraintFromString(groups[endLinkGroupName])
		material  = data.Materials[groups[materialGroupName]]
		section   = data.Sections[groups[sectionGroupName]]
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

	return builder.Build()
}
