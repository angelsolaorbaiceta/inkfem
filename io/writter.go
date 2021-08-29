package io

import (
	"bufio"
	_ "embed"
	"os"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/process"
)

//go:embed templates/solution.template.txt
var solutionTemplateBytes []byte

/*
StructureSolutionToFile writes the solution of a structure to a file with the
given path.
*/
func StructureSolutionToFile(solution *process.Solution, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could not create file for the structure solution")
	}
	defer file.Close()

	var (
		tmpl   = template.Must(template.New("solution").Parse(string(solutionTemplateBytes)))
		writer = bufio.NewWriter(file)
	)

	tmpl.Execute(writer, solution)
	writer.Flush()
}
