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
	`(?P<id>\d+)(?:\s*->\s*)` +
		`(?P<start_node>\d+)(?:\s*)(?P<start_link>{.*})(?:\s+)` +
		`(?P<end_node>\d+)(?:\s*)(?P<end_link>{.*})(?:\s+)` +
		`(?P<material>'[A-Za-z0-9_ ]+')(?:\s+)` +
		`(?P<section>'[A-Za-z0-9_ ]+')`)

func readElements(
	scanner *bufio.Scanner,
	count int,
	nodes *map[int]*structure.Node,
	materials *map[string]*structure.Material,
	sections *map[string]*structure.Section,
	loads *map[int][]load.Load,
) *[]*structure.Element {
	lines := definitionLines(scanner, count)
	return deserializeElements(lines, nodes, materials, sections, loads)
}

func deserializeElements(
	lines []string,
	nodes *map[int]*structure.Node,
	materials *map[string]*structure.Material,
	sections *map[string]*structure.Section,
	loads *map[int][]load.Load,
) *[]*structure.Element {
	var (
		id, startNodeID, endNodeID int
		startNode, endNode         *structure.Node
		startLink, endLink         string
		material                   *structure.Material
		section                    *structure.Section
		ok                         bool
		elements                   = make([]*structure.Element, len(lines))
		groupName                  string
	)

	for i, line := range lines {
		if !elementDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found element with wrong format: '%s'", line))
		}

		groups := elementDefinitionRegex.FindStringSubmatch(line)

		groupName = groups[startNodeIDIndex]
		startNodeID = ensureParseInt(groupName, "element start node id")
		startNode, ok = (*nodes)[startNodeID]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown start node id: %d", id, startNodeID))
		}

		groupName = groups[endNodeIDIndex]
		endNodeID = ensureParseInt(groupName, "element end node id")
		endNode, ok = (*nodes)[endNodeID]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown end node id: %d", id, endNodeID))
		}

		groupName = groups[materialNameIndex]
		material, ok = (*materials)[groupName]
		if !ok {
			panic(fmt.Sprintf("Element %d: unknown material name: %s", id, groupName))
		}

		groupName = groups[sectionNameIndex]
		section, ok = (*sections)[groupName]
		if !ok {
			panic(fmt.Sprintf("Element %d: unknown section name: %s", id, groupName))
		}

		id = ensureParseInt(groups[idIndex], "element id")
		startLink = groups[startLinkIndex]
		endLink = groups[endLinkIndex]

		elements[i] = structure.MakeElement(
			id, startNode, endNode,
			constraintFromString(startLink),
			constraintFromString(endLink),
			material,
			section,
			(*loads)[id])
	}

	return &elements
}

func deserializeElement(
	definition string,
	nodes *map[int]*structure.Node,
	materials *map[string]*structure.Material,
	sections *map[string]*structure.Section,
	loads *map[int][]load.Load,
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
				"Element %d: unknown material name: %s", components.id, components.materialName,
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
				"Element %d: unknown section name: %s", components.id, components.sectionName,
			),
		)
	}

	return section
}
