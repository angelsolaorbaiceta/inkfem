package io

import (
	"bufio"
	_ "embed"
	"io"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

//go:embed templates/preprocess.template.txt
var preprocessTemplateBytes []byte

// WritePreprocessedStructure Writes the given preprocessed structure to the passed in writer.
func WritePreprocessedStructure(structure *preprocess.Structure, writer io.Writer) {
	var (
		tmpl       = template.Must(template.New("preprocess").Parse(string(preprocessTemplateBytes)))
		buffWriter = bufio.NewWriter(writer)
	)

	tmpl.Execute(buffWriter, structure)
	buffWriter.Flush()
}
