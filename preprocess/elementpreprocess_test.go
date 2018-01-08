package preprocess

import (
	// "fmt"

	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath"
)

/* Axial Member */
func TestSliceAxialMemberNodePositions(t *testing.T) {
	element := makeElementWithLoads(make([]load.Load, 0))
	slicedEl := sliceAxialElement(element)

	if len(slicedEl.Nodes) != 2 {
		t.Error("Expected element to have two nodes")
	}

	if !slicedEl.Nodes[0].Position.Equals(element.StartPoint()) {
		t.Error("First node's position was not as expected")
	}
	if !slicedEl.Nodes[1].Position.Equals(element.EndPoint()) {
		t.Error("Last node's position was not as expected")
	}
}

func TestSliceAxialMemberStartNodeLoads(t *testing.T) {
	loads := []load.Load{
		load.MakeConcentrated(load.FX, true, inkgeom.MIN_T, 50.0),
		load.MakeConcentrated(load.FY, true, inkgeom.MIN_T, 75.0)}
	element := makeElementWithLoads(loads)
	slicedEl := sliceAxialElement(element)

	if slicedEl.Nodes[0].LocalFx() != 50.0 {
		t.Error("Node Fx value not as expected")
	}
	if slicedEl.Nodes[0].LocalFy() != 75.0 {
		t.Error("Node Fy value not as expected")
	}
}

func TestSliceAxialMemberEndNodeLoads(t *testing.T) {
	loads := []load.Load{
		load.MakeConcentrated(load.FX, true, inkgeom.MAX_T, 50.0),
		load.MakeConcentrated(load.FY, true, inkgeom.MAX_T, 75.0)}
	element := makeElementWithLoads(loads)
	slicedEl := sliceAxialElement(element)

	if slicedEl.Nodes[1].LocalFx() != 50.0 {
		t.Error("Node Fx value not as expected")
	}
	if slicedEl.Nodes[1].LocalFy() != 75.0 {
		t.Error("Node Fy value not as expected")
	}
}

func TestSliceAxialMemberGlobalLoadProjected(t *testing.T) {
	loads := []load.Load{load.MakeConcentrated(load.FY, false, inkgeom.MIN_T, 100.0)}
	element := makeElementWithLoads(loads)
	slicedEl := sliceAxialElement(element)
	expectedProjLoadX := inkgeom.MakeVector(0, 100).DotTimes(element.Geometry.DirectionVersor())
	expectedProjLoadY := inkgeom.MakeVector(0, 100).DotTimes(element.Geometry.NormalVersor())

	if !inkmath.FuzzyEqual(slicedEl.Nodes[0].LocalFx(), expectedProjLoadX) {
		t.Error("Node projected Fx value was not as expected")
	}
	if !inkmath.FuzzyEqual(slicedEl.Nodes[0].LocalFy(), expectedProjLoadY) {
		t.Error("Node projected Fy value was not as expected")
	}
}

/* Non Axial Member */
func TestSliceNonAxialMemberNodePositions(t *testing.T) {
	element := makeElementWithLoads(make([]load.Load, 0))
	slicedEl := sliceElement(element, 2)

	if len(slicedEl.Nodes) != 3 {
		t.Error("Expected element to have three nodes")
	}

	if !slicedEl.Nodes[0].Position.Equals(element.StartPoint()) {
		t.Error("First node's position was not as expected")
	}
	if !slicedEl.Nodes[1].Position.Equals(element.PointAt(inkgeom.MakeTParam(0.5))) {
		t.Error("Middle node's position was not as expected")
	}
	if !slicedEl.Nodes[2].Position.Equals(element.EndPoint()) {
		t.Error("Last node's position was not as expected")
	}
}

/* Utils */
func makeElementWithLoads(loads []load.Load) structure.Element {
	return structure.MakeElement(
		1,
		structure.MakeFreeNodeFromProjs(1, 1.0, 2.0),
		structure.MakeFreeNodeFromProjs(2, 3.0, 4.0),
		structure.MakeDispConstraint(),
		structure.MakeDispConstraint(),
		structure.MakeUnitMaterial(),
		structure.MakeUnitSection(),
		loads)
}
