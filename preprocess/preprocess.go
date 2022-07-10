package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// StructureModel preprocesses the structure by concurrently slicing each of the structural members.
// The resulting sliced structure includes the degrees of freedom numbering.
func StructureModel(str *structure.Structure) *Structure {
	var (
		numOfBars      = str.ElementsCount()
		channel        = make(chan *Element, numOfBars)
		slicedElements = make([]*Element, numOfBars)
		metadata       = structure.StrMetadata{
			MajorVersion: build.Info.MajorVersion,
			MinorVersion: build.Info.MinorVersion,
		}
	)

	for _, element := range str.Elements() {
		go sliceElement(element, channel)
	}

	for i := 0; i < numOfBars; i++ {
		slicedElements[i] = <-channel
	}
	close(channel)

	return MakeStructure(metadata, str.NodesById, slicedElements).AssignDof()
}
