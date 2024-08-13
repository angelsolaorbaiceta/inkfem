package math

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// NilTorsor is a torsor with all components to zero.
var NilTorsor = MakeNilTorsor()

// A Torsor is a force-moment tuple: {fx, fy, mz}, where {fx, fy} are the x and y components
// of the force and mz is the moment around z.
//
// The moment can be represented as the result of a pair of parallel forces separated
// a given distance.
type Torsor struct {
	fx, fy, mz float64
}

var nilTorsor = &Torsor{0.0, 0.0, 0.0}

// MakeTorsor creates a new torsor with the given components.
func MakeTorsor(fx, fy, mz float64) *Torsor {
	return &Torsor{fx, fy, mz}
}

// MakeNilTorsor creates a new torsor with all components to zero.
func MakeNilTorsor() *Torsor {
	return nilTorsor
}

// Fx returns the force x component.
func (torsor *Torsor) Fx() float64 {
	return torsor.fx
}

// Fy returns the force y component.
func (torsor *Torsor) Fy() float64 {
	return torsor.fy
}

// Mz returns the moment around z component.
func (torsor *Torsor) Mz() float64 {
	return torsor.mz
}

// Plus creates a new torsor result of adding this and another.
func (augend *Torsor) Plus(addend *Torsor) *Torsor {
	return MakeTorsor(
		augend.fx+addend.fx,
		augend.fy+addend.fy,
		augend.mz+addend.mz,
	)
}

// PlusComponents creates a new torsor result of adding the passed in torsor components to this one.
func (augend *Torsor) PlusComponents(fx, fy, mz float64) *Torsor {
	return MakeTorsor(
		augend.fx+fx,
		augend.fy+fy,
		augend.mz+mz,
	)
}

// Minus creates a new torsor result of subtracting another to this one.
func (minuend *Torsor) Minus(subtrahend *Torsor) *Torsor {
	return MakeTorsor(
		minuend.fx-subtrahend.fx,
		minuend.fy-subtrahend.fy,
		minuend.mz-subtrahend.mz,
	)
}

// ProjectedToGlobal creates a new torsor with the values projected to the global reference frame
// assuming it was originally projected in the passed in reference frame.
func (torsor *Torsor) ProjectedToGlobal(refFrame *g2d.RefFrame) *Torsor {
	projectedForces := refFrame.ProjectionsToGlobal(torsor.fx, torsor.fy)
	return MakeTorsor(projectedForces.X(), projectedForces.Y(), torsor.mz)
}

// ProjectedTo creates a new torsor with the values projected to the passed in reference frame
// assuming it was originally projected in the global reference frame.
func (torsor *Torsor) ProjectedTo(refFrame *g2d.RefFrame) *Torsor {
	projectedForces := refFrame.ProjectProjections(torsor.fx, torsor.fy)
	return MakeTorsor(projectedForces.X(), projectedForces.Y(), torsor.mz)
}

// Equals checks whether this an another torsors are equal.
func (torsor *Torsor) Equals(other *Torsor) bool {
	return nums.FloatsEqual(torsor.fx, other.fx) &&
		nums.FloatsEqual(torsor.fy, other.fy) &&
		nums.FloatsEqual(torsor.mz, other.mz)
}

// String representation of the torsor.
func (torsor *Torsor) String() string {
	return fmt.Sprintf("{%f %f %f}", torsor.fx, torsor.fy, torsor.mz)
}
