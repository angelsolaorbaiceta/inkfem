package io

import (
	"bufio"
	_ "embed"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/process"
)

//go:embed templates/solution.template.txt
var solutionTemplateBytes []byte

// StructureSolutionToFile writes the solution of a structure to a file with the given path.
func StructureSolutionToFile(solution *process.Solution, filePath string) {
	var (
		file   = CreateFile(filePath)
		tmpl   = template.Must(template.New("solution").Parse(string(solutionTemplateBytes)))
		writer = bufio.NewWriter(file)
	)
	defer file.Close()

	tmpl.Execute(writer, solution)
	writer.Flush()
}
