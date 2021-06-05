package math

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
)

// A Torsor is a force-moment tuple.
type Torsor struct {
	fx, fy, mz float64
}

var nilTorsor = &Torsor{0.0, 0.0, 0.0}

func MakeTorsor(fx, fy, mz float64) *Torsor {
	return &Torsor{fx, fy, mz}
}

func MakeNilTorsor() *Torsor {
	return nilTorsor
}

func (torsor *Torsor) Fx() float64 {
	return torsor.fx
}

func (torsor *Torsor) Fy() float64 {
	return torsor.fy
}

func (torsor *Torsor) Mz() float64 {
	return torsor.mz
}

func (augend *Torsor) Plus(addend *Torsor) *Torsor {
	return MakeTorsor(
		augend.fx+addend.fx,
		augend.fy+addend.fy,
		augend.mz+addend.mz,
	)
}

func (minuend *Torsor) Minus(subtrahend *Torsor) *Torsor {
	return MakeTorsor(
		minuend.fx-subtrahend.fx,
		minuend.fy-subtrahend.fy,
		minuend.mz-subtrahend.mz,
	)
}

func (torsor *Torsor) Equals(other *Torsor) bool {
	return inkgeom.FloatsEqual(torsor.fx, other.fx) &&
		inkgeom.FloatsEqual(torsor.fy, other.fy) &&
		inkgeom.FloatsEqual(torsor.mz, other.mz)
}

func (torsor *Torsor) String() string {
	return fmt.Sprintf("%f %f %f", torsor.fx, torsor.fy, torsor.mz)
}
