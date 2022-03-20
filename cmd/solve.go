package cmd

import (
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
	iopre "github.com/angelsolaorbaiceta/inkfem/io/pre"
	"github.com/angelsolaorbaiceta/inkfem/log"
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
		Use:   "solve <inkfem file path>",
		Short: "solves the structure",
		Long:  "solves the structure given in an .inkfem file and saves the result in an .inkfemsol file.",
		Args:  cobra.ExactArgs(1),
		Run:   solveStructure,
	}
)

func init() {
	solveCommand.
		Flags().
		BoolVarP(&solveIncludeOwnWeight, "weight", "w", false, "include the weight of each bars as a distributed load")

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
		outPath       = strings.TrimSuffix(inputFilePath, io.InputFileExt)
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: solveIncludeOwnWeight}
		structure     = readStructureFromFile(inputFilePath, readerOptions)
		preStructure  = preprocessStructure(structure)
	)

	if solvePreprocessToFile {
		go (func() {
			file := io.CreateFile(outPath + io.PreFileExt)
			defer file.Close()
			iopre.Write(preStructure, file)
		})()
	}

	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    solveSysMatrixToPng,
		OutputPath:            outPath,
		SafeChecks:            solveSafeChecks,
		MaxDisplacementsError: solveDispMaxError,
	}

	solution := process.Solve(preStructure, solveOptions)
	io.StructureSolutionToFile(solution, outPath+io.SolFileExt)

	log.Result()
}
