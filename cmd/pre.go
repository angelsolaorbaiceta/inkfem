package cmd

import (
	"strings"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	iopre "github.com/angelsolaorbaiceta/inkfem/io/pre"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/spf13/cobra"
)

var (
	preIncludeOwnWeight bool
	preUseVerbose       bool

	preCommand = &cobra.Command{
		Use:   "pre <inkfem file path>",
		Short: "Preprocess structure",
		Long:  "Preprocesses the structure definition (.inkfem file) and saves it as a .inkfempre file.",
		Args:  cobra.ExactArgs(1),
		Run:   preStructure,
	}
)

func init() {
	preCommand.
		Flags().
		BoolVarP(&preIncludeOwnWeight, "weight", "w", false, "include the weight of each barsas a distributed load")

	preCommand.
		Flags().
		BoolVarP(&preUseVerbose, "verbose", "v", false, "use verbose output")

	rootCmd.AddCommand(preCommand)
}

func preStructure(cmd *cobra.Command, args []string) {
	log.SetVerbosity(preUseVerbose)
	log.StartProcess()

	var (
		inputFilePath = args[0]
		outPath       = strings.TrimSuffix(inputFilePath, inkio.DefinitionFileExt)
		readerOptions = inkio.ReaderOptions{ShouldIncludeOwnWeight: solveIncludeOwnWeight}
		structure     = readStructureFromFile(inputFilePath, readerOptions)
		preStructure  = preprocessStructure(structure)
	)

	file := inkio.CreateFile(outPath + inkio.PreFileExt)
	defer file.Close()
	iopre.Write(preStructure, file)

	log.Result()
}
