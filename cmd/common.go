package cmd

import (
	"os"

	"github.com/angelsolaorbaiceta/inkfem/io"
	iopre "github.com/angelsolaorbaiceta/inkfem/io/pre"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func readStructureFromFile(filePath string, readerOptions io.ReaderOptions) *structure.Structure {
	log.StartReadFile()
	structure := io.StructureFromFile(filePath, readerOptions)
	log.EndReadFile(structure.NodesCount(), structure.ElementsCount())

	return structure
}

func readPreprocessedStructureFromFile(filePath string) *preprocess.Structure {
	log.StartReadFile()

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	preStructure := iopre.Read(file)
	log.EndReadFile(preStructure.NodesCount(), preStructure.ElementsCount())

	return preStructure
}

func preprocessStructure(structure *structure.Structure) *preprocess.Structure {
	log.StartPreprocess()
	preprocessedStructure := preprocess.StructureModel(structure)
	log.EndPreprocess()

	return preprocessedStructure
}
