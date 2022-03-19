package pre

import (
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func readBars(
	linesReader *inkio.LinesReader,
	count int,
	data *structure.StructureData,
) []*preprocess.Element {
	var (
		lines = linesReader.GetNextLines(count)
		bars  = make([]*preprocess.Element, count)
	)

	for i, line := range lines {
		bars[i] = deserializeBar(line, data)
	}

	return bars
}

func deserializeBar(line string, data *structure.StructureData) *preprocess.Element {
	originalElement := inkio.DeserializeBar(line, data, inkio.ReaderOptions{ShouldIncludeOwnWeight: false})

	return preprocess.MakeElement(originalElement, []*preprocess.Node{})
}
