package pre

import (
	"bufio"
	"io"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Read parses a preprocessed structure from a file.
// This function requires access to the original structure in order to fill in the details.
func Read(st structure.Structure, reader io.Reader) *preprocess.Structure {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	metadata := parseMetadata(scanner)

	return preprocess.MakeStructure(metadata, st.NodesById, []*preprocess.Element{})
}

func parseMetadata(scanner *bufio.Scanner) structure.StrMetadata {
	// First line must be "inkfem vM.m"
	scanner.Scan()
	majorVersion, minorVersion := inkio.ParseVersionNumbers(scanner.Text())
	return structure.StrMetadata{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
	}
}
