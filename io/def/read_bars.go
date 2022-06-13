package def

import (
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
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

type DeserializeBarDTO struct {
	Id           contracts.StrID
	StartNodeId  contracts.StrID
	StartLink    *structure.Constraint
	EndNodeId    contracts.StrID
	EndLink      *structure.Constraint
	MaterialName string
	SectionName  string
}

func (bar *DeserializeBarDTO) Equals(other *DeserializeBarDTO) bool {
	return bar.Id == other.Id &&
		bar.StartNodeId == other.StartNodeId &&
		bar.EndNodeId == other.EndNodeId &&
		bar.MaterialName == other.MaterialName &&
		bar.SectionName == other.SectionName &&
		bar.StartLink.Equals(other.StartLink) &&
		bar.EndLink.Equals(other.EndLink)
}

// DeserializeBar parses a bar from the definition line. The bar is a deserialization data transfer
// object containing the data needed to create the bar. It references the ids of the nodes and names
// of the materials and sections.
//
// If the bar has the preprocess format, it also reads the number of nodes of the sliced bar
// and returns the number as the second argument.
func DeserializeBar(line string) (*DeserializeBarDTO, int) {
	var (
		groups        = inkio.ExtractNamedGroups(elementDefinitionRegex, line)
		id            = groups[inkio.IdGrpName]
		startNodeId   = groups[startNodeGroupName]
		startLink     = constraintFromString(groups[startLinkGroupName])
		endNodeId     = groups[endNodeGroupName]
		endLink       = constraintFromString(groups[endLinkGroupName])
		materialName  = groups[materialGroupName]
		sectionName   = groups[sectionGroupName]
		numberOfNodes = 2
	)

	bar := &DeserializeBarDTO{
		Id:           id,
		StartNodeId:  startNodeId,
		StartLink:    startLink,
		EndNodeId:    endNodeId,
		EndLink:      endLink,
		MaterialName: materialName,
		SectionName:  sectionName,
	}

	if nNodesString, isPreprocessed := groups[numNodesGroupName]; isPreprocessed {
		numberOfNodes = inkio.EnsureParseInt(nNodesString, "bar")
	}

	return bar, numberOfNodes
}

// BarsFromDeserialization maps the deserialization data transfer objects to the structure elements
// given the structure data (nodes, sections, materials and loads).
//
// Using the reader options, the bar can be added loads for its own weight.
func BarsFromDeserialization(
	deserializedBars []*DeserializeBarDTO,
	data *structure.StructureData,
	readerOptions inkio.ReaderOptions,
) []*structure.Element {
	var bars []*structure.Element

	for _, deserializedBar := range deserializedBars {
		bars = append(bars, BarFromDeserialization(deserializedBar, data, readerOptions))
	}

	return bars
}

func BarFromDeserialization(
	bar *DeserializeBarDTO,
	data *structure.StructureData,
	readerOptions inkio.ReaderOptions,
) *structure.Element {
	var (
		startNode = data.Nodes[bar.StartNodeId]
		endNode   = data.Nodes[bar.EndNodeId]
		material  = data.Materials[bar.MaterialName]
		section   = data.Sections[bar.SectionName]
	)

	builder := structure.MakeElementBuilder(bar.Id).
		WithStartNode(startNode, bar.StartLink).
		WithEndNode(endNode, bar.EndLink).
		WithMaterial(material).
		WithSection(section).
		AddConcentratedLoads(data.ConcentratedLoads[bar.Id]).
		AddDistributedLoads(data.DistributedLoads[bar.Id])

	if readerOptions.ShouldIncludeOwnWeight {
		builder.IncludeOwnWeightLoad()
	}

	return builder.Build()
}
