package cmd

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/spf13/cobra"
)

var (
	plotInputFilePath    string
	plotScale            float64
	plotIncludeOwnWeight bool

	plotCommand = &cobra.Command{
		Use:   "plot",
		Short: "plots the structure to one or multiple SVG files",
		Long:  "plots the structure to one of multiple SVG files.",
		Run:   plotStructure,
	}
)

func init() {
	plotCommand.
		Flags().
		StringVarP(&plotInputFilePath, "input", "i", "", "Input file path (required)")
	plotCommand.MarkFlagRequired("input")

	plotCommand.
		Flags().
		Float64VarP(&plotScale, "scale", "s", 1.0, "Plot scale")

	plotCommand.
		Flags().
		BoolVarP(&plotIncludeOwnWeight, "weight", "w", false, "include the weight of each bars as a distributed load")

	rootCmd.AddCommand(plotCommand)
}

func plotStructure(cmd *cobra.Command, args []string) {
	var (
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: plotIncludeOwnWeight}
		structure     = io.StructureFromFile(plotInputFilePath, readerOptions)
	)

	fmt.Printf("%v\n", structure)
}
