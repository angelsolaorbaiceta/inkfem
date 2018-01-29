package contracts

import (
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

/*
StiffnessComputer should be implemented by anyone which can generate a
global stiffness matrix between two positions.
*/
type StiffnessComputer interface {
	StiffnessGlobalMat(startT, entT inkgeom.TParam) mat.Matrixable
}
