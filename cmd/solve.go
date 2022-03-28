package cmd

import (
	"fmt"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
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
	solveSysMatrixToPng   bool
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
		BoolVarP(&solveUseVerbose, "verbose", "v", false, "use verbose output")

	solveCommand.
		Flags().
		BoolVarP(&solvePreprocessToFile, "preprocess", "p", false, "dump preprocessed structure to file")

	solveCommand.
		Flags().
		BoolVarP(&solveSysMatrixToPng, "matrix", "m", false, "save system matrix as a PNG image")

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
		outPath       = strings.TrimSuffix(inputFilePath, io.DefinitionFileExt)
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: solveIncludeOwnWeight}
		preStructure  *preprocess.Structure
	)

	if io.IsDefinitionFile(inputFilePath) {
		structure := readStructureFromFile(inputFilePath, readerOptions)
		preStructure = preprocessStructure(structure)

		if solvePreprocessToFile {
			go (func() {
				file := io.CreateFile(outPath + io.PreFileExt)
				defer file.Close()
				iopre.Write(preStructure, file)
			})()
		}
	} else if io.IsPreprocessedFile(inputFilePath) {
		preStructure = readPreprocessedStructureFromFile(inputFilePath)
	} else {
		panic(fmt.Sprintf("Unsuported file type: %s", inputFilePath))
	}

	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    solveSysMatrixToPng,
		OutputPath:            outPath,
		SafeChecks:            solveSafeChecks,
		MaxDisplacementsError: solveDispMaxError,
	}

	solution := process.Solve(preStructure, solveOptions)
	iosol.Write(solution, outPath+io.SolFileExt)

	log.Result()
}
