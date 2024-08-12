package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// PreprocessOptions contains the options to configure the preprocessing of the
// structure.
type PreprocessOptions struct {
	// IncludeOwnWeight indicates whether the weight of each bar should be included
	// as a distributed load.
	IncludeOwnWeight bool
}

// StructureModel preprocesses the structure by concurrently slicing each of the
// structural members: the bars.
// The resulting sliced structure includes the degrees of freedom numbering needed
// in the resolution of the system of equations.
//
// The passed in options are used to configure the preprocessing.
// See the PreprocessOptions struct for more information.
func StructureModel(str *structure.Structure, options *PreprocessOptions) *Structure {
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
		if options.IncludeOwnWeight {
			element.AddOwnWeight()
		}

		go sliceElement(element, channel)
	}

	for i := 0; i < numOfBars; i++ {
		slicedElements[i] = <-channel
	}
	close(channel)

	return MakeStructure(metadata, str.NodesById.Copy(), slicedElements).AssignDof()
}
