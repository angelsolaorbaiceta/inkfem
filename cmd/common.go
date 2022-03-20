package cmd

import (
	"github.com/angelsolaorbaiceta/inkfem/io"
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

func preprocessStructure(structure *structure.Structure) *preprocess.Structure {
	log.StartPreprocess()
	preprocessedStructure := preprocess.StructureModel(structure)
	log.EndPreprocess()

	return preprocessedStructure
}
