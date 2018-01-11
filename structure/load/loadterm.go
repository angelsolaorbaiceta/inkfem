package load

import "fmt"

/*
Term represents the available terms for which loads can exist:
    - Force in X
    - Force in Y
    - Moment about Z
*/
type Term string

const (
	// FX is a force in the direction of X
	FX = Term("fx")

	// FY is a force in the direction of Y
	FY = Term("fy")

	// MZ is a moment about Z
	MZ = Term("mz")
)

// IsValidTerm validates the load term, to make sure it represents a known term.
func IsValidTerm(term Term) bool {
	return (term == FX) || (term == FY) || (term == MZ)
}

// EnsureValidTerm validates the load term and panics if an unknown term happens.
func EnsureValidTerm(term Term) {
	if !IsValidTerm(term) {
		panic(fmt.Sprintf("Invalid load term: '%s'", term))
	}
}
