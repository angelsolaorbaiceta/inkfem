package process

import (
	"flag"
	"fmt"
	"os"
)

// CLIFlags contains the input flags to the CLI program
type CLIFlags struct {
	InputFilePath  *string
	Preprocess     *bool
	SysMatrixToPng *bool
	SafeChecks     *bool
}

// ParseOrShowUsage reads the program flags and parses them.
// If the required flags are not passed, shows the usage of the cli.
func ParseOrShowUsage() CLIFlags {
	flags := CLIFlags{
		InputFilePath:  flag.String("i", "", "input file path"),
		Preprocess:     flag.Bool("p", false, "dump preprocessed structure to file?"),
		SysMatrixToPng: flag.Bool("m", false, "save system of equations matrix to png image file?"),
		SafeChecks:     flag.Bool("safe", false, "perform safety checks?"),
	}

	flag.Parse()

	if len(*flags.InputFilePath) == 0 {
		printUsage()
		os.Exit(1)
	}

	return flags
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
