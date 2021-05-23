package contracts

// StrID is the type used for structural data ids.
type StrID = string

/*
Identifiable is anything that can be referenced using an integer number.
*/
type Identifiable interface {
	GetID() StrID
}
