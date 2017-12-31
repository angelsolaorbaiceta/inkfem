package main

import (
    "fmt"
    "os"
    "sync"
    "flag"
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

    fmt.Println("***************************************")
    fmt.Println(preStructure)
}

func printUsage() {
    fmt.Println("InkFEM usage:")
    fmt.Println("\tinkfem -i=<input_file_path> [-p]")
}
