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
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
)

var (
	versionRegex         = regexp.MustCompile(`(?:inkfem\s+v)(\d+)(?:[.])(\d+)`)
	nodesHeaderRegex     = regexp.MustCompile(`(?:\|nodes\|\s*)(\d+)`)
	materialsHeaderRegex = regexp.MustCompile(`(?:\|materials\|\s*)(\d+)`)
	sectionsHeaderRegex  = regexp.MustCompile(`(?:\|sections\|\s*)(\d+)`)
	loadsHeaderRegex     = regexp.MustCompile(`(?:\|loads\|\s*)(\d+)`)
	elementsHeaderRegex  = regexp.MustCompile(`(?:\|elements\|\s*)(\d+)`)
)

/*
StructureFromFile Reads the given .inkfem file and tries to parse a structure from the
data defined.

The first line in the file should be as follows: 'inkfem vM.m', where 'M' and 'm' are the
major and minor version numbers of inkfem used to produce the file or required to compute
the structure.
*/
func StructureFromFile(filePath string) structure.Structure {
	file, error := os.Open(filePath)
	if error != nil {
		log.Fatal(error)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return parseStructure(scanner)
}

func parseStructure(scanner *bufio.Scanner) structure.Structure {
	var (
		nodesDefined               = false
		materialsDefined           = false
		sectionsDefined            = false
		loadsDefined               = false
		majorVersion, minorVersion int
		nodes                      *map[contracts.StrID]*structure.Node
		materials                  *map[string]*structure.Material
		sections                   *map[string]*structure.Section
		loads                      map[contracts.StrID][]load.Load
		elements                   *[]*structure.Element
	)

	// First line must be "inkfem vM.m"
	scanner.Scan()
	majorVersion, minorVersion = parseVersionNumbers(scanner.Text())

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if lineIsComment(line) || lineIsEmpty(line) {
			continue
		}

		switch {
		case nodesHeaderRegex.MatchString(line):
			{
				nodesCount, _ := strconv.Atoi(nodesHeaderRegex.FindStringSubmatch(line)[1])
				nodes = readNodes(scanner, nodesCount)
				nodesDefined = true
			}

		case materialsHeaderRegex.MatchString(line):
			{
				materialsCount, _ := strconv.Atoi(materialsHeaderRegex.FindStringSubmatch(line)[1])
				materials = readMaterials(scanner, materialsCount)
				materialsDefined = true
			}

		case sectionsHeaderRegex.MatchString(line):
			{
				sectionsCount, _ := strconv.Atoi(sectionsHeaderRegex.FindStringSubmatch(line)[1])
				sections = readSections(scanner, sectionsCount)
				sectionsDefined = true
			}

		case loadsHeaderRegex.MatchString(line):
			{
				loadsCount, _ := strconv.Atoi(loadsHeaderRegex.FindStringSubmatch(line)[1])
				loads = readLoads(scanner, loadsCount)
				loadsDefined = true
			}

		case elementsHeaderRegex.MatchString(line):
			{
				if !(nodesDefined && materialsDefined && sectionsDefined && loadsDefined) {
					panic(
						"Can't' define elements if some of the following not already defined: " +
							"nodes, materials, sections and loads",
					)
				}

				elementsCount, _ := strconv.Atoi(elementsHeaderRegex.FindStringSubmatch(line)[1])
				elements = readElements(scanner, elementsCount, nodes, materials, sections, &loads)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return structure.Structure{
		Metadata: structure.StrMetadata{
			MajorVersion: majorVersion,
			MinorVersion: minorVersion},
		Nodes:    *nodes,
		Elements: *elements}
}

/* <-- READ : Version Numbers --> */
func parseVersionNumbers(firstLine string) (majorVersion, minorVersion int) {
	if foundMatch := versionRegex.MatchString(firstLine); !foundMatch {
		panic(
			"Could not parse major and minor version numbers." +
				"Are you missing 'inkfem vM.m' in your file's first line?",
		)
	}

	versions := versionRegex.FindStringSubmatch(firstLine)
	majorVersion, _ = strconv.Atoi(versions[1])
	minorVersion, _ = strconv.Atoi(versions[2])

	return
}
