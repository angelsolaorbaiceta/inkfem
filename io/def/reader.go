package def

import (
	"fmt"
	"io"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Reads the given .inkfem file and tries to parse a structure from the data defined.
//
// The first line in the file should be as follows: 'inkfem vM.m', where 'M' and 'm'
// are the major and minor version numbers of inkfem used to produce the file or
// required to compute the structure.
func Read(reader io.Reader) *structure.Structure {
	linesReader := inkio.MakeLinesReader(reader)
	return parseStructure(linesReader)
}

func parseStructure(linesReader *inkio.LinesReader) *structure.Structure {
	var (
		line              string
		nodes             = make(map[contracts.StrID]*structure.Node)
		materials         = make(structure.MaterialsByName)
		sections          = make(structure.SectionsByName)
		concentratedLoads = make(structure.ConcLoadsById)
		distributedLoads  = make(structure.DistLoadsById)
		deserializedBars  = make([]*DeserializedBarDTO, 0)
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
				}

			case inkio.MaterialsHeader:
				{
					material := DeserializeMaterial(line)
					materials[material.Name] = material
				}

			case inkio.SectionsHeader:
				{
					section := DeserializeSection(line)
					sections[section.Name] = section
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
				}

			case inkio.BarsHeader:
				{
					bar, _ := DeserializeBar(line)
					deserializedBars = append(deserializedBars, bar)
				}

			default:
				panic(fmt.Sprintf("Unknown header in file: '%s'", line))
			}
		}
	}

	data := &structure.StructureData{
		Nodes:             nodes,
		Materials:         materials,
		Sections:          sections,
		ConcentratedLoads: concentratedLoads,
		DistributedLoads:  distributedLoads,
	}
	bars := BarsFromDeserialization(deserializedBars, data)

	// TODO: lines reader error handling?
	// if err := scanner.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	return structure.Make(metadata, nodes, bars)
}
