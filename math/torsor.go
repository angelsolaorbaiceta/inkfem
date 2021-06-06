package math

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
)

/*
A Torsor is a force-moment tuple: {fx, fy, mz}, where {fx, fy} are the x and y components
of the force and mz is the moment around z.

The moment can be represented as the result of a pair of parallel forces separated
a given distance.
*/
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

// Minus creates a new torsor result of subtracting another to this one.
func (minuend *Torsor) Minus(subtrahend *Torsor) *Torsor {
	return MakeTorsor(
		minuend.fx-subtrahend.fx,
		minuend.fy-subtrahend.fy,
		minuend.mz-subtrahend.mz,
	)
}

// Equals checks whether this an another torsors are equal.
func (torsor *Torsor) Equals(other *Torsor) bool {
	return inkgeom.FloatsEqual(torsor.fx, other.fx) &&
		inkgeom.FloatsEqual(torsor.fy, other.fy) &&
		inkgeom.FloatsEqual(torsor.mz, other.mz)
}

func (torsor *Torsor) String() string {
	return fmt.Sprintf("%f %f %f", torsor.fx, torsor.fy, torsor.mz)
}
