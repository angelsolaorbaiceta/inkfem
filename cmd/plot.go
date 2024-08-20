package cmd

import (
	"os"

	"github.com/angelsolaorbaiceta/inkfem/plot"
	"github.com/spf13/cobra"
)

var (
	plotScale            float64
	plotPreprocessedFile bool
	plotUseDarkTheme     bool
	plotDistLoadScale    float64

	plotCommand = &cobra.Command{
		Use:   "plot <inkfem file path>",
		Short: "Plot the structure to one or multiple SVG files",
		Long: `Plot the structure to one of multiple SVG files.
	
The original structure definition (.inkfem file) is plotted to an SVG file with the same name, but with the .svg extension.
This plot includes the bars, node supports, and loads.
		`,
		Args: cobra.ExactArgs(1),
		Run:  plotStructure,
	}
)

func init() {
	plotCommand.
		Flags().
		Float64VarP(&plotScale, "scale", "s", 0.25, "Plot scale")
	plotCommand.
		Flags().
		Float64Var(&plotDistLoadScale, "dload-scale", 0.5, "Scale factor for distributed loads")

	plotCommand.
		Flags().
		BoolVarP(&plotPreprocessedFile, "preprocess", "p", false, "plot the preprocessed structure (if the .inkfempre file can be found)")

	plotCommand.
		Flags().
		BoolVarP(&plotUseDarkTheme, "dark", "d", false, "use a dark theme for the plot")

	rootCmd.AddCommand(plotCommand)
}

func plotStructure(cmd *cobra.Command, args []string) {
	var (
		inputFilePath         = args[0]
		structurePlotFilePath = inputFilePath + ".svg"
		structure             = readStructureFromFile(inputFilePath)
		strPlotOptions        = &plot.StructurePlotOps{
			Scale: plotScale,
			// TODO: read these two values from the command line
			DistLoadScale: plotDistLoadScale,
			MinMargin:     150,
		}
		plotConfig *plot.PlotConfig
	)

	if plotUseDarkTheme {
		plotConfig = plot.DarkPlotConfig()
	} else {
		plotConfig = plot.DefaultPlotConfig()
	}

	strPlotFile, err := os.Create(structurePlotFilePath)
	if err != nil {
		panic("Could not create file for the structure drawing")
	}
	defer strPlotFile.Close()

	plot.StructureToSVG(structure, strPlotOptions, plotConfig, strPlotFile)
}
