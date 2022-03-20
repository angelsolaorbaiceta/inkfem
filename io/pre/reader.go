package pre

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var (
	dofRegex = regexp.MustCompile(`dof_count:\s*(\d+)`)
)

// Read parses a preprocessed structure from a file.
func Read(reader io.Reader) *preprocess.Structure {
	linesReader := inkio.MakeLinesReader(reader)

	var (
		metadata         = inkio.ParseMetadata(linesReader)
		numberOfDof      = extractNumberOfDof(linesReader)
		nodes            map[contracts.StrID]*structure.Node
		materials        structure.MaterialsByName
		sections         structure.SectionsByName
		bars             []*preprocess.Element
		nodesDefined     = false
		materialsDefined = false
		sectionsDefined  = false
		line             string
	)

	for linesReader.ReadNext() {
		line = linesReader.GetNextLine()

		if inkio.ShouldIgnoreLine(line) {
			continue
		}

		switch {
		case inkio.IsNodesHeader(line):
			{
				nodesCount := inkio.ExtractNodesCount(line)
				nodes = inkio.ReadNodes(linesReader, nodesCount)
				nodesDefined = true
			}

		case inkio.IsMaterialsHeader(line):
			{
				materialsCount := inkio.ExtractMaterialsCount(line)
				materials = inkio.ReadMaterials(linesReader, materialsCount)
				materialsDefined = true
			}

		case inkio.IsSectionsHeader(line):
			{
				sectionsCount := inkio.ExtractSectionsCount(line)
				sections = inkio.ReadSections(linesReader, sectionsCount)
				sectionsDefined = true
			}

		case inkio.IsBarsHeader(line):
			{
				if !(nodesDefined && materialsDefined && sectionsDefined) {
					panic(
						"Can't' parse the bars if some of the following isn't already parsed: " +
							"nodes, materials and sections",
					)
				}

				barsCount := inkio.ExtractBarsCount(line)
				data := &structure.StructureData{
					Nodes:             nodes,
					Materials:         materials,
					Sections:          sections,
					ConcentratedLoads: structure.ConcLoadsById{},
					DistributedLoads:  structure.DistLoadsById{},
				}
				bars = readBars(linesReader, barsCount, data)
			}
		}

	}

	return preprocess.MakeStructure(
		metadata,
		structure.MakeNodesById(nodes),
		bars,
	).SetDofsCount(numberOfDof)
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
