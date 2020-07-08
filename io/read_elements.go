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
	"strconv"

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
	var (
		id, startNodeID, endNodeID int
		startNode, endNode         *structure.Node
		startLink, endLink         string
		material                   *structure.Material
		section                    *structure.Section
		ok                         bool
		elements                   = make([]*structure.Element, count)
		groupName                  string
	)

	for i, line := range definitionLines(scanner, count) {
		if !elementDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found element with wrong format: '%s'", line))
		}

		groups := elementDefinitionRegex.FindStringSubmatch(line)

		groupName = groups[startNodeIDIndex]
		startNodeID, _ = strconv.Atoi(groupName)
		startNode, ok = (*nodes)[startNodeID]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown start node id: %d", id, startNodeID))
		}

		groupName = groups[endNodeIDIndex]
		endNodeID, _ = strconv.Atoi(groupName)
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

		id, _ = strconv.Atoi(groups[idIndex])
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
