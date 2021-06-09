package math

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestTorsorComponents(t *testing.T) {
	torsor := MakeTorsor(1, 2, 3)

	if gotFx := torsor.Fx(); gotFx != 1 {
		t.Errorf("Expected Fx = 1, but got %f", gotFx)
	}
	if gotFy := torsor.Fy(); gotFy != 2 {
		t.Errorf("Expected Fy = 2, but got %f", gotFy)
	}
	if gotMz := torsor.Mz(); gotMz != 3 {
		t.Errorf("Expected Mz = 3, but got %f", gotMz)
	}
}

func TestTorsorOperations(t *testing.T) {
	var (
		t1 = MakeTorsor(1, 2, 3)
		t2 = MakeTorsor(4, 8, 12)
	)

	t.Run("addition", func(t *testing.T) {
		expected := MakeTorsor(
			t1.Fx()+t2.Fx(),
			t1.Fy()+t2.Fy(),
			t1.Mz()+t2.Mz(),
		)

		if got := t1.Plus(t2); !got.Equals(expected) {
			t.Errorf("Expected %v, but got %v", expected, got)
		}
	})

	t.Run("subtraction", func(t *testing.T) {
		expected := MakeTorsor(
			t1.Fx()-t2.Fx(),
			t1.Fy()-t2.Fy(),
			t1.Mz()-t2.Mz(),
		)

		if got := t1.Minus(t2); !got.Equals(expected) {
			t.Errorf("Expected %v, but got %v", expected, got)
		}
	})
}

func TestTorsorProjection(t *testing.T) {
	var (
		refFrame = g2d.MakeRefFrameWithIVersor(g2d.MakeVersor(1, 1))
		torsor   = MakeTorsor(10, 20, 50)
	)

	t.Run("project from local to global", func(t *testing.T) {
		expected := MakeTorsor(
			-5.0*math.Sqrt2,
			15.0*math.Sqrt2,
			50.0,
		)

		if got := torsor.ProjectedToGlobal(refFrame); !got.Equals(expected) {
			t.Errorf("Expected %v, but got %v", expected, got)
		}
	})

	t.Run("project from global to local", func(t *testing.T) {
		expected := MakeTorsor(
			15.0*math.Sqrt2,
			5.0*math.Sqrt2,
			50.0,
		)

		if got := torsor.ProjectedTo(refFrame); !got.Equals(expected) {
			t.Errorf("Expected %v, but got %v", expected, got)
		}
	})
}
