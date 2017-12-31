package preprocess

import (
    "github.com/angelsolaorbaiceta/inkfem/structure"
)

// Preprocessed structure ready to be solved.
// The elements of a preprocessed structure are already sliced.
type Structure struct {
    Metadata structure.StrMetadata
    Nodes map[int]structure.Node
    Elements []Element
}
