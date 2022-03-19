package io

import (
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	idIndex = iota + 1
	startNodeIDIndex
	startLinkIndex
	endNodeIDIndex
	endLinkIndex
	materialNameIndex
	sectionNameIndex
)

type MaterialsByName = map[string]*structure.Material
type SectionsByName = map[string]*structure.Section

// <id> -> <s_node> {[dx dy rz]} <e_node> {[dx dy rz]} <material> <section>
var elementDefinitionRegex = regexp.MustCompile(
	"^" + idGrpExpr + arrowExpr +
		idGroupExpr("start_node") + optionalSpaceExpr +
		constraintGroupExpr("start_link") + spaceExpr +
		idGroupExpr("end_node") + optionalSpaceExpr +
		constraintGroupExpr("end_link") + spaceExpr +
		nameGroupExpr("material") + spaceExpr +
		nameGroupExpr("section") + optionalSpaceExpr + "$")

func readElements(
	linesReader *LinesReader,
	count int,
	nodes map[contracts.StrID]*structure.Node,
	materials *MaterialsByName,
	sections *SectionsByName,
	concentratedLoads *ConcLoadsById,
	distributedLoads *DistLoadsById,
	readerOptions ReaderOptions,
) []*structure.Element {
	lines := linesReader.GetNextLines(count)
	return deserializeElements(
		lines,
		nodes,
		materials,
		sections,
		concentratedLoads,
		distributedLoads,
		readerOptions,
	)
}

func deserializeElements(
	lines []string,
	nodes map[contracts.StrID]*structure.Node,
	materials *MaterialsByName,
	sections *SectionsByName,
	concentratedLoads *ConcLoadsById,
	distributedLoads *DistLoadsById,
	readerOptions ReaderOptions,
) []*structure.Element {
	elements := make([]*structure.Element, len(lines))
	for i, line := range lines {
		elements[i] = deserializeElement(
			line,
			nodes,
			materials,
			sections,
			concentratedLoads,
			distributedLoads,
			readerOptions,
		)
	}

	return elements
}

func deserializeElement(
	definition string,
	nodes map[contracts.StrID]*structure.Node,
	materials *MaterialsByName,
	sections *SectionsByName,
	concentratedLoads *ConcLoadsById,
	distributedLoads *DistLoadsById,
	readerOptions ReaderOptions,
) *structure.Element {
	var (
		components         = readElementComponents(definition)
		startNode, endNode = extractNodesForElement(components, nodes)
		material           = extractMaterialForElement(components, materials)
		section            = extractSectionForElement(components, sections)
	)

	builder := structure.MakeElementBuilder(
		components.id,
	).WithStartNode(
		startNode, components.startLink,
	).WithEndNode(
		endNode, components.endLink,
	).WithMaterial(
		material,
	).WithSection(
		section,
	).AddConcentratedLoads(
		(*concentratedLoads)[components.id],
	).AddDistributedLoads(
		(*distributedLoads)[components.id],
	)

	if readerOptions.ShouldIncludeOwnWeight {
		builder.IncludeOwnWeightLoad()
	}

	return builder.Build()
}

type elementComponents struct {
	id, startNodeID, endNodeID contracts.StrID
	materialName, sectionName  string
	startLink, endLink         *structure.Constraint
}

func readElementComponents(definition string) *elementComponents {
	if !elementDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found element with wrong format: '%s'", definition))
	}

	groups := elementDefinitionRegex.FindStringSubmatch(definition)

	return &elementComponents{
		id:           groups[idIndex],
		startNodeID:  groups[startNodeIDIndex],
		endNodeID:    groups[endNodeIDIndex],
		materialName: groups[materialNameIndex],
		sectionName:  groups[sectionNameIndex],
		startLink:    constraintFromString(groups[startLinkIndex]),
		endLink:      constraintFromString(groups[endLinkIndex]),
	}
}

func extractNodesForElement(
	components *elementComponents,
	nodes map[contracts.StrID]*structure.Node,
) (startNode, endNode *structure.Node) {
	var ok bool

	startNode, ok = nodes[components.startNodeID]
	if !ok {
		panic(
			fmt.Sprintf(
				"Element %s with unknown start node id: %s", components.id, components.startNodeID,
			),
		)
	}

	endNode, ok = nodes[components.endNodeID]
	if !ok {
		panic(
			fmt.Sprintf(
				"Element %s with unknown end node id: %s", components.id, components.endNodeID,
			),
		)
	}

	return
}

func extractMaterialForElement(
	components *elementComponents,
	materials *map[string]*structure.Material,
) *structure.Material {
	material, ok := (*materials)[components.materialName]
	if !ok {
		panic(
			fmt.Sprintf(
				"Element %s: couldn't find material with name '%s'",
				components.id,
				components.materialName,
			),
		)
	}

	return material
}

func extractSectionForElement(
	components *elementComponents,
	sections *map[string]*structure.Section,
) *structure.Section {
	section, ok := (*sections)[components.sectionName]
	if !ok {
		panic(
			fmt.Sprintf(
				"Element %s: couldn't find section with name '%s'",
				components.id,
				components.sectionName,
			),
		)
	}

	return section
}
