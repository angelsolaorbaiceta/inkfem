package structure

import (
	"bytes"
)

var (
	nilConstraint  = Constraint{false, false, false}
	dispConstraint = Constraint{true, true, false}
	fullConstraint = Constraint{true, true, true}
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
		return nilConstraint

	case isDxConstr && isDyConstr && !isRzConst:
		return dispConstraint

	case isDxConstr && isDyConstr && isRzConst:
		return fullConstraint

	default:
		return Constraint{isDxConstr, isDyConstr, isRzConst}
	}
}

/*
MakeNilConstraint returns a constraint which imposes no conditions of the
degrees of freedom.
*/
func MakeNilConstraint() Constraint {
	return nilConstraint
}

/*
MakeDispConstraint returns a constraint which imposes conditions on the
displacement degrees of freedom, but rotations are left free.
*/
func MakeDispConstraint() Constraint {
	return dispConstraint
}

/*
MakeFullConstraint returns a constraint which imposes conditions on all
degrees of freedom.
*/
func MakeFullConstraint() Constraint {
	return fullConstraint
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
		buffer.WriteString("rz")
	}

	buffer.WriteString("}")
	return buffer.String()
}
