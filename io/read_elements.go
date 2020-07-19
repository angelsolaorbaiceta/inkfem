/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import (
	"bufio"
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
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
	scanner *bufio.Scanner,
	count int,
	nodes *map[contracts.StrID]*structure.Node,
	materials *map[string]*structure.Material,
	sections *map[string]*structure.Section,
	loads *map[contracts.StrID][]load.Load,
) *[]*structure.Element {
	lines := definitionLines(scanner, count)
	return deserializeElements(lines, nodes, materials, sections, loads)
}

func deserializeElements(
	lines []string,
	nodes *map[contracts.StrID]*structure.Node,
	materials *map[string]*structure.Material,
	sections *map[string]*structure.Section,
	loads *map[contracts.StrID][]load.Load,
) *[]*structure.Element {
	var (
		element  *structure.Element
		elements = make([]*structure.Element, len(lines))
	)

	for i, line := range lines {
		element = deserializeElement(line, nodes, materials, sections, loads)
		elements[i] = element
	}

	return &elements
}

func deserializeElement(
	definition string,
	nodes *map[contracts.StrID]*structure.Node,
	materials *map[string]*structure.Material,
	sections *map[string]*structure.Section,
	loads *map[contracts.StrID][]load.Load,
) *structure.Element {
	components := readElementComponents(definition)
	startNode, endNode := extractNodesForElement(components, nodes)
	material := extractMaterialForElement(components, materials)
	section := extractSectionForElement(components, sections)

	return structure.MakeElement(
		components.id,
		startNode,
		endNode,
		components.startLink,
		components.endLink,
		material,
		section,
		(*loads)[components.id],
	)
}

type elementComponents struct {
	id, startNodeID, endNodeID contracts.StrID
	materialName, sectionName  string
	startLink, endLink         structure.Constraint
}

func readElementComponents(definition string) *elementComponents {
	if !elementDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found element with wrong format: '%s'", definition))
	}

	groups := elementDefinitionRegex.FindStringSubmatch(definition)

	return &elementComponents{
		id:           ensureParseInt(groups[idIndex], "element id"),
		startNodeID:  ensureParseInt(groups[startNodeIDIndex], "element start node id"),
		endNodeID:    ensureParseInt(groups[endNodeIDIndex], "element end node id"),
		materialName: groups[materialNameIndex],
		sectionName:  groups[sectionNameIndex],
		startLink:    constraintFromString(groups[startLinkIndex]),
		endLink:      constraintFromString(groups[endLinkIndex]),
	}
}

func extractNodesForElement(
	components *elementComponents,
	nodes *map[int]*structure.Node,
) (startNode, endNode *structure.Node) {
	var ok bool

	startNode, ok = (*nodes)[components.startNodeID]
	if !ok {
		panic(
			fmt.Sprintf(
				"Element %d with unknown start node id: %d", components.id, components.startNodeID,
			),
		)
	}

	endNode, ok = (*nodes)[components.endNodeID]
	if !ok {
		panic(
			fmt.Sprintf(
				"Element %d with unknown end node id: %d", components.id, components.endNodeID,
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
				"Element %d: couldn't find material with name '%s'",
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
				"Element %d: couldn't find section with name '%s'",
				components.id,
				components.sectionName,
			),
		)
	}

	return section
}
