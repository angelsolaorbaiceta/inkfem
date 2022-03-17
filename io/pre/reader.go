package pre

import (
	"bufio"
	"io"
	"strings"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Read parses a preprocessed structure from a file.
// This function requires access to the original structure in order to fill in the details.
func Read(st structure.Structure, reader io.Reader) *preprocess.Structure {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	metadata := inkio.ParseMetadata(scanner)

	var (
		line string
	)

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())

		if inkio.ShouldIgnoreLine(line) {
			continue
		}
	}

	return preprocess.MakeStructure(
		metadata,
		st.NodesById,
		[]*preprocess.Element{},
	)
}
