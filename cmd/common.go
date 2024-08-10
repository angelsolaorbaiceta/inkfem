package cmd

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/io"
	iodef "github.com/angelsolaorbaiceta/inkfem/io/def"
	iopre "github.com/angelsolaorbaiceta/inkfem/io/pre"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// readStructureFromFile reads the structure definition from the given .inkfem file.
// If the file is not a .inkfem file, it panics.
func readStructureFromFile(filePath string, readerOptions io.ReaderOptions) *structure.Structure {
	if !io.IsDefinitionFile(filePath) {
		panic(fmt.Sprintf("Expected %s file: %s", io.DefinitionFileExt, filePath))
	}

	log.StartReadFile()

	file := io.OpenFile(filePath)
	defer file.Close()

	structure := iodef.Read(file, readerOptions)
	log.EndReadFile(io.DefinitionFileExt, structure.NodesCount(), structure.ElementsCount())

	return structure
}

func readPreprocessedStructureFromFile(filePath string) *preprocess.Structure {
	if !io.IsPreprocessedFile(filePath) {
		panic(fmt.Sprintf("Expected %s file: %s", io.PreFileExt, filePath))
	}

	log.StartReadFile()

	file := io.OpenFile(filePath)
	defer file.Close()

	preStructure := iopre.Read(file)
	log.EndReadFile(io.PreFileExt, preStructure.NodesCount(), preStructure.ElementsCount())

	return preStructure
}

func preprocessStructure(structure *structure.Structure) *preprocess.Structure {
	log.StartPreprocess()
	preprocessedStructure := preprocess.StructureModel(structure)
	log.EndPreprocess()

	return preprocessedStructure
}
