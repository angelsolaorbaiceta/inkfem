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
StructureFromFile Reads the given .inkfem file and tries to parse a structure from the data defined.

The first line in the file should be as follows: 'inkfem vM.m', where 'M' and 'm' are the major and
minor version numbers of inkfem used to produce the file or required to compute the structure.
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
		materials                  *MaterialsByName
		sections                   *SectionsByName
		concentratedLoads          ConcLoadsById
		distributedLoads           DistLoadsById
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
				concentratedLoads, distributedLoads = readLoads(scanner, loadsCount)
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
				elements = readElements(
					scanner,
					elementsCount,
					nodes,
					materials,
					sections,
					&concentratedLoads,
					&distributedLoads,
				)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return structure.Structure{
		Metadata: structure.StrMetadata{
			MajorVersion: majorVersion,
			MinorVersion: minorVersion,
		},
		Nodes:    *nodes,
		Elements: *elements,
	}
}

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
