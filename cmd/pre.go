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
		Use:   "pre [-w] [-v] <.inkfem file path>",
		Short: "Preprocess structure",
		Long: `Preprocess the structure definition (.inkfem file), slicing it and distributing the loads into the nodes, and saves it as a .inkfempre file.

Each bar is sliced into a number of elements, and the loads are distributed into the nodes of the structure.
How many elements each bar is sliced into depends on its end supports and applied loads.
There are three different cases:

1. Axial bars: These are not sliced at all.
2. Bars without loads: These are sliced into a maximum of 6 elements.
3. Bars with loads: These are sliced into a maximum of 10 elements.

Intermediate points where a concentrated load is applied also generate intermediate nodes for the load to be included.

When the -w flag is used, the weight of each bar is included as a distributed load.
		`,
		Args: cobra.ExactArgs(1),
		Run:  preStructure,
	}
)

func init() {
	preCommand.
		Flags().
		BoolVarP(&preIncludeOwnWeight, "weight", "w", false, "include the weight of each bar as a distributed load")

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
		structure     = readStructureFromFile(inputFilePath)
		preStructure  = preprocessStructure(structure)
	)

	file := inkio.CreateFile(outPath + inkio.PreFileExt)
	defer file.Close()
	iopre.Write(preStructure, file)

	log.Result()
}
