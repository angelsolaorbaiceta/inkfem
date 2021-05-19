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
A ConcentratedLoad is a load applied in a specific point.
*/
type ConcentratedLoad struct {
	Term            Term
	IsInLocalCoords bool
	T               inkgeom.TParam
	Value           float64
}

/*
MakeConcentrated creates a concentrated load for the given term (FX, FY or MZ)
which may be defined locally to the element it will be applied to or referenced
in global coordinates.

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
*/
func (load *ConcentratedLoad) IsNodal() bool {
	return load.T.IsMin() || load.T.IsMax()
}

/*
AsVector returns a vector with the components of the load.
*/
func (load *ConcentratedLoad) AsVector() [3]float64 {
	switch load.Term {
	case FX:
		return [3]float64{load.Value, 0.0, 0.0}

	case FY:
		return [3]float64{0.0, load.Value, 0.0}

	case MZ:
		return [3]float64{0.0, 0.0, load.Value}

	default:
		panic("Unknown load term: " + load.Term)
	}
}

/*
ForcesVector returns a vector for a concentrated load with the components of {Fx, Fy}.
*/
func (load *ConcentratedLoad) ForcesVector() g2d.Projectable {
	switch load.Term {
	case FX:
		return g2d.MakeVector(load.Value, 0.0)

	case FY:
		return g2d.MakeVector(0.0, load.Value)

	case MZ:
		return g2d.MakeVector(0.0, 0.0)

	default:
		panic("Unknown load term: " + load.Term)
	}
}

/*
Equals tests whether the two loads are equal or not.
*/
func (load *ConcentratedLoad) Equals(other *ConcentratedLoad) bool {
	return load.Term == other.Term &&
		load.IsInLocalCoords == other.IsInLocalCoords &&
		load.T.Equals(other.T) &&
		nums.FuzzyEqual(load.Value, other.Value)
}
