package def

import (
	"fmt"
	"io"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Read Reads the given .inkfem file and tries to parse a structure from the data defined.
//
// The first line in the file should be as follows: 'inkfem vM.m', where 'M' and 'm' are the major and
// minor version numbers of inkfem used to produce the file or required to compute the structure.
func Read(reader io.Reader, options inkio.ReaderOptions) *structure.Structure {
	linesReader := inkio.MakeLinesReader(reader)
	return parseStructure(linesReader, options)
}

func parseStructure(
	linesReader *inkio.LinesReader,
	options inkio.ReaderOptions,
) *structure.Structure {
	var (
		line              string
		nodesDefined      = false
		materialsDefined  = false
		sectionsDefined   = false
		loadsDefined      = false
		nodes             map[contracts.StrID]*structure.Node
		materials         structure.MaterialsByName
		sections          structure.SectionsByName
		concentratedLoads structure.ConcLoadsById
		distributedLoads  structure.DistLoadsById
		bars              []*structure.Element
	)

	// First line must be "inkfem vM.m"
	metadata := inkio.ParseMetadata(linesReader)

	for linesReader.ReadNext() {
		line = linesReader.GetNextLine()

		switch {
		case inkio.IsNodesHeader(line):
			{
				nodesCount := inkio.ExtractNodesCount(line)
				nodes = ReadNodes(linesReader, nodesCount)
				nodesDefined = true
			}

		case inkio.IsMaterialsHeader(line):
			{
				materialsCount := inkio.ExtractMaterialsCount(line)
				materials = ReadMaterials(linesReader, materialsCount)
				materialsDefined = true
			}

		case inkio.IsSectionsHeader(line):
			{
				sectionsCount := inkio.ExtractSectionsCount(line)
				sections = ReadSections(linesReader, sectionsCount)
				sectionsDefined = true
			}

		case inkio.IsLoadsHeader(line):
			{
				loadsCount := inkio.ExtractLoadsCount(line)
				concentratedLoads, distributedLoads = readLoads(linesReader, loadsCount)
				loadsDefined = true
			}

		case inkio.IsBarsHeader(line):
			{
				if !(nodesDefined && materialsDefined && sectionsDefined && loadsDefined) {
					panic(
						"Can't' parse the bars if any of the following isn't already parsed: " +
							"nodes, materials, sections and loads",
					)
				}

				barsCount := inkio.ExtractBarsCount(line)
				data := &structure.StructureData{
					Nodes:             nodes,
					Materials:         materials,
					Sections:          sections,
					ConcentratedLoads: concentratedLoads,
					DistributedLoads:  distributedLoads,
				}
				bars = readBars(linesReader, barsCount, data, options)
			}

		default:
			panic(fmt.Sprintf("Unknown header in file: '%s'", line))
		}

	}

	// TODO: lines reader error handling?
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	return structure.Make(metadata, nodes, bars)
}
