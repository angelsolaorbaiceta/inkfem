package io

import (
	"bufio"
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

	writer := bufio.NewWriter(file)

	// Write header
	writer.WriteString(
		fmt.Sprintf("inkfem v%d.%d\n", solution.Metadata.MajorVersion, solution.Metadata.MinorVersion))
	writer.WriteString(
		fmt.Sprintf("|solution| %d Elements\n", len(solution.Elements)))

	writeElementSolutionsToFile(solution.Elements, writer)
	writer.Flush()
}

func writeElementSolutionsToFile(elementsSolution []process.ElementSolution, writter *bufio.Writer) {
	for _, element := range elementsSolution {
		writter.WriteString("\n" + element.Element.OriginalElementString())

		// Points
		// file.WriteString("\n\t")
		// for _, point := range element.Points {
		// 	file.WriteString(fmt.Sprintf("T = %f : %s ", t.Value(), point.String()))
		// }

		// Global Displacements
		writter.WriteString("\n\tgDx >> ")
		for _, disp := range element.GlobalXDispl {
			writter.WriteString(disp.String() + " ")
		}

		writter.WriteString("\n\tgDy >> ")
		for _, disp := range element.GlobalYDispl {
			writter.WriteString(disp.String() + " ")
		}

		writter.WriteString("\n\tgRz >> ")
		for _, rot := range element.GlobalZRot {
			writter.WriteString(rot.String() + " ")
		}

		// Local Displacements
		writter.WriteString("\n\tlDx >> ")
		for _, disp := range element.LocalXDispl {
			writter.WriteString(disp.String() + " ")
		}

		writter.WriteString("\n\tlDy >> ")
		for _, disp := range element.LocalYDispl {
			writter.WriteString(disp.String() + " ")
		}

		writter.WriteString("\n\tlRz >> ")
		for _, rot := range element.LocalZRot {
			writter.WriteString(rot.String() + " ")
		}

		// Axial Stress
		writter.WriteString("\n\tN >> ")
		for _, axial := range element.AxialStress {
			writter.WriteString(axial.String() + " ")
		}

		// Shear Stress
		writter.WriteString("\n\tV >> ")
		for _, shear := range element.ShearStress {
			writter.WriteString(shear.String() + " ")
		}

		// Bending Moment
		writter.WriteString("\n\tM >> ")
		for _, bending := range element.BendingMoment {
			writter.WriteString(bending.String() + " ")
		}

		writter.WriteString("\n")
	}
}
