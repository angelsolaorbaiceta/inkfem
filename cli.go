package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
)

func main() {
	var (
		inputFilePathFlag  = flag.String("i", "", "input file path")
		preprocessFlag     = flag.Bool("p", false, "should dump preprocessed structure to file?")
		sysMatrixToPngFlag = flag.Bool("m", false, "should save system of equations matrix to png image file?")
		safeFlag           = flag.Bool("safe", false, "should perform safety checks?")
	)
	flag.Parse()

	if len(*inputFilePathFlag) == 0 {
		printUsage()
		return
	}

	var (
		outPath      = strings.TrimSuffix(*inputFilePathFlag, ".inkfem")
		structure    = io.StructureFromFile(*inputFilePathFlag)
		preStructure = preprocess.DoStructure(structure)
	)

	if *preprocessFlag {
		filePath := outPath + "_sliced"
		io.PreprocessedStructureToFile(preStructure, filePath)
	}

	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    *sysMatrixToPngFlag,
		OutputPath:            outPath,
		SafeChecks:            *safeFlag,
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
