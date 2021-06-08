package tests

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/process"
)

var solution *process.Solution

func BenchmarkSolveStructure(b *testing.B) {
	str := io.StructureFromFile("./retic_10x5.inkfem")

	for n := 0; n < b.N; n++ {
		solution = solveStructure(&str)
	}
}
