package io

import (
	"bufio"
	"fmt"
	"os"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/process"
)

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
		tmpl   = template.Must(template.ParseFiles("io/solution.template.txt"))
		writer = bufio.NewWriter(file)
	)

	// TODO: remove
	for _, bar := range solution.Elements {
		fmt.Printf("-> bar solution %p %s\n", bar, bar.OriginalElementString())
	}

	tmpl.Execute(writer, solution)
	writer.Flush()
}
