package io

import (
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

	// Write header
	file.WriteString(
		fmt.Sprintf("inkfem v%d.%d\n", structure.Metadata.MajorVersion, structure.Metadata.MinorVersion))
	file.WriteString(
		fmt.Sprintf("|sliced_structure| %d DOFs |\n", structure.DofsCount))

	writeNodesToFile(structure.Nodes, file)
	writeElementsToFile(structure.Elements, file)
}

func writeNodesToFile(nodes map[int]structure.Node, file *os.File) {
	file.WriteString(fmt.Sprintf("\n|nodes| %d\n", len(nodes)))
	for _, val := range nodes {
		file.WriteString(val.String() + "\n")
	}
}

func writeElementsToFile(elements []preprocess.Element, file *os.File) {
	// utils.SortById(utils.ByID(elements))
	file.WriteString(fmt.Sprintf("\n|elements| %d\n", len(elements)))

	for _, element := range elements {
		file.WriteString(
			fmt.Sprintf("%s (%d)\n", element.OriginalElementString(), len(element.Nodes)))
		for _, node := range element.Nodes {
			file.WriteString("\t" + node.String() + "\n")
		}
	}
}
