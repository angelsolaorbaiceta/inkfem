/*
Package load contains definition of loads applied to structural members.
*/
package load

import (
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/nums"
)

// A ConcentratedLoad is a load applied in a specific point.
type ConcentratedLoad struct {
	Term            Term
	IsInLocalCoords bool
	T               inkgeom.TParam
	Value           float64
}

/*
MakeConcentrated creates a concentrated load for the given term (FX, FY or MZ) which may be defined
locally to the element it will be applied to or referenced in global coordinates.

Concentrated loads are defined by a position - value tuple.
*/
func MakeConcentrated(
	term Term,
	isInLocalCoords bool,
	t inkgeom.TParam,
	value float64,
) *ConcentratedLoad {
	return &ConcentratedLoad{term, isInLocalCoords, t, value}
}

/*
IsNodal returns true if the load is applied in extreme values of T.

When `true`, this means that the load applied to an element is acting on one of its end nodes.
*/
func (load *ConcentratedLoad) IsNodal() bool {
	return load.T.IsMin() || load.T.IsMax()
}

// AsVector returns a vector with the components of the load.
func (load *ConcentratedLoad) AsVector() [3]float64 {
	return [3]float64{load.LocalFx(), load.LocalFy(), load.LocalMz()}
}

// ForcesVector returns a vector for a concentrated load with the components of {Fx, Fy}.
func (load *ConcentratedLoad) ForcesVector() g2d.Projectable {
	return g2d.MakeVector(load.LocalFx(), load.LocalFy())
}

func (load *ConcentratedLoad) LocalFx() float64 {
	if load.Term == FX {
		return load.Value
	}

	return 0.0
}

func (load *ConcentratedLoad) LocalFy() float64 {
	if load.Term == FY {
		return load.Value
	}

	return 0.0
}

func (load *ConcentratedLoad) LocalMz() float64 {
	if load.Term == MZ {
		return load.Value
	}

	return 0.0
}

// ProjectedVectorValue returns the concentrated load vector projected in a reference frame.
func (load *ConcentratedLoad) ProjectedVectorValue(refFrame g2d.RefFrame) [3]float64 {
	projectedVector := refFrame.ProjectVector(load.ForcesVector())
	return [3]float64{projectedVector.X, projectedVector.Y, load.LocalMz()}
}

// Equals tests whether the two loads are equal or not.
func (load *ConcentratedLoad) Equals(other *ConcentratedLoad) bool {
	return load.Term == other.Term &&
		load.IsInLocalCoords == other.IsInLocalCoords &&
		load.T.Equals(other.T) &&
		nums.FuzzyEqual(load.Value, other.Value)
}
