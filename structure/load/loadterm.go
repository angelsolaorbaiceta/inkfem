/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
