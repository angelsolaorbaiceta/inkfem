package contracts

// StrID is the type used for structural data ids.
type StrID = string

// Identifiable is anything that can be identified using a string.
type Identifiable interface {
	GetID() StrID
}
