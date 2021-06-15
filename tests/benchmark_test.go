package tests

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/process"
)

var solution *process.Solution

func BenchmarkSolveStructure(b *testing.B) {
	var (
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: true}
		str           = io.StructureFromFile("./retic_10x5.inkfem", readerOptions)
	)

	for n := 0; n < b.N; n++ {
		solution = solveStructure(&str)
	}
}
