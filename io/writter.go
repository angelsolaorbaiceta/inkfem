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

	writeElementSolutionsToFile(solution.Elements, file)
}

func writeElementSolutionsToFile(elementsSolution []process.ElementSolution, file *os.File) {
	for _, element := range elementsSolution {
		file.WriteString("\n" + element.Element.OriginalElementString())

		// Points
		// file.WriteString("\n\t")
		// for _, point := range element.Points {
		// 	file.WriteString(fmt.Sprintf("T = %f : %s ", t.Value(), point.String()))
		// }

		// Global Displacements
		file.WriteString("\n\tgDx >> ")
		for _, disp := range element.GlobalXDispl {
			file.WriteString(disp.String() + " ")
		}

		file.WriteString("\n\tgDy >> ")
		for _, disp := range element.GlobalYDispl {
			file.WriteString(disp.String() + " ")
		}

		file.WriteString("\n\tgRz >> ")
		for _, rot := range element.GlobalZRot {
			file.WriteString(rot.String() + " ")
		}

		// Local Displacements
		file.WriteString("\n\tlDx >> ")
		for _, disp := range element.LocalXDispl {
			file.WriteString(disp.String() + " ")
		}

		file.WriteString("\n\tlDy >> ")
		for _, disp := range element.LocalYDispl {
			file.WriteString(disp.String() + " ")
		}

		file.WriteString("\n\tlRz >> ")
		for _, rot := range element.LocalZRot {
			file.WriteString(rot.String() + " ")
		}

		// Axial Stress
		file.WriteString("\n\tN >> ")
		for _, axial := range element.AxialStress {
			file.WriteString(axial.String() + " ")
		}

		// Shear Stress
		file.WriteString("\n\tV >> ")
		for _, shear := range element.ShearStress {
			file.WriteString(shear.String() + " ")
		}

		// Bending Moment
		file.WriteString("\n\tM >> ")
		for _, bending := range element.BendingMoment {
			file.WriteString(bending.String() + " ")
		}

		file.WriteString("\n")
	}
}
