package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/inkfem/generate"
	iodef "github.com/angelsolaorbaiceta/inkfem/io/def"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/spf13/cobra"
)

const (
	reticType     = "retic"
	reticLongType = "reticular"
)

var (
	generateType        string
	generateSpans       int
	generateSpanLength  float64
	generateLevels      int
	generateLevelHeight float64
	generateloadValue   float64

	generateCommand = &cobra.Command{
		Use:   "generate --type=<type>",
		Short: "Generates a structure",
		Long:  "Generates a structure with the given typology.",
		Run:   generateStructure,
	}
)

func init() {
	generateCommand.
		Flags().
		StringVarP(&generateType, "type", "t", reticType, "the typology of structure to generate")
	generateCommand.MarkFlagRequired("type")

	generateCommand.
		Flags().
		IntVarP(&generateSpans, "spans", "s", 10, "the number of horizontal spans")

	generateCommand.
		Flags().
		Float64VarP(&generateSpanLength, "span", "p", 400.0, "the length of each span")

	generateCommand.
		Flags().
		IntVarP(&generateLevels, "levels", "l", 5, "the number of vertical levels")

	generateCommand.
		Flags().
		Float64VarP(&generateLevelHeight, "level", "e", 300.0, "the height of each level")

	generateCommand.
		Flags().
		Float64VarP(&generateloadValue, "load", "o", 50.0, "the value of the vertical load (distributed or concentrated)")

	rootCmd.AddCommand(generateCommand)
}

func generateStructure(cmd *cobra.Command, args []string) {
	switch generateType {
	case reticType, reticLongType:
		{
			str := generate.Reticular(generate.ReticStructureParams{
				Spans:         generateSpans,
				Span:          generateSpanLength,
				Levels:        generateLevels,
				Height:        generateLevelHeight,
				LoadDistValue: generateloadValue,
				Section:       structure.MakeSection("sec", 1.0, 1.0, 1.0, 1, 0),
				Material:      structure.MakeMaterial("mat", 1.0, 1.0, 1.0, 1.0, 1.0, 1.0),
			})
			iodef.Write(str, os.Stdout)
		}
	default:
		{
			fmt.Printf("Unknown structure typology: \"%s\".\n", generateType)
			fmt.Println("Use one of the following:")
			fmt.Printf(
				"\t- \"%s\" or \"%s\": A reticular frame made of beams and columns with \"spans\" beams per level and \"levels\" columns.\n",
				reticType,
				reticLongType,
			)

			os.Exit(1)
		}
	}
}
