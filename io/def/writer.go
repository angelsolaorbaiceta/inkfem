package def

import (
	"bufio"
	_ "embed"
	"io"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

//go:embed definition.template.txt
var definitionTemplateBytes []byte

// Write writes the given structure to the passed in writer.
func Write(structure *structure.Structure, writer io.Writer) {
	var (
		tmpl       = template.Must(template.New("definition").Parse(string(definitionTemplateBytes)))
		buffWriter = bufio.NewWriter(writer)
	)

	tmpl.Execute(buffWriter, structure)
	buffWriter.Flush()
}
