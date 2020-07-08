/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"flag"
	"os"
)

// CLIFlags contains the input flags to the CLI program
type CLIFlags struct {
	InputFilePath  *string
	Verbose        *bool
	Preprocess     *bool
	SysMatrixToPng *bool
	SafeChecks     *bool
	DispMaxError   *float64
}

// ParseOrShowUsage reads the program flags and parses them.
// If the required flags are not passed, shows the usage of the cli.
func ParseOrShowUsage() CLIFlags {
	flags := CLIFlags{
		InputFilePath:  flag.String("i", "", "input file path"),
		Verbose:        flag.Bool("v", false, "verbose?"),
		Preprocess:     flag.Bool("p", false, "dump preprocessed structure to file?"),
		SysMatrixToPng: flag.Bool("m", false, "save system of equations matrix to png image file?"),
		SafeChecks:     flag.Bool("safe", false, "perform safety checks?"),
		DispMaxError:   flag.Float64("e", 1e-5, "maximum allowed displacements error"),
	}

	flag.Parse()

	if len(*flags.InputFilePath) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	return flags
}
