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

package structure

import (
	"bytes"
)

var (
	// NilConstraint is a constraint where all DOF are free.
	NilConstraint = Constraint{false, false, false}

	// DispConstraint is a constraint where the displacement DOFs are constrained.
	DispConstraint = Constraint{true, true, false}

	// FullConstraint is a constraint where all the DOFs are constrained.
	FullConstraint = Constraint{true, true, true}
)

/*
A Constraint represents a condition on displacements and rotations.
*/
type Constraint struct {
	isDxConstr, isDyConstr, isRzConst bool
}

/* <-- Construction --> */

/*
MakeConstraint creates a new constraint with the given degrees of freedom
constrained of freed.
*/
func MakeConstraint(isDxConstr, isDyConstr, isRzConst bool) Constraint {
	switch {
	case !isDxConstr && !isDyConstr && !isRzConst:
		return NilConstraint

	case isDxConstr && isDyConstr && !isRzConst:
		return DispConstraint

	case isDxConstr && isDyConstr && isRzConst:
		return FullConstraint

	default:
		return Constraint{isDxConstr, isDyConstr, isRzConst}
	}
}

/* <-- Properties --> */

/*
AllowsRotation returns true is rotation degree of freedom is not constrained.
*/
func (c Constraint) AllowsRotation() bool {
	return !c.isRzConst
}

/*
AllowsDispX returns true if displacement in x degree of freedom is not
constrainted.
*/
func (c Constraint) AllowsDispX() bool {
	return !c.isDxConstr
}

/*
AllowsDispY returns true if displacement in y degree of freedom is not
constrainted.
*/
func (c Constraint) AllowsDispY() bool {
	return !c.isDyConstr
}

/*
Equals tests whether this constraint equals other.
*/
func (c Constraint) Equals(other Constraint) bool {
	return c.isDxConstr == other.isDxConstr &&
		c.isDyConstr == other.isDyConstr &&
		c.isRzConst == other.isRzConst
}

/* <-- Stringer --> */

func (c Constraint) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{ ")

	if c.isDxConstr {
		buffer.WriteString("dx ")
	}
	if c.isDyConstr {
		buffer.WriteString("dy ")
	}
	if c.isRzConst {
		buffer.WriteString("rz ")
	}

	buffer.WriteString("}")
	return buffer.String()
}
