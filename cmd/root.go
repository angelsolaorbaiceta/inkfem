package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "inkfem",
	Short: "Solves a structure",
	Long:  "Finite Element Method CLI to solve linear two-dimensional structures defined as .inkfem files.",
}

// Execute adds all child commands to the root command and sets flags appropiately.
// This is called by inkfem.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {}
