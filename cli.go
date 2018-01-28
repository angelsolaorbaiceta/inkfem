package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
)

func main() {
	inputFilePathFlagPtr := flag.String("i", "", "input file path")
	preprocessFlagPtr := flag.Bool("p", false, "should dump preprocessed structure to file?")
	flag.Parse()

	if len(*inputFilePathFlagPtr) == 0 {
		printUsage()
		os.Exit(1)
	}

	fmt.Println("FILE:", *inputFilePathFlagPtr)
	fmt.Println("PREPROCESS:", *preprocessFlagPtr)

	structure := io.StructureFromFile(*inputFilePathFlagPtr)
	preStructure := preprocess.DoStructure(structure)

	if *preprocessFlagPtr {
		fileNameWithoutExtension := strings.TrimSuffix(*inputFilePathFlagPtr, ".inkfem")
		filePath := fileNameWithoutExtension + "_sliced"
		io.PreprocessedStructureToFile(preStructure, filePath)
	}

	process.Solve(&preStructure)
}

func printUsage() {
	fmt.Println("InkFEM usage:")
	fmt.Println("\tinkfem -i=<input_file_path> [-p]")
}
