package io

import (
	"fmt"
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

	// Write header
	file.WriteString(
		fmt.Sprintf("inkfem v%d.%d\n", solution.Metadata.MajorVersion, solution.Metadata.MinorVersion))
	file.WriteString(
		fmt.Sprintf("|solution| %d Elements\n", len(solution.Elements)))
	file.WriteString("[points, globalDisp, localDisp, axial, shear, bending]\n")

	writeElementSolutionsToFile(solution.Elements, file)
}

func writeElementSolutionsToFile(elementsSolution []process.ElementSolution, file *os.File) {
	for _, element := range elementsSolution {
		file.WriteString("\n" + element.Element.OriginalElementString())

		file.WriteString("\n\t")
		for t, point := range element.Points {
			file.WriteString(fmt.Sprintf("T = %f : %s ", t.Value(), point.String()))
		}

		file.WriteString("\n\t")
		for t, disp := range element.GlobalDispl {
			file.WriteString(fmt.Sprintf("T = %f : {%f, %f, %f} ", t.Value(), disp[0], disp[1], disp[2]))
		}

		file.WriteString("\n")
	}
}
