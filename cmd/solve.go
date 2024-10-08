package cmd

import (
	"fmt"
	"strings"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	iopre "github.com/angelsolaorbaiceta/inkfem/io/pre"
	iosol "github.com/angelsolaorbaiceta/inkfem/io/sol"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/spf13/cobra"
)

var (
	solveIncludeOwnWeight bool
	solveDispMaxError     float64
	solveUseVerbose       bool
	solvePreprocessToFile bool
	solveSafeChecks       bool

	solveCommand = &cobra.Command{
		Use:   "solve <inkfem|inkfempre file path>",
		Short: "Solves the structure",
		Long:  "Solves the structure given in an .inkfem or preprocessed .inkfempre file and saves the result in an .inkfemsol file.",
		Args:  cobra.ExactArgs(1),
		Run:   solveStructure,
	}
)

func init() {
	solveCommand.
		Flags().
		BoolVarP(&solveIncludeOwnWeight, "weight", "w", false, "include the weight of each bar as a distributed load")

	solveCommand.
		Flags().
		Float64VarP(&solveDispMaxError, "error", "e", 1e-5, "maximum allowed displacement error")

	solveCommand.
		Flags().
		BoolVarP(&solveUseVerbose, "verbose", "v", false, "use verbose output; logs the progress of the resolution")

	solveCommand.
		Flags().
		BoolVarP(&solvePreprocessToFile, "preprocess", "p", false, "dump the preprocessed structure to an .inkfempre file")

	solveCommand.
		Flags().
		BoolVarP(&solveSafeChecks, "safe", "s", false, "perform safety checks")

	rootCmd.AddCommand(solveCommand)
}

func solveStructure(cmd *cobra.Command, args []string) {
	log.SetVerbosity(solveUseVerbose)
	log.StartProcess()

	var (
		inputFilePath = args[0]
		outPath       = strings.TrimSuffix(inputFilePath, inkio.DefinitionFileExt)
		preStructure  *preprocess.Structure
	)

	if inkio.IsDefinitionFile(inputFilePath) {
		var (
			structure = readStructureFromFile(inputFilePath)
			options   = &preprocess.PreprocessOptions{
				IncludeOwnWeight: solveIncludeOwnWeight,
			}
		)

		preStructure = preprocessStructure(structure, options)

		if solvePreprocessToFile {
			go (func() {
				file := inkio.CreateFile(outPath + inkio.PreFileExt)
				defer file.Close()
				iopre.Write(preStructure, file)
			})()
		}
	} else if inkio.IsPreprocessedFile(inputFilePath) {
		preStructure = readPreprocessedStructureFromFile(inputFilePath)
	} else {
		panic(
			fmt.Sprintf(
				"Unsupported file type: %s. Expected %s or %s\n",
				inputFilePath, inkio.DefinitionFileExt, inkio.PreFileExt,
			),
		)
	}

	solveOptions := process.SolveOptions{
		OutputPath:            outPath,
		SafeChecks:            solveSafeChecks,
		MaxDisplacementsError: solveDispMaxError,
	}

	var (
		solution = process.Solve(preStructure, solveOptions)
		solFile  = inkio.CreateFile(outPath + inkio.SolFileExt)
	)
	defer solFile.Close()

	iosol.Write(solution, solFile)

	log.Result()
}
