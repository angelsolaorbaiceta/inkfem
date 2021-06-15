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

package main

import (
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func main() {
	flags := process.ParseOrShowUsage()
	log.SetVerbosity(*flags.Verbose)
	log.StartProcess()

	var (
		outPath       = strings.TrimSuffix(*flags.InputFilePath, ".inkfem")
		readerOptions = io.ReaderOptions{ShouldIncludeOwnWeight: *flags.ShouldIncludeOwnWeight}
		structure     = readStructureFromFile(*flags.InputFilePath, readerOptions)
		preStructure  = preprocessStructure(structure)
	)

	if *flags.Preprocess {
		go io.PreprocessedStructureToFile(preStructure, outPath+".inkfempre")
	}

	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    *flags.SysMatrixToPng,
		OutputPath:            outPath,
		SafeChecks:            *flags.SafeChecks,
		MaxDisplacementsError: *flags.DispMaxError,
	}

	solution := process.Solve(preStructure, solveOptions)
	io.StructureSolutionToFile(solution, outPath+".inkfemsol")

	log.ResultTable()
}

func readStructureFromFile(filePath string, readerOptions io.ReaderOptions) *structure.Structure {
	log.StartReadFile()
	structure := io.StructureFromFile(filePath, readerOptions)
	log.EndReadFile(structure.NodesCount(), structure.ElementsCount())

	return &structure
}

func preprocessStructure(structure *structure.Structure) *preprocess.Structure {
	log.StartPreprocess()
	preprocessedStructure := preprocess.DoStructure(structure)
	log.EndPreprocess()

	return preprocessedStructure
}
