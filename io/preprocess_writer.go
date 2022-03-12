package io

import (
	"bufio"
	_ "embed"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

//go:embed templates/preprocess.template.txt
var preprocessTemplateBytes []byte

// PreprocessedStructureToFile Writes the given preprocessed structure to a file.
func PreprocessedStructureToFile(structure *preprocess.Structure, filePath string) {
	var (
		file   = createFile(filePath)
		tmpl   = template.Must(template.New("preprocess").Parse(string(preprocessTemplateBytes)))
		writer = bufio.NewWriter(file)
	)
	defer file.Close()

	tmpl.Execute(writer, structure)
	writer.Flush()
}
