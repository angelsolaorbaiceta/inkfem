package process

import (
	"flag"
	"os"
)

// CLIFlags contains the input flags to the CLI program
type CLIFlags struct {
	InputFilePath          *string
	Verbose                *bool
	Preprocess             *bool
	SysMatrixToPng         *bool
	SafeChecks             *bool
	DispMaxError           *float64
	ShouldIncludeOwnWeight *bool
}

/*
ParseOrShowUsage reads the program flags and parses them.
If the required flags are not passed, shows the usage of the cli.
*/
func ParseOrShowUsage() CLIFlags {
	flags := CLIFlags{
		InputFilePath:          flag.String("i", "", "input file path"),
		Verbose:                flag.Bool("v", false, "verbose?"),
		Preprocess:             flag.Bool("p", false, "dump preprocessed structure to file?"),
		SysMatrixToPng:         flag.Bool("mat", false, "save system of equations matrix to png image file?"),
		SafeChecks:             flag.Bool("safe", false, "perform safety checks?"),
		DispMaxError:           flag.Float64("error", 1e-5, "maximum allowed displacements error"),
		ShouldIncludeOwnWeight: flag.Bool("weight", false, "include the elements self weight as distributed load?"),
	}

	flag.Parse()

	if len(*flags.InputFilePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	return flags
}
