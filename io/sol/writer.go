package sol

import (
	"bufio"
	_ "embed"
	"io"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/process"
)

//go:embed solution.template.txt
var solutionTemplateBytes []byte

// Write writes the solution of a structure to the passed in writer.
func Write(solution *process.Solution, writer io.Writer) {
	var (
		tmpl       = template.Must(template.New("solution").Parse(string(solutionTemplateBytes)))
		buffWriter = bufio.NewWriter(writer)
	)

	tmpl.Execute(buffWriter, solution)
	buffWriter.Flush()
}
