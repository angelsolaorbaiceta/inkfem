package io

import (
	"bufio"
	"os"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

// PreprocessedStructureToFile Writes the given preprocessed structure to a file.
func PreprocessedStructureToFile(structure *preprocess.Structure, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could not create file for preprocessed structure")
	}
	defer file.Close()

	var (
		tmpl   = template.Must(template.ParseFiles("io/preprocess.template.txt"))
		writer = bufio.NewWriter(file)
	)

	tmpl.Execute(writer, structure)
	writer.Flush()
}
