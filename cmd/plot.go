package cmd

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/plot"
	"github.com/spf13/cobra"
)

var (
	plotScale            float64
	plotIncludeOwnWeight bool
	plotPreprocessedFile bool

	plotCommand = &cobra.Command{
		Use:   "plot <inkfem file path>",
		Short: "plots the structure to one or multiple SVG files",
		Long:  "plots the structure to one of multiple SVG files.",
		Args:  cobra.ExactArgs(1),
		Run:   plotStructure,
	}
)

func init() {
	plotCommand.
		Flags().
		Float64VarP(&plotScale, "scale", "s", 1.0, "Plot scale")

	plotCommand.
		Flags().
		BoolVarP(&plotIncludeOwnWeight, "weight", "w", false, "include the weight of each bars as a distributed load")

	plotCommand.
		Flags().
		BoolVarP(&plotPreprocessedFile, "preprocess", "p", false, "plot the preprocessed structure")

	rootCmd.AddCommand(plotCommand)
}

func plotStructure(cmd *cobra.Command, args []string) {
	var (
		inputFilePath = args[0]
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: plotIncludeOwnWeight}
		structure     = io.StructureFromFile(inputFilePath, readerOptions)
	)

	fmt.Printf("The args: %v, the file path: %s\n", args, inputFilePath)
	// fmt.Printf("%v\n", structure)

	plot.StructureToSVG(structure)
}
