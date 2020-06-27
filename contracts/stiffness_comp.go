package contracts

import (
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

/*
A StiffnessComputer can generate a global stiffness matrix between two positions
of a directrix defined by the startT and endT parameter values.
*/
type StiffnessComputer interface {
	StiffnessGlobalMat(startT, entT inkgeom.TParam) mat.ReadOnlyMatrix
}
