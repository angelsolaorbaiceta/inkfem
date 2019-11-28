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

func writeElementSolutionsToFile(elementsSolution []process.ElementSolution, writer *bufio.Writer) {
	for _, element := range elementsSolution {
		writer.WriteString("\n" + element.Element.OriginalElementString())

		writeGlobalDispl(&element, writer)
		writeLocalDispl(&element, writer)
		writeAxialStress(&element, writer)
		writeShearStress(&element, writer)
		writeBendingMoment(&element, writer)
		writer.WriteString("\n")
	}
}

func writeGlobalDispl(element *process.ElementSolution, writter *bufio.Writer) {
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
}

func writeLocalDispl(element *process.ElementSolution, writter *bufio.Writer) {
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
}

func writeAxialStress(element *process.ElementSolution, writter *bufio.Writer) {
	writter.WriteString("\n\tN >> ")
	for _, axial := range element.AxialStress {
		writter.WriteString(axial.String() + " ")
	}
}

func writeShearStress(element *process.ElementSolution, writter *bufio.Writer) {
	writter.WriteString("\n\tV >> ")
	for _, shear := range element.ShearStress {
		writter.WriteString(shear.String() + " ")
	}
}

func writeBendingMoment(element *process.ElementSolution, writter *bufio.Writer) {
	writter.WriteString("\n\tM >> ")
	for _, bending := range element.BendingMoment {
		writter.WriteString(bending.String() + " ")
	}
}
