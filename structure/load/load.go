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

/*
Package load contains definition of loads applied to structural members.
*/
package load

import (
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/nums"
)

/*
Load is a distributed or concentrated load.

Distributed loads are linear: the start and end values are interpolated linearly.

A load is expressed as:
	- a term of application, which in 2D can be: Force in X, Force in Y or
	Moment about Z
	- a projection frame, which can be local to the element to which load is
	applied or global
  - start/end position and value
*/
type Load struct {
	Term                 Term
	IsInLocalCoords      bool
	startT, endT         inkgeom.TParam
	startValue, endValue float64
}

/* <-- Creation --> */

/*
MakeDistributed creates a distributed load for the given term (FX, FY, MZ) which
may be defined locally to the element it will be applied to or referenced in
global coordinates.

Distributed loads are defined by a start position - value and an end position -
value tuples.
*/
func MakeDistributed(
	term Term,
	isInLocalCoords bool,
	startT inkgeom.TParam,
	startValue float64,
	endT inkgeom.TParam,
	endValue float64,
) Load {
	return Load{term, isInLocalCoords, startT, endT, startValue, endValue}
}

/*
MakeConcentrated creates a concentrated load for the given term (FX, FY, MZ)
which may be defined locally to the element it will be applied to or referenced
in global coordinates.

Concentrated loads are defined by a position - value tuple.
*/
func MakeConcentrated(
	term Term,
	isInLocalCoords bool,
	t inkgeom.TParam,
	value float64,
) Load {
	return Load{term, isInLocalCoords, t, t, value, value}
}

/* <-- Properties --> */

/*
IsConcentrated returns true if the load is concentrated.
*/
func (load Load) IsConcentrated() bool {
	return load.startT.Equals(load.endT)
}

/*
IsDistributed returns true if the load is distributed.
*/
func (load Load) IsDistributed() bool {
	return !load.IsConcentrated()
}

/*
IsNodal returns true if the load is concentrated and applied in extreme values of T.
*/
func (load Load) IsNodal() bool {
	return load.IsConcentrated() && (load.T().IsMin() || load.T().IsMax())
}

/*
T returns the t parameter value for a concentrated load. Panics if the load is
distributed.
*/
func (load Load) T() inkgeom.TParam {
	if load.IsDistributed() {
		panic("Can't get T value for distributed load. Use StartT and EndT instead")
	}

	return load.startT
}

/*
Value returns the value for a concentrated load. Panics if the load is distributed.
*/
func (load Load) Value() float64 {
	if load.IsDistributed() {
		panic("Can't get T value for distributed load. Use StartT and EndT instead")
	}

	return load.startValue
}

/*
VectorValue returns a vector for a concentrated load with the components of the load.
*/
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

/*
ForcesVector returns a vector for a concentrated load with the components of {Fx, Fy}.
*/
func (load Load) ForcesVector() g2d.Projectable {
	switch load.Term {
	case FX:
		return g2d.MakeVector(load.Value(), 0.0)

	case FY:
		return g2d.MakeVector(0.0, load.Value())

	default:
		return g2d.MakeVector(0.0, 0.0)
	}
}

/* <-- Methods --> */

/*
ValueAt returns the value of the distributed load at a given t Parameter value.
*/
func (load Load) ValueAt(t inkgeom.TParam) float64 {
	if load.IsConcentrated() {
		panic("Can't get value at of a concentrated load at a given point")
	}

	if t.IsLessThan(load.startT) || t.IsGreaterThan(load.endT) {
		return 0.0
	}

	return nums.LinInterpol(
		load.startT.Value(),
		load.startValue,
		load.endT.Value(),
		load.endValue,
		t.Value(),
	)
}

/*
VectorValueAt returns the the distributed load vector at a given position.
*/
func (load Load) VectorValueAt(t inkgeom.TParam) [3]float64 {
	value := load.ValueAt(t)

	switch load.Term {
	case FX:
		return [3]float64{value, 0.0, 0.0}

	case FY:
		return [3]float64{0.0, value, 0.0}

	case MZ:
		return [3]float64{0.0, 0.0, value}

	default:
		panic("Unknown load term")
	}
}

/*
ProjectedVectorValueAt returns the distributed load vector at a given position projected
in a reference frame.
*/
func (load Load) ProjectedVectorValueAt(t inkgeom.TParam, refFrame g2d.RefFrame) [3]float64 {
	var (
		vectorValue     = load.VectorValueAt(t)
		projectedVector = refFrame.ProjectVector(g2d.MakeVector(vectorValue[0], vectorValue[1]))
	)

	vectorValue[0] = projectedVector.X
	vectorValue[1] = projectedVector.Y

	return vectorValue
}

/*
StartT returns the start T parameter value for distributed loads.
*/
func (load Load) StartT() inkgeom.TParam {
	return load.startT
}

/*
EndT returns the end T parameter value for the distributed load.
*/
func (load Load) EndT() inkgeom.TParam {
	return load.endT
}

/*
Equals tests whether the two loads are equal or not.
*/
func (load Load) Equals(other Load) bool {
	return load.Term == other.Term &&
		load.IsInLocalCoords == other.IsInLocalCoords &&
		load.startT.Equals(other.startT) &&
		nums.FuzzyEqual(load.startValue, other.startValue) &&
		load.endT.Equals(other.endT) &&
		nums.FuzzyEqual(load.endValue, other.endValue)
}
