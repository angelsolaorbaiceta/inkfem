package io

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var (
	nodesHeaderRegex     = regexp.MustCompile(`(?:\|nodes\|\s*)(\d+)`)
	materialsHeaderRegex = regexp.MustCompile(`(?:\|materials\|\s*)(\d+)`)
	sectionsHeaderRegex  = regexp.MustCompile(`(?:\|sections\|\s*)(\d+)`)
	loadsHeaderRegex     = regexp.MustCompile(`(?:\|loads\|\s*)(\d+)`)
	barsHeaderRegex      = regexp.MustCompile(`(?:\|bars\|\s*)(\d+)`)
)

// StructureFromFile Reads the given .inkfem file and tries to parse a structure from the data defined.
//
// The first line in the file should be as follows: 'inkfem vM.m', where 'M' and 'm' are the major and
// minor version numbers of inkfem used to produce the file or required to compute the structure.
func StructureFromFile(filePath string, options ReaderOptions) *structure.Structure {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return parseStructure(scanner, options)
}

func parseStructure(scanner *bufio.Scanner, options ReaderOptions) *structure.Structure {
	var (
		nodesDefined      = false
		materialsDefined  = false
		sectionsDefined   = false
		loadsDefined      = false
		nodes             *map[contracts.StrID]*structure.Node
		materials         *MaterialsByName
		sections          *SectionsByName
		concentratedLoads ConcLoadsById
		distributedLoads  DistLoadsById
		elements          *[]*structure.Element
	)

	// First line must be "inkfem vM.m"
	metadata := ParseMetadata(scanner)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if ShouldIgnoreLine(line) {
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

		case barsHeaderRegex.MatchString(line):
			{
				if !(nodesDefined && materialsDefined && sectionsDefined && loadsDefined) {
					panic(
						"Can't' define elements if some of the following not already defined: " +
							"nodes, materials, sections and loads",
					)
				}

				elementsCount, _ := strconv.Atoi(barsHeaderRegex.FindStringSubmatch(line)[1])
				elements = readElements(
					scanner,
					elementsCount,
					nodes,
					materials,
					sections,
					&concentratedLoads,
					&distributedLoads,
					options,
				)
			}

		default:
			panic(fmt.Sprintf("Unknown header in file: '%s'", line))
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return structure.Make(metadata, *nodes, *elements)
}
