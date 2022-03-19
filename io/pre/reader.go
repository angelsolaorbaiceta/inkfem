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
	// scanner := bufio.NewScanner(reader)
	// scanner.Split(bufio.ScanLines)

	var (
		metadata    = inkio.ParseMetadata(linesReader)
		numberOfDof = extractNumberOfDof(linesReader)
		nodes       map[contracts.StrID]*structure.Node
		// nodesDefined = false
		line string
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
				// nodesDefined = true
			}
		}
	}

	return preprocess.MakeStructure(
		metadata,
		structure.MakeNodesById(nodes),
		[]*preprocess.Element{},
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
