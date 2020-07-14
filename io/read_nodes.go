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
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

// <id> -> <xCoord> <yCoord> {[dx dy rz]}
var nodeDefinitionRegex = regexp.MustCompile(
	`(?P<id>\d+)(?:\s*->\s*)` +
		`(?P<x>\d+\.*\d*)(?:\s+)` +
		`(?P<y>\d+\.*\d*)(?:\s+)` +
		`(?P<constraints>{.*})`)

func readNodes(scanner *bufio.Scanner, count int) *map[int]*structure.Node {
	var (
		node  *structure.Node
		nodes = make(map[int]*structure.Node)
	)

	for _, line := range definitionLines(scanner, count) {
		node = deserializeNode(line)
		nodes[node.Id] = node
	}

	return &nodes
}

func deserializeNode(definition string) *structure.Node {
	if !nodeDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found node with wrong format: '%s'", definition))
	}

	groups := nodeDefinitionRegex.FindStringSubmatch(definition)

	id, _ := strconv.Atoi(groups[1])
	x, _ := strconv.ParseFloat(groups[2], 64)
	y, _ := strconv.ParseFloat(groups[3], 64)
	externalConstraint := groups[4]

	return structure.MakeNode(
		id,
		g2d.MakePoint(x, y),
		constraintFromString(externalConstraint),
	)
}
