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
		nodes             = make(map[contracts.StrID]*structure.Node)
		materials         = make(structure.MaterialsByName)
		sections          = make(structure.SectionsByName)
		concentratedLoads = make(structure.ConcLoadsById)
		distributedLoads  = make(structure.DistLoadsById)
		bars              = make([]*structure.Element, 0)
		currentSection    string
	)

	// First line must be "inkfem vM.m"
	metadata := inkio.ParseMetadata(linesReader)

	for linesReader.ReadNext() {
		line = linesReader.GetNextLine()

		if inkio.IsSectionHeaderLine(line) {
			currentSection = inkio.ParseSectionHeader(line)
		} else {
			switch currentSection {
			case inkio.NodesHeader:
				{
					node := DeserializeNode(line)
					nodes[node.GetID()] = node
					nodesDefined = true
				}

			case inkio.MaterialsHeader:
				{
					material := DeserializeMaterial(line)
					materials[material.Name] = material
					materialsDefined = true
				}

			case inkio.SectionsHeader:
				{
					section := DeserializeSection(line)
					sections[section.Name] = section
					sectionsDefined = true
				}

			case inkio.LoadsHeader:
				{
					barId, distLoad, concLoad := DeserializeLoad(line)
					if distLoad != nil {
						distributedLoads[barId] = append(distributedLoads[barId], distLoad)
					}
					if concLoad != nil {
						concentratedLoads[barId] = append(concentratedLoads[barId], concLoad)
					}
					loadsDefined = true
				}

			case inkio.BarsHeader:
				{
					if !(nodesDefined && materialsDefined && sectionsDefined && loadsDefined) {
						panic(
							"Can't' parse the bars if any of the following isn't already parsed: " +
								"nodes, materials, sections and loads",
						)
					}

					data := &structure.StructureData{
						Nodes:             nodes,
						Materials:         materials,
						Sections:          sections,
						ConcentratedLoads: concentratedLoads,
						DistributedLoads:  distributedLoads,
					}
					bar, _ := DeserializeBar(line, data, options)
					bars = append(bars, bar)
				}

			default:
				panic(fmt.Sprintf("Unknown header in file: '%s'", line))
			}
		}
	}

	// TODO: lines reader error handling?
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	return structure.Make(metadata, nodes, bars)
}
