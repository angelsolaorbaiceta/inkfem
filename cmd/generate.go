package cmd

import (
	"fmt"
	"os"

	"github.com/angelsolaorbaiceta/inkfem/generate"
	iodef "github.com/angelsolaorbaiceta/inkfem/io/def"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/spf13/cobra"
)

type GenerateTypology string

const (
	// The reticular typology is a frame made of beams and columns in a grid-like pattern.
	reticularTypology GenerateTypology = "retic"
)

func parseTypology(s string) (GenerateTypology, error) {
	switch s {
	case string(reticularTypology):
		return reticularTypology, nil
	case "reticular":
		return reticularTypology, nil
	default:
		return "", fmt.Errorf("unknown typology: \"%s\"", s)
	}
}

var (
	generateType        string
	generateSpans       int
	generateSpanLength  float64
	generateLevels      int
	generateLevelHeight float64
	generateLoadValue   float64

	generateCommand = &cobra.Command{
		Use:   "generate --type=<type>",
		Short: "Generate a structure of a given typology.",
		Long: `Generate a structure of a given typology. All bars will get assigned the same section and material.

The resulting structure will be written to the standard output in the INKFEM definition format.
Redirection to a file can be done by using the ">" operator.

The typology of structure to generate can be one of the following:
	- "retic" or "reticular": A reticular frame made of beams and columns with "spans" beams per level and "levels" columns.
		`,
		Run: generateStructure,
	}
)

func init() {
	generateCommand.
		Flags().
		StringVarP(&generateType, "type", "t", string(reticularTypology), "the typology of structure to generate. Use one of: retic, reticular")
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
		Float64VarP(&generateLoadValue, "load", "o", 50.0, "the value of the vertical load (distributed or concentrated)")

	rootCmd.AddCommand(generateCommand)
}

func generateStructure(cmd *cobra.Command, args []string) {
	var (
		ipe100Section     = structure.MakeSection("sec", 10.3, 171.0, 15.92, 34.2, 5.79)
		steelS275Material = structure.MakeMaterial("mat", 0.00000785, 21000000, 8100000, 0.3, 27500, 43000)
	)

	typology, err := parseTypology(generateType)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch typology {
	case reticularTypology:
		{
			str := generate.Reticular(generate.ReticStructureParams{
				Spans:         generateSpans,
				Span:          generateSpanLength,
				Levels:        generateLevels,
				Height:        generateLevelHeight,
				LoadDistValue: generateLoadValue,
				Section:       ipe100Section,
				Material:      steelS275Material,
			})
			iodef.Write(str, os.Stdout)
		}
	}
}
