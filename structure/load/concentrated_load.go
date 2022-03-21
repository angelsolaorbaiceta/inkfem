// Package load contains definition of loads applied to structural members.
package load

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// A ConcentratedLoad is a load applied in a specific point.
type ConcentratedLoad struct {
	Term            Term
	IsInLocalCoords bool
	T               nums.TParam
	Value           float64
}

// MakeConcentrated creates a concentrated load for the given term (FX, FY or MZ) which may be defined
// localy to the element it will be applied to or referenced in global coordinates.
//
// Concentrated loads are defined by a position - value tuple.
func MakeConcentrated(
	term Term,
	isInLocalCoords bool,
	t nums.TParam,
	value float64,
) *ConcentratedLoad {
	return &ConcentratedLoad{term, isInLocalCoords, t, value}
}

// IsNodal returns true if the load is applied in extreme values of T.
// When `true`, this means that the load applied to an element is acting on one of its end nodes.
func (load *ConcentratedLoad) IsNodal() bool {
	return load.T.IsMin() || load.T.IsMax()
}

// AsTorsor returns a vector with the components of the load.
func (load *ConcentratedLoad) AsTorsor() *math.Torsor {
	return math.MakeTorsor(load.LocalFx(), load.LocalFy(), load.LocalMz())
}

// AsTorsorProjectedTo returns the concentrated load vector projected in a reference frame.
func (load *ConcentratedLoad) AsTorsorProjectedTo(refFrame *g2d.RefFrame) *math.Torsor {
	return load.AsTorsor().ProjectedTo(refFrame)
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

// Equals tests whether the two loads are equal or not.
func (load *ConcentratedLoad) Equals(other *ConcentratedLoad) bool {
	return load.Term == other.Term &&
		load.IsInLocalCoords == other.IsInLocalCoords &&
		load.T.Equals(other.T) &&
		nums.FloatsEqual(load.Value, other.Value)
}

func ConcentratedLoadsEqual(a, b []*ConcentratedLoad) bool {
	if len(a) != len(b) {
		return false
	}

	for i, load := range a {
		if !load.Equals(b[i]) {
			return false
		}
	}

	return true
}
