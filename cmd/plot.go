package cmd

import (
	"os"

	"github.com/angelsolaorbaiceta/inkfem/plot"
	"github.com/spf13/cobra"
)

var (
	plotScale            float64
	plotPreprocessedFile bool

	plotCommand = &cobra.Command{
		Use:   "plot <inkfem file path>",
		Short: "Plots the structure to one or multiple SVG files",
		Long:  "Plots the structure to one of multiple SVG files.",
		Args:  cobra.ExactArgs(1),
		Run:   plotStructure,
	}
)

func init() {
	plotCommand.
		Flags().
		Float64VarP(&plotScale, "scale", "s", 0.25, "Plot scale")

	plotCommand.
		Flags().
		BoolVarP(&plotPreprocessedFile, "preprocess", "p", false, "plot the preprocessed structure")

	rootCmd.AddCommand(plotCommand)
}

func plotStructure(cmd *cobra.Command, args []string) {
	var (
		inputFilePath         = args[0]
		structurePlotFilePath = inputFilePath + ".svg"
		structure             = readStructureFromFile(inputFilePath)
		strPlotOptions        = plot.StructurePlotOps{
			Scale:     plotScale,
			MinMargin: 100,
		}
	)

	strPlotFile, err := os.Create(structurePlotFilePath)
	if err != nil {
		panic("Could not create file for the structure drawing")
	}
	defer strPlotFile.Close()

	plot.StructureToSVG(structure, strPlotOptions, strPlotFile)
}
