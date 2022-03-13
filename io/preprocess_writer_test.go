package io

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func TestWritePreprocessedStructure(t *testing.T) {
	var (
		metadata = structure.StrMetadata{
			MajorVersion: 2,
			MinorVersion: 3,
		}
		nodesById = structure.MakeNodesById(map[contracts.StrID]*structure.Node{})
		elements  = []*preprocess.Element{}
		str       = preprocess.MakeStructure(metadata, nodesById, elements)
		writer    bytes.Buffer
	)

	WritePreprocessedStructure(str, &writer)
	gotLines := strings.Split(writer.String(), "\n")

	fmt.Println(gotLines)

	t.Run("first line is always the header with the version", func(t *testing.T) {
		want := fmt.Sprintf("inkfem v%d.%d", metadata.MajorVersion, metadata.MinorVersion)
		if got := gotLines[0]; got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})
}
