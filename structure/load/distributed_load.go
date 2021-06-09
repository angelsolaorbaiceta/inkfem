package load

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/nums"
)

/*
Load is a distributed or concentrated load.

Distributed loads are linear: the start and end values are interpolated linearly.

A load is expressed as:
	- a term of application, which in 2D can be: Force in X, Force in Y or Moment about Z
	- a projection frame, which can be local to the element to which load is applied or global
  - start/end position and value
*/
type DistributedLoad struct {
	Term                 Term
	IsInLocalCoords      bool
	StartT, EndT         inkgeom.TParam
	StartValue, EndValue float64
}

/*
MakeDistributed creates a distributed load for the given term (FX, FY, MZ) which may be defined
locally to the element it will be applied to or referenced in global coordinates.

Distributed loads are defined by a start position - value and an end position - value tuples.
*/
func MakeDistributed(
	term Term,
	isInLocalCoords bool,
	startT inkgeom.TParam,
	startValue float64,
	endT inkgeom.TParam,
	endValue float64,
) *DistributedLoad {
	return &DistributedLoad{term, isInLocalCoords, startT, endT, startValue, endValue}
}

// ValueAt returns the value of the load at a given t Parameter value.
func (load *DistributedLoad) ValueAt(t inkgeom.TParam) float64 {
	if t.IsLessThan(load.StartT) || t.IsGreaterThan(load.EndT) {
		return 0.0
	}

	return nums.LinInterpol(
		load.StartT.Value(),
		load.StartValue,
		load.EndT.Value(),
		load.EndValue,
		t.Value(),
	)
}

// AsVectorAt returns the the distributed load vector at a given position.
func (load *DistributedLoad) AsVectorAt(t inkgeom.TParam) *math.Torsor {
	value := load.ValueAt(t)

	switch load.Term {
	case FX:
		return math.MakeTorsor(value, 0.0, 0.0)

	case FY:
		return math.MakeTorsor(0.0, value, 0.0)

	case MZ:
		return math.MakeTorsor(0.0, 0.0, value)

	default:
		panic("Unknown load term: " + load.Term)
	}
}

/*
ProjectedVectorAt returns the distributed load vector at a given position projected in a
reference frame.
*/
func (load *DistributedLoad) ProjectedVectorAt(t inkgeom.TParam, refFrame g2d.RefFrame) *math.Torsor {
	return load.AsVectorAt(t).ProjectedTo(refFrame)
}

// Equals tests whether the two loads are equal or not.
func (load *DistributedLoad) Equals(other *DistributedLoad) bool {
	return load.Term == other.Term &&
		load.IsInLocalCoords == other.IsInLocalCoords &&
		load.StartT.Equals(other.StartT) &&
		nums.FuzzyEqual(load.StartValue, other.StartValue) &&
		load.EndT.Equals(other.EndT) &&
		nums.FuzzyEqual(load.EndValue, other.EndValue)
}
