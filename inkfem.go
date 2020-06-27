package main

import (
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
)

func main() {
	flags := process.ParseOrShowUsage()

	var (
		outPath      = strings.TrimSuffix(*flags.InputFilePath, ".inkfem")
		structure    = io.StructureFromFile(*flags.InputFilePath)
		preStructure = preprocess.DoStructure(structure)
	)

	if *flags.Preprocess {
		go io.PreprocessedStructureToFile(&preStructure, outPath+".inkfempre")
	}

	// TODO: allow displacement error as flag
	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    *flags.SysMatrixToPng,
		OutputPath:            outPath,
		SafeChecks:            *flags.SafeChecks,
		MaxDisplacementsError: 1e-5,
	}

	solution := process.Solve(&preStructure, solveOptions)
	io.StructureSolutionToFile(solution, outPath+".inkfemsol")
}