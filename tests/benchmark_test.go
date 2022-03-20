package tests

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/io"
	iodef "github.com/angelsolaorbaiceta/inkfem/io/def"
	"github.com/angelsolaorbaiceta/inkfem/process"
)

var solution *process.Solution

func BenchmarkSolveStructure(b *testing.B) {
	var (
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: true}
		file          = io.OpenFile("./retic_10x5.inkfem")
		str           = iodef.Read(file, readerOptions)
	)
	defer file.Close()

	for n := 0; n < b.N; n++ {
		solution = solveStructure(str)
	}
}
