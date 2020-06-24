package io

import (
	"bufio"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	commentDeclaration = "#"
	dispX              = "dx"
	dispY              = "dy"
	rotZ               = "rz"
)

func lineIsComment(line string) bool {
	return strings.HasPrefix(line, commentDeclaration)
}

func lineIsEmpty(line string) bool {
	return len(line) < 1
}

func definitionLines(scanner *bufio.Scanner, count int) []string {
	var (
		line  string
		lines = make([]string, count)
	)

	for i := 0; i < count; {
		if !scanner.Scan() {
			panic("Couldn't read all expected lines")
		}

		line = scanner.Text()
		if lineIsComment(line) {
			continue
		}

		lines[i] = line
		i++
	}

	return lines
}

func constraintFromString(str string) *structure.Constraint {
	var (
		dxConst = strings.Contains(str, dispX)
		dyConst = strings.Contains(str, dispY)
		rzConst = strings.Contains(str, rotZ)
	)

	return structure.MakeConstraint(dxConst, dyConst, rzConst)
}
