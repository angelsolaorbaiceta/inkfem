package main

import (
	"fmt"
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
		filePath := outPath + "_sliced"
		go io.PreprocessedStructureToFile(preStructure, filePath)
	}

	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    *flags.SysMatrixToPng,
		OutputPath:            outPath,
		SafeChecks:            *flags.SafeChecks,
		MaxDisplacementsError: 1e-5,
	}

	solution := process.Solve(&preStructure, solveOptions)
	io.StructureSolutionToFile(solution, outPath+".inksol")
}

func printUsage() {
	fmt.Println("InkFEM usage:")
	fmt.Println("\tinkfem -i=<input_file_path> [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("\t-p: save preprocessed structure to file")
	fmt.Println("\t-m: save system of equations matrix to png image file")
	fmt.Println("\t-safe: do safe checks for conditions that must be satisfied during analysis")
}
