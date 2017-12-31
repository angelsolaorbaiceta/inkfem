/*
Definition of loads applied to structural members.
*/
package load

import (
    "github.com/angelsolaorbaiceta/inkgeom"
)

// Distributed or Concentrated load definition. A load is expressed as:
//      - a term of application, which in 2D can be: Force in X, Force in Y and Moment about Z
//      - a projection frame, which can be local to the element to which load is applied or global
//      - start/end position and value
type Load struct {
    Term LoadTerm
    IsInLocalCoords bool

    startT, endT inkgeom.TParam
    startValue, endValue float64
}

/* Creation */

// Creates a distributed load for the given term (FX, FY, MZ) which may be defined
// locally to the element it will be applied to or referenced in global coordinates.
// Distributed loads are defined by a start position - value and an end position - value
// tuples.
func MakeDistributed(
    term LoadTerm,
    isInLocalCoords bool,
    startT inkgeom.TParam, /* -> */ startValue float64,
    endT inkgeom.TParam, /* -> */ endValue float64) Load {
    return Load{term, isInLocalCoords, startT, endT, startValue, endValue}
}

// Creates a concentrated load for the given term (FX, FY, MZ) which may be defined
// locally to the element it will be applied to or referenced in global coordinates.
// Concentrated loads are defined by a position - value tuple.
func MakeConcentrated(
    term LoadTerm,
    isInLocalCoords bool,
    t inkgeom.TParam,
    value float64) Load {
    return Load{term, isInLocalCoords, t, t, value, value}
}

/* Properties */
func (load Load) IsConcentrated() bool {
    return load.startT.Equals(load.endT)
}

func (load Load) IsDistributed() bool {
    return !load.IsConcentrated()
}

func (load Load) IsNodal() bool {
    return load.IsConcentrated() && (load.T().IsMin() || load.T().IsMax())
}

func (load Load) T() inkgeom.TParam {
    if load.IsDistributed() {
        panic("Can't get T value for distributed load. Use StartT and EndT instead")
    }

    return load.startT
}

func (load Load) StartT() inkgeom.TParam {
    return load.startT
}

func (load Load) EndT() inkgeom.TParam {
    return load.endT
}
