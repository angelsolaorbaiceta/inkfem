package structure

import "bytes"

var (
	// NilConstraint is a constraint where all DOF are free.
	NilConstraint = Constraint{false, false, false}

	// DispConstraint is a constraint where the displacement DOFs are constrained.
	DispConstraint = Constraint{true, true, false}

	// FullConstraint is a constraint where all the DOFs are constrained.
	FullConstraint = Constraint{true, true, true}
)

// A Constraint represents a condition on displacements and rotations.
//
// Constraints are immutable, and therefore can be shared amont the elements that use them.
// Use the `MakeConstraint` factory function to get an existing instance of a constraint.
type Constraint struct {
	isDxConstr, isDyConstr, isRzConst bool
}

// MakeConstraint creates a new constraint with the given degrees of freedom constrained of free.
func MakeConstraint(isDxConstr, isDyConstr, isRzConst bool) *Constraint {
	switch {
	case !isDxConstr && !isDyConstr && !isRzConst:
		return &NilConstraint

	case isDxConstr && isDyConstr && !isRzConst:
		return &DispConstraint

	case isDxConstr && isDyConstr && isRzConst:
		return &FullConstraint

	default:
		return &Constraint{isDxConstr, isDyConstr, isRzConst}
	}
}

// AllowsRotation returns true is rotation degree of freedom is not constrained.
func (c Constraint) AllowsRotation() bool {
	return !c.isRzConst
}

// AllowsDispX returns true if displacement in x degree of freedom is not constrainted.
func (c Constraint) AllowsDispX() bool {
	return !c.isDxConstr
}

// AllowsDispY returns true if displacement in y degree of freedom is not constrainted.
func (c Constraint) AllowsDispY() bool {
	return !c.isDyConstr
}

// Equals tests whether this constraint equals other.
func (c *Constraint) Equals(other *Constraint) bool {
	return c.isDxConstr == other.isDxConstr &&
		c.isDyConstr == other.isDyConstr &&
		c.isRzConst == other.isRzConst
}

// String representation of the constraint.
// Used in the serialization format.
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
