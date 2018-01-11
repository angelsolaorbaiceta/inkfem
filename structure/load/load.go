/*
Package load contains definition of loads applied to structural members.
*/
package load

import (
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath"
)

/*
Load represents a distributed or concentrated load definition. A load is expressed as:
    - a term of application, which in 2D can be: Force in X, Force in Y and Moment about Z
    - a projection frame, which can be local to the element to which load is applied or global
    - start/end position and value
*/
type Load struct {
	Term            Term
	IsInLocalCoords bool

	startT, endT         inkgeom.TParam
	startValue, endValue float64
}

/* Creation */

/*
MakeDistributed creates a distributed load for the given term (FX, FY, MZ) which may be defined
locally to the element it will be applied to or referenced in global coordinates.
Distributed loads are defined by a start position - value and an end position - value
tuples.
*/
func MakeDistributed(
	term Term,
	isInLocalCoords bool,
	startT inkgeom.TParam /* -> */, startValue float64,
	endT inkgeom.TParam /* -> */, endValue float64) Load {
	return Load{term, isInLocalCoords, startT, endT, startValue, endValue}
}

/*
MakeConcentrated creates a concentrated load for the given term (FX, FY, MZ) which may be defined
locally to the element it will be applied to or referenced in global coordinates.
Concentrated loads are defined by a position - value tuple.
*/
func MakeConcentrated(
	term Term,
	isInLocalCoords bool,
	t inkgeom.TParam,
	value float64) Load {
	return Load{term, isInLocalCoords, t, t, value, value}
}

/* Properties */

// IsConcentrated returns true if the load is concentrated.
func (load Load) IsConcentrated() bool {
	return load.startT.Equals(load.endT)
}

// IsDistributed returns true if the load is distributed.
func (load Load) IsDistributed() bool {
	return !load.IsConcentrated()
}

// IsNodal returns true if the load is concentrated and applied in extreme values of T.
func (load Load) IsNodal() bool {
	return load.IsConcentrated() && (load.T().IsMin() || load.T().IsMax())
}

// T returns the t parameter value for a concentrated load. Panics if the load is distributed.
func (load Load) T() inkgeom.TParam {
	if load.IsDistributed() {
		panic("Can't get T value for distributed load. Use StartT and EndT instead")
	}

	return load.startT
}

// Value returns the value for a concentrated load. Panics if the load is distributed.
func (load Load) Value() float64 {
	if load.IsDistributed() {
		panic("Can't get T value for distributed load. Use StartT and EndT instead")
	}

	return load.startValue
}

// VectorValue returns a vector for a concentrated load with the components of the load.
func (load Load) VectorValue() [3]float64 {
	switch load.Term {
	case FX:
		return [3]float64{load.Value(), 0.0, 0.0}
	case FY:
		return [3]float64{0.0, load.Value(), 0.0}
	case MZ:
		return [3]float64{0.0, 0.0, load.Value()}
	default:
		panic("Unknown load term: " + load.Term)
	}
}

// ForcesVector returns a vector for a concentrated load with the components of {Fx, Fy}.
func (load Load) ForcesVector() inkgeom.Projectable {
	switch load.Term {
	case FX:
		return inkgeom.MakeVector(load.Value(), 0.0)
	case FY:
		return inkgeom.MakeVector(0.0, load.Value())
	default:
		return inkgeom.MakeVector(0.0, 0.0)
	}
}

// ValueAt returns the value of the distributed load at a given t Parameter value.
func (load Load) ValueAt(t inkgeom.TParam) float64 {
	if load.IsConcentrated() {
		panic("Can't get value at of a concentrated load at a given point")
	}

	if t.IsLessThan(load.StartT()) || t.IsGreaterThan(load.EndT()) {
		return 0.0
	}

	return inkmath.LinInterpol(load.startT.Value(), load.startValue, load.endT.Value(), load.endValue, t.Value())
}

// AvgValueBetween return the average load value inside the given range for the distributed load.
func (load Load) AvgValueBetween(startT, endT inkgeom.TParam) float64 {
	ok, maxStartT, minEndT := inkmath.RangesOverlap(load.startT.Value(), load.endT.Value(), startT.Value(), endT.Value())
	if !ok {
		return 0.0
	}

	startVal := load.ValueAt(inkgeom.MakeTParam(maxStartT))
	endVal := load.ValueAt(inkgeom.MakeTParam(minEndT))
	applicationLength := minEndT - maxStartT
	rangeLength := startT.DistanceTo(endT)

	// one of both ends has a zero value load -> No need to average values
	if inkmath.IsCloseToZero(startVal * endVal) {
		return (startVal + endVal) * applicationLength / rangeLength
	}

	return applicationLength * 0.5 * (startVal + endVal) / rangeLength
}

// StartT returns the start T parameter value for distributed loads.
func (load Load) StartT() inkgeom.TParam {
	return load.startT
}

// EndT returns the end T parameter value for the distributed load.
func (load Load) EndT() inkgeom.TParam {
	return load.endT
}
