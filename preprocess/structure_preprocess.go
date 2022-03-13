package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// StructureModel preprocesses the structure by concurrently slicing each of the structural members.
// The resulting sliced structure includes the degrees of freedom numbering.
func StructureModel(str *structure.Structure) *Structure {
	var (
		channel        = make(chan *Element)
		slicedElements []*Element
	)

	for _, element := range str.Elements {
		go sliceElement(element, channel)
	}

	for i := 0; i < str.ElementsCount(); i++ {
		slicedElements = append(slicedElements, <-channel)
	}

	return MakeStructure(str.Metadata, str.NodesById, slicedElements)
}
