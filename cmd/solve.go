package cmd

import (
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/spf13/cobra"
)

var (
	inputFilePath    string
	includeOwnWeight bool
	dispMaxError     float64
	useVerbose       bool
	preprocessToFile bool
	sysMatrixToPng   bool
	safeChecks       bool

	solveCommand = &cobra.Command{
		Use:   "solve",
		Short: "",
		Long:  "",
		Run:   solveStructure,
	}
)

func init() {
	solveCommand.
		Flags().
		StringVarP(&inputFilePath, "input", "i", "", "Input file path (required)")
	solveCommand.MarkFlagRequired("input")

	solveCommand.
		Flags().
		BoolVarP(&includeOwnWeight, "weight", "w", false, "include the weight of the bars as a distributed load")

	solveCommand.
		Flags().
		Float64VarP(&dispMaxError, "error", "e", 1e-5, "maximum allowed displacement error")

	solveCommand.
		Flags().
		BoolVarP(&useVerbose, "verbose", "v", false, "use verbose output")

	solveCommand.
		Flags().
		BoolVarP(&preprocessToFile, "preprocess", "p", false, "dump preprocessed structure to file")

	solveCommand.
		Flags().
		BoolVarP(&sysMatrixToPng, "matrix", "m", false, "save system matrix as a PNG image")

	solveCommand.
		Flags().
		BoolVarP(&safeChecks, "safe", "s", false, "perform safety checks")

	rootCmd.AddCommand(solveCommand)
}

func solveStructure(cmd *cobra.Command, args []string) {
	log.SetVerbosity(useVerbose)
	log.StartProcess()

	var (
		outPath       = strings.TrimSuffix(inputFilePath, io.InputFileExt)
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: includeOwnWeight}
		structure     = readStructureFromFile(inputFilePath, readerOptions)
		preStructure  = preprocessStructure(structure)
	)

	if preprocessToFile {
		go io.PreprocessedStructureToFile(preStructure, outPath+io.PreFileExt)
	}

	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    sysMatrixToPng,
		OutputPath:            outPath,
		SafeChecks:            safeChecks,
		MaxDisplacementsError: dispMaxError,
	}

	solution := process.Solve(preStructure, solveOptions)
	io.StructureSolutionToFile(solution, outPath+io.SolFileExt)

	log.Result()
}

func readStructureFromFile(filePath string, readerOptions io.ReaderOptions) *structure.Structure {
	log.StartReadFile()
	structure := io.StructureFromFile(filePath, readerOptions)
	log.EndReadFile(structure.NodesCount(), structure.ElementsCount())

	return &structure
}

func preprocessStructure(structure *structure.Structure) *preprocess.Structure {
	log.StartPreprocess()
	preprocessedStructure := preprocess.StructureModel(structure)
	log.EndPreprocess()

	return preprocessedStructure
}
