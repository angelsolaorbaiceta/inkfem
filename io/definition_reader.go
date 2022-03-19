package io

import (
	"fmt"
	"log"
	"os"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
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

	linesReader := MakeLinesReader(file)
	return parseStructure(linesReader, options)
}

func parseStructure(linesReader *LinesReader, options ReaderOptions) *structure.Structure {
	var (
		line              string
		nodesDefined      = false
		materialsDefined  = false
		sectionsDefined   = false
		loadsDefined      = false
		nodes             map[contracts.StrID]*structure.Node
		materials         *MaterialsByName
		sections          *SectionsByName
		concentratedLoads ConcLoadsById
		distributedLoads  DistLoadsById
		elements          []*structure.Element
	)

	// First line must be "inkfem vM.m"
	metadata := ParseMetadata(linesReader)

	for linesReader.ReadNext() {
		line = linesReader.GetNextLine()

		switch {
		case IsNodesHeader(line):
			{
				nodesCount := ExtractNodesCount(line)
				nodes = ReadNodes(linesReader, nodesCount)
				nodesDefined = true
			}

		case IsMaterialsHeader(line):
			{
				materialsCount := ExtractMaterialsCount(line)
				materials = ReadMaterials(linesReader, materialsCount)
				materialsDefined = true
			}

		case IsSectionsHeader(line):
			{
				sectionsCount := ExtractSectionsCount(line)
				sections = ReadSections(linesReader, sectionsCount)
				sectionsDefined = true
			}

		case IsLoadsHeader(line):
			{
				loadsCount := ExtractLoadsCount(line)
				concentratedLoads, distributedLoads = readLoads(linesReader, loadsCount)
				loadsDefined = true
			}

		case IsBarsHeader(line):
			{
				if !(nodesDefined && materialsDefined && sectionsDefined && loadsDefined) {
					panic(
						"Can't' define elements if some of the following not already defined: " +
							"nodes, materials, sections and loads",
					)
				}

				elementsCount := ExtractBarsCount(line)
				elements = readElements(
					linesReader,
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

	// TODO: lines reader error handling?
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	return structure.Make(metadata, nodes, elements)
}
