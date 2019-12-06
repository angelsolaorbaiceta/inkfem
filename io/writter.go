package io

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/angelsolaorbaiceta/inkfem/process"
)

// StructureSolutionToFile writes the solution of a structure to a file with the given path.
func StructureSolutionToFile(solution *process.Solution, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could not create file for the structure solution")
	}
	defer file.Close()

	solutionJSON, err := json.Marshal(solution)
	if err != nil {
		panic("Could not convert the structure solution to JSON")
	}

	writer := bufio.NewWriter(file)
	writer.Write(solutionJSON)
	writer.Flush()
}
