package pre

import (
	"bufio"
	"io"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Read parses a preprocessed structure from a file.
// This function requires access to the original structure in order to fill in the details.
func Read(st structure.Structure, reader io.Reader) *preprocess.Structure {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	return parse(scanner)
}

func parse(scanner *bufio.Scanner) *preprocess.Structure {
	// First line must be "inkfem vM.m"
	// scanner.Scan()
	// majorVersion, minorVersion := inkfemio.ParseVersionNumbers(scanner.Text())
	// metadata := structure.StrMetadata{
	// 	MajorVersion: majorVersion,
	// 	MinorVersion: minorVersion,
	// }

	return nil
}
