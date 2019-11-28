package io

import (
	"bufio"
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// PreprocessedStructureToFile Writes the given preprocessed structure to a file.
func PreprocessedStructureToFile(structure preprocess.Structure, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could not create file for preprocessed structure")
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write header
	writer.WriteString(
		fmt.Sprintf("inkfem v%d.%d\n", structure.Metadata.MajorVersion, structure.Metadata.MinorVersion))
	writer.WriteString(
		fmt.Sprintf("|sliced_structure| %d DOFs\n", structure.DofsCount))

	writeNodesToFile(structure.Nodes, writer)
	writeElementsToFile(structure.Elements, writer)
	writer.Flush()
}

func writeNodesToFile(nodes map[int]structure.Node, writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("\n|nodes| %d\n", len(nodes)))
	for _, val := range nodes {
		writer.WriteString(val.String() + "\n")
	}
}

func writeElementsToFile(elements []preprocess.Element, writer *bufio.Writer) {
	// utils.SortById(utils.ByID(elements))
	writer.WriteString(fmt.Sprintf("\n|elements| %d\n", len(elements)))

	for _, element := range elements {
		writer.WriteString(
			fmt.Sprintf("%s (%d)\n", element.OriginalElementString(), len(element.Nodes)))
		for _, node := range element.Nodes {
			writer.WriteString("\t" + node.String() + "\n")
		}
	}
}
