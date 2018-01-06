package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
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

	wg := new(sync.WaitGroup)
	structure := io.StructureFromFile(*inputFilePathFlagPtr)
	preStructure := preprocess.DoStructure(structure, wg)

	if *preprocessFlagPtr {
		fileNameWithoutExtension := strings.TrimSuffix(*inputFilePathFlagPtr, ".inkfem")
		filePath := fileNameWithoutExtension + "_sliced"
		io.PreprocessedStructureToFile(preStructure, filePath)
	}
}

func printUsage() {
	fmt.Println("InkFEM usage:")
	fmt.Println("\tinkfem -i=<input_file_path> [-p]")
}
