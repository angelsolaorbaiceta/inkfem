package pre

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

var (
	dofRegex = regexp.MustCompile(`dof_count:\s*(\d+)`)
)

// Read parses a preprocessed structure from a file.
// This function requires access to the original structure in order to fill in the details.
func Read(st structure.Structure, reader io.Reader) *preprocess.Structure {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var (
		metadata    = inkio.ParseMetadata(scanner)
		numberOfDof = extractNumberOfDof(scanner)
	)

	var (
		line string
	)

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())

		if inkio.ShouldIgnoreLine(line) {
			continue
		}

		switch {

		}
	}

	return preprocess.MakeStructure(
		metadata,
		st.NodesById,
		[]*preprocess.Element{},
	).SetDofsCount(numberOfDof)
}

func extractNumberOfDof(scanner *bufio.Scanner) int {
	var line string

	for scanner.Scan() {
		line = scanner.Text()

		if dofRegex.MatchString(line) {
			dofs, err := strconv.Atoi(dofRegex.FindStringSubmatch(line)[1])
			if err != nil {
				panic(fmt.Sprintf("Can't read number of degrees of freedom from '%s'", line))
			}

			return dofs
		}
	}

	panic("Preprocessed file without 'dof_count' set")
}
