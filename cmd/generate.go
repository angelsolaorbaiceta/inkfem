package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	reticType     = "retic"
	reticLongType = "reticular"
)

var (
	generateType  string
	generateSpans int

	generateCommand = &cobra.Command{
		Use:   "generate --type=<type>",
		Short: "Generates a structure",
		Long:  "Generates a structure of the given type.",
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

	rootCmd.AddCommand(generateCommand)
}

func generateStructure(cmd *cobra.Command, args []string) {
	switch generateType {
	case reticType:
	case reticLongType:
		{
			fmt.Printf("generating %d spans\n", generateSpans)

		}
	default:
		{
			fmt.Printf("Unknown structure typology: \"%s\"\n", generateType)
			fmt.Println("Please use one of the following:")
			fmt.Printf("\t- \"%s\" or \"%s\": A reticular frame made of beams and columns\n", reticType, reticLongType)

			os.Exit(1)
		}
	}
}
