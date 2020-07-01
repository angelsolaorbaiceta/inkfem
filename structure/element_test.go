package structure

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/nums"
)

const (
	elementID   = 11
	startNodeID = 21
	endNodeID   = 22
)

var (
	startPoint = inkgeom.MakePoint(0, 0)
	startNode  = MakeNode(startNodeID, startPoint, nilConstraint)

	endPoint = inkgeom.MakePoint(100, 0)
	endNode  = MakeNode(endNodeID, endPoint, nilConstraint)

	material = &Material{"material", 2, 3, 4, 5, 6, 7}
	section  = &Section{"section", 2, 3, 4, 5, 6}
)

func TestElementStartPoint(t *testing.T) {
	element := makeElement()

	if point := element.StartPoint(); !point.Equals(startPoint) {
		t.Errorf("Wrong start point: %v", point)
	}
}

func TestElementEndPoint(t *testing.T) {
	element := makeElement()

	if point := element.EndPoint(); !point.Equals(endPoint) {
		t.Errorf("Wrong end point: %v", point)
	}
}

func TestElementHasLoadsApplied(t *testing.T) {
	t.Run("has loads applied", func(t *testing.T) {
		l := load.MakeConcentrated(load.FX, true, inkgeom.MinT, 10)
		element := makeLoadedElement(l)

		if !element.HasLoadsApplied() {
			t.Error("Element has loads applied")
		}
	})

	t.Run("has no loads applied", func(t *testing.T) {
		element := makeElement()

		if element.HasLoadsApplied() {
			t.Error("Element doesn't have loads applied")
		}
	})
}

func TestElementIsAxial(t *testing.T) {
	t.Run("isn't axial if start link allows rotation", func(t *testing.T) {
		element := makeElement()
		element.StartLink = dispConstraint
		element.EndLink = fullConstraint

		if element.IsAxialMember() {
			t.Error("Element shouln't be axial")
		}
	})

	t.Run("isn't axial if end link allows rotation", func(t *testing.T) {
		element := makeElement()
		element.StartLink = fullConstraint
		element.EndLink = dispConstraint

		if element.IsAxialMember() {
			t.Error("Element shouln't be axial")
		}
	})

	t.Run("isn't axial if has at least a distributed load", func(t *testing.T) {
		l := load.MakeDistributed(load.FX, true, inkgeom.MinT, 20, inkgeom.MaxT, 40)
		element := makeLoadedElement(l)
		element.StartLink = dispConstraint
		element.EndLink = dispConstraint

		if element.IsAxialMember() {
			t.Error("Element shouln't be axial")
		}
	})

	t.Run("isn't axial if has at least a concentrated non-nodal load", func(t *testing.T) {
		l := load.MakeConcentrated(load.MZ, true, inkgeom.HalfT, 10)
		element := makeLoadedElement(l)
		element.StartLink = dispConstraint
		element.EndLink = dispConstraint

		if element.IsAxialMember() {
			t.Error("Element shouln't be axial")
		}
	})

	t.Run("isn't axial if has at least a nodal MZ load", func(t *testing.T) {
		l := load.MakeConcentrated(load.MZ, true, inkgeom.MinT, 10)
		element := makeLoadedElement(l)
		element.StartLink = dispConstraint
		element.EndLink = dispConstraint

		if element.IsAxialMember() {
			t.Error("Element shouln't be axial")
		}
	})

	t.Run("is axial if pinned and all loads are nodal and not MZ", func(t *testing.T) {
		l := load.MakeConcentrated(load.FY, true, inkgeom.MinT, 10)
		element := makeLoadedElement(l)
		element.StartLink = dispConstraint
		element.EndLink = dispConstraint

		if !element.IsAxialMember() {
			t.Error("Element should be axial")
		}
	})
}

func TestHorizontalElementGlobalStiffnessMatrix(t *testing.T) {
	var (
		element = makeElement()
		matrix  = element.StiffnessGlobalMat(inkgeom.MinT, inkgeom.MaxT)
		e       = material.YoungMod
		i       = section.IStrong
		a       = section.Area
		l       = element.Geometry.Length()
	)

	t.Run("Fx -> Dx terms", func(t *testing.T) {
		want := e * a / l

		if got := matrix.Value(0, 0); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}

		if got := matrix.Value(0, 3); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}

		if got := matrix.Value(3, 0); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}

		if got := matrix.Value(3, 3); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}
	})

	t.Run("Fy -> Dy terms", func(t *testing.T) {
		want := (12.0 * e * i) / (l * l * l)

		if got := matrix.Value(1, 1); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}

		if got := matrix.Value(1, 4); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}

		if got := matrix.Value(4, 1); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}

		if got := matrix.Value(4, 4); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}
	})

	t.Run("Fy -> Rz terms", func(t *testing.T) {
		want := (6.0 * e * i) / (l * l)

		if got := matrix.Value(1, 2); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}

		if got := matrix.Value(1, 5); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}

		if got := matrix.Value(4, 2); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}

		if got := matrix.Value(4, 5); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}
	})

	t.Run("Mz -> Dy terms", func(t *testing.T) {
		want := (6.0 * e * i) / (l * l)

		if got := matrix.Value(2, 1); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}

		if got := matrix.Value(2, 4); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}

		if got := matrix.Value(5, 1); !nums.FuzzyEqual(want, got) {
			t.Errorf("Expected term to be %f, but got %f", want, got)
		}

		if got := matrix.Value(5, 4); !nums.FuzzyEqual(-want, got) {
			t.Errorf("Expected term to be %f, but got %f", -want, got)
		}
	})

	t.Run("Mz -> Rz terms", func(t *testing.T) {
		wantOne := 2.0 * e * i / l
		wantTwo := 4.0 * e * i / l

		if got := matrix.Value(2, 2); !nums.FuzzyEqual(wantTwo, got) {
			t.Errorf("Expected term to be %f, but got %f", wantTwo, got)
		}

		if got := matrix.Value(2, 5); !nums.FuzzyEqual(wantOne, got) {
			t.Errorf("Expected term to be %f, but got %f", wantOne, got)
		}

		if got := matrix.Value(5, 2); !nums.FuzzyEqual(wantOne, got) {
			t.Errorf("Expected term to be %f, but got %f", wantOne, got)
		}

		if got := matrix.Value(5, 5); !nums.FuzzyEqual(wantTwo, got) {
			t.Errorf("Expected term to be %f, but got %f", wantTwo, got)
		}
	})
}

func makeElement() *Element {
	return MakeElement(
		elementID,
		startNode, endNode,
		fullConstraint, fullConstraint,
		material, section,
		[]load.Load{},
	)
}

func makeLoadedElement(l load.Load) *Element {
	return MakeElement(
		elementID,
		startNode, endNode,
		fullConstraint, fullConstraint,
		material, section,
		[]load.Load{l},
	)
}
