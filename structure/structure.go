package structure

/*
Structure is a groups of resistant elements joint together to withstand external loads.
*/
type Structure struct {
	Metadata StrMetadata
	Nodes    map[int]Node
	Elements []Element
}
