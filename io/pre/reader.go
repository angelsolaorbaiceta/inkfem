package pre

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	iodef "github.com/angelsolaorbaiceta/inkfem/io/def"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var (
	dofRegex       = regexp.MustCompile(`dof_count:\s*(\d+)`)
	ownWeightRegex = regexp.MustCompile(`includes_own_weight:\s*(yes|no)`)
)

// Read parses a preprocessed structure from an .inkfempre file.
func Read(reader io.Reader) *preprocess.Structure {
	linesReader := inkio.MakeLinesReader(reader)

	var (
		metadata          = inkio.ParseMetadata(linesReader)
		numberOfDof       = extractNumberOfDof(linesReader)
		includesOwnWeight = extractIncludesOwnWeight(linesReader)
		nodes             = make(map[contracts.StrID]*structure.Node)
		materials         = make(structure.MaterialsByName)
		sections          = make(structure.SectionsByName)
		bars              = make([]*preprocess.Element, 0)
		nodesDefined      = false
		materialsDefined  = false
		sectionsDefined   = false
		line              string
		currentSection    string
	)

	for linesReader.ReadNext() {
		line = linesReader.GetNextLine()

		if inkio.IsSectionHeaderLine(line) {
			currentSection = inkio.ParseSectionHeader(line)
		} else {
			switch currentSection {
			case inkio.NodesHeader:
				{
					node := iodef.DeserializeNode(line)
					nodes[node.GetID()] = node
					nodesDefined = true
				}

			case inkio.MaterialsHeader:
				{
					material := iodef.DeserializeMaterial(line)
					materials[material.Name] = material
					materialsDefined = true
				}

			case inkio.SectionsHeader:
				{
					section := iodef.DeserializeSection(line)
					sections[section.Name] = section
					sectionsDefined = true
				}

			case inkio.BarsHeader:
				{
					if !(nodesDefined && materialsDefined && sectionsDefined) {
						panic(
							"Can't' parse the bars if some of the following isn't already parsed: " +
								"nodes, materials and sections",
						)
					}

					data := &structure.StructureData{
						Nodes:             nodes,
						Materials:         materials,
						Sections:          sections,
						ConcentratedLoads: structure.ConcLoadsById{},
						DistributedLoads:  structure.DistLoadsById{},
					}
					bar := DeserializeBar(linesReader, data)
					bars = append(bars, bar)
				}

			default:
				panic(fmt.Sprintf("Unknown header in file: '%s'", line))
			}
		}
	}

	return preprocess.MakeStructure(
		metadata,
		structure.MakeNodesById(nodes),
		bars,
		includesOwnWeight,
	).SetDofsCount(numberOfDof) // TODO: should read the DOFs from the file, not reassign them
}

func extractNumberOfDof(linesReader *inkio.LinesReader) int {
	linesReader.ReadNext()

	line := linesReader.GetNextLine()
	if dofRegex.MatchString(line) {
		dofs, err := strconv.Atoi(dofRegex.FindStringSubmatch(line)[1])
		if err != nil {
			panic(fmt.Sprintf("Can't read number of degrees of freedom from '%s'", line))
		}

		return dofs
	}

	panic("Preprocessed file without 'dof_count' set")
}

func extractIncludesOwnWeight(linesReader *inkio.LinesReader) bool {
	linesReader.ReadNext()

	line := linesReader.GetNextLine()
	if ownWeightRegex.MatchString(line) {
		return ownWeightRegex.FindStringSubmatch(line)[1] == "yes"
	}

	panic("Preprocessed file without 'includes_own_weight' set")
}
