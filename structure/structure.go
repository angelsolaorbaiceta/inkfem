package structure

/*
Structure is a group of resistant elements joined together designed to withstand
the application of external loads.
*/
type Structure struct {
	Metadata StrMetadata
	Nodes    map[int]*Node
	Elements []*Element
}
