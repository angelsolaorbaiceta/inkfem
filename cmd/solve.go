package cmd

import (
	"fmt"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
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
	major, minor := inkio.GetBinaryVersion()
	log.SetVerbosity(solveUseVerbose)
	log.StartProcess(major, minor)

	var (
		inputFilePath = args[0]
		outPath       = strings.TrimSuffix(inputFilePath, io.DefinitionFileExt)
		readerOptions = inkio.ReaderOptions{ShouldIncludeOwnWeight: solveIncludeOwnWeight}
		preStructure  *preprocess.Structure
	)

	if inkio.IsDefinitionFile(inputFilePath) {
		structure := readStructureFromFile(inputFilePath, readerOptions)
		preStructure = preprocessStructure(structure)

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
		panic(fmt.Sprintf("Unsuported file type: %s", inputFilePath))
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
