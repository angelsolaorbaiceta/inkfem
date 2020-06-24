package preprocess

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/nums"
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
		load.MakeConcentrated(load.FX, true, inkgeom.MinT, 50.0),
		load.MakeConcentrated(load.FY, true, inkgeom.MinT, 75.0)}
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
		load.MakeConcentrated(load.FX, true, inkgeom.MaxT, 50.0),
		load.MakeConcentrated(load.FY, true, inkgeom.MaxT, 75.0)}
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
	loads := []load.Load{load.MakeConcentrated(load.FY, false, inkgeom.MinT, 100.0)}
	element := makeElementWithLoads(loads)
	slicedEl := sliceAxialElement(element)
	expectedProjLoadX := inkgeom.MakeVector(0, 100).DotTimes(element.Geometry.DirectionVersor())
	expectedProjLoadY := inkgeom.MakeVector(0, 100).DotTimes(element.Geometry.NormalVersor())

	if !nums.FuzzyEqual(slicedEl.Nodes[0].LocalFx(), expectedProjLoadX) {
		t.Error("Node projected Fx value was not as expected")
	}
	if !nums.FuzzyEqual(slicedEl.Nodes[0].LocalFy(), expectedProjLoadY) {
		t.Error("Node projected Fy value was not as expected")
	}
}

/* Non Axial Member : Unloaded */
func TestSliceNonAxialUnloadedMemberNodePositions(t *testing.T) {
	element := makeElementWithLoads(make([]load.Load, 0))
	slicedEl := sliceUnloadedElement(element, 2)

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

/* Non Axial Member : Loaded -> slicing */
func TestDistributedLoadInEntireLengthAddsNoPositions(t *testing.T) {
	loads := []load.Load{load.MakeDistributed(load.FY, true, inkgeom.MinT, 45.0, inkgeom.MaxT, 55.0)}
	tPos := sliceLoadedElementPositions(loads, 2)

	if posCount := len(tPos); posCount != 3 {
		t.Errorf("Expected 3 positions, got %d", posCount)
	}
}

func TestConcentratedLoadAddsPosition(t *testing.T) {
	loads := []load.Load{load.MakeConcentrated(load.FY, true, inkgeom.MakeTParam(0.75), 45.0)}
	tPos := sliceLoadedElementPositions(loads, 2)

	if posCount := len(tPos); posCount != 4 {
		t.Errorf("Expected 4 positions, got %d", posCount)
	}
	if loadPos := tPos[2]; loadPos.Value() != 0.75 {
		t.Errorf("Expected load position to be at 0.75, fount it at %f", loadPos)
	}
}

func TestDistributedLoadAddsTwoPositions(t *testing.T) {
	loads := []load.Load{load.MakeDistributed(load.FY, true, inkgeom.MakeTParam(0.25), 45.0, inkgeom.MakeTParam(0.75), 55.0)}
	tPos := sliceLoadedElementPositions(loads, 2)

	if posCount := len(tPos); posCount != 5 {
		t.Errorf("Expected 5 positions, got %d", posCount)
	}
	if loadPos := tPos[1]; loadPos.Value() != 0.25 {
		t.Errorf("Expected load position to be at 0.25, fount it at %f", loadPos)
	}
	if loadPos := tPos[3]; loadPos.Value() != 0.75 {
		t.Errorf("Expected load position to be at 0.75, fount it at %f", loadPos)
	}
}

func TestMultipleLoadsNotAddingPositionTwice(t *testing.T) {
	loads := []load.Load{
		load.MakeDistributed(load.FY, true, inkgeom.MakeTParam(0.25), 45.0, inkgeom.MakeTParam(0.75), 55.0),
		load.MakeConcentrated(load.FY, true, inkgeom.MakeTParam(0.75), 45.0),
	}
	tPos := sliceLoadedElementPositions(loads, 2)

	if posCount := len(tPos); posCount != 5 {
		t.Errorf("Expected 5 positions, got %d", posCount)
	}
}

/* Non Axial Member : Loaded -> loads */
func TestDistributedLocalLoadDistribution(t *testing.T) {
	element := structure.MakeElement(
		1,
		structure.MakeFreeNodeFromProjs(1, 0.0, 0.0),
		structure.MakeFreeNodeFromProjs(2, 4.0, 0.0),
		structure.MakeDispConstraint(),
		structure.MakeDispConstraint(),
		*structure.MakeUnitMaterial(),
		*structure.MakeUnitSection(),
		[]load.Load{load.MakeDistributed(load.FY, true, inkgeom.MinT, 5.0, inkgeom.MaxT, 5.0)},
	)
	slicedEl := sliceLoadedElement(element, 2)

	// First Node
	if fx := slicedEl.Nodes[0].LocalFx(); fx != 0.0 {
		t.Errorf("First node Fx expected to be 0.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[0].LocalFy(); fy != 5.0 {
		t.Errorf("First node Fy expected to be 5.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[0].LocalMz(); !nums.FuzzyEqual(mz, 5.0/3.0) {
		t.Errorf("First node Mz expected to be %f, but was %f", 5.0/3.0, mz)
	}

	// Second Node
	if fx := slicedEl.Nodes[1].LocalFx(); fx != 0.0 {
		t.Errorf("Second node Fx expected to be 0.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[1].LocalFy(); fy != 10.0 {
		t.Errorf("Second node Fy expected to be 10.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[1].LocalMz(); mz != 0.0 {
		t.Errorf("Second node Mz expected to be 0.0, but was %f", mz)
	}

	// Third Node
	if fx := slicedEl.Nodes[2].LocalFx(); fx != 0.0 {
		t.Errorf("Third node Fx expected to be 0.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[2].LocalFy(); fy != 5.0 {
		t.Errorf("Third node Fy expected to be 5.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[2].LocalMz(); !nums.FuzzyEqual(mz, -5.0/3.0) {
		t.Errorf("Third node Mz expected to be %f, but was %f", -5.0/3.0, mz)
	}
}

func TestDistributedGlobalLoadDistribution(t *testing.T) {
	element := structure.MakeElement(
		1,
		structure.MakeFreeNodeFromProjs(1, 0.0, 0.0),
		structure.MakeFreeNodeFromProjs(2, 4.0, 4.0),
		structure.MakeDispConstraint(),
		structure.MakeDispConstraint(),
		*structure.MakeUnitMaterial(),
		*structure.MakeUnitSection(),
		[]load.Load{load.MakeDistributed(load.FY, false, inkgeom.MinT, 5.0, inkgeom.MaxT, 5.0)},
	)
	slicedEl := sliceLoadedElement(element, 2)

	// First Node
	if fx := slicedEl.Nodes[0].LocalFx(); !nums.FuzzyEqual(fx, 5.0) {
		t.Errorf("First node Fx expected to be 5.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[0].LocalFy(); !nums.FuzzyEqual(fy, 5.0) {
		t.Errorf("First node Fy expected to be 5.0, but was %f", fy)
	}
	if mz, expected := slicedEl.Nodes[0].LocalMz(), 10.0/math.Sqrt(18.0); !nums.FuzzyEqual(mz, expected) {
		t.Errorf("First node Mz expected to be %f, but was %f", expected, mz)
	}

	// Second Node
	if fx := slicedEl.Nodes[1].LocalFx(); !nums.FuzzyEqual(fx, 10.0) {
		t.Errorf("Second node Fx expected to be 10.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[1].LocalFy(); !nums.FuzzyEqual(fy, 10.0) {
		t.Errorf("Second node Fy expected to be 10.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[1].LocalMz(); mz != 0.0 {
		t.Errorf("Second node Mz expected to be 0.0, but was %f", mz)
	}

	// Third Node
	if fx := slicedEl.Nodes[2].LocalFx(); !nums.FuzzyEqual(fx, 5.0) {
		t.Errorf("Third node Fx expected to be 5.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[2].LocalFy(); !nums.FuzzyEqual(fy, 5.0) {
		t.Errorf("Third node Fy expected to be 5.0, but was %f", fy)
	}
	if mz, expected := slicedEl.Nodes[2].LocalMz(), -10.0/math.Sqrt(18.0); !nums.FuzzyEqual(mz, expected) {
		t.Errorf("Third node Mz expected to be %f, but was %f", expected, mz)
	}
}

func TestConcentratedLocalLoadDistribution(t *testing.T) {
	element := structure.MakeElement(
		1,
		structure.MakeFreeNodeFromProjs(1, 0.0, 0.0),
		structure.MakeFreeNodeFromProjs(2, 4.0, 0.0),
		structure.MakeDispConstraint(),
		structure.MakeDispConstraint(),
		*structure.MakeUnitMaterial(),
		*structure.MakeUnitSection(),
		[]load.Load{
			load.MakeConcentrated(load.FX, true, inkgeom.MakeTParam(0.25), 3.0),
			load.MakeConcentrated(load.FY, true, inkgeom.MakeTParam(0.25), 5.0),
			load.MakeConcentrated(load.MZ, true, inkgeom.MakeTParam(0.25), 7.0),
		},
	)
	slicedEl := sliceLoadedElement(element, 2)

	if fx := slicedEl.Nodes[1].LocalFx(); fx != 3.0 {
		t.Errorf("First node Fx expected to be 3.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[1].LocalFy(); fy != 5.0 {
		t.Errorf("First node Fy expected to be 5.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[1].LocalMz(); mz != 7.0 {
		t.Errorf("First node Mz expected to be 7.0, but was %f", mz)
	}
}

func TestConcentratedGlobalLoadDistribution(t *testing.T) {
	element := structure.MakeElement(
		1,
		structure.MakeFreeNodeFromProjs(1, 0.0, 0.0),
		structure.MakeFreeNodeFromProjs(2, 4.0, 4.0),
		structure.MakeDispConstraint(),
		structure.MakeDispConstraint(),
		*structure.MakeUnitMaterial(),
		*structure.MakeUnitSection(),
		[]load.Load{
			load.MakeConcentrated(load.FY, false, inkgeom.MakeTParam(0.25), 5.0),
		},
	)
	slicedEl := sliceLoadedElement(element, 2)

	if fx := slicedEl.Nodes[1].LocalFx(); !nums.FuzzyEqual(fx, 5.0/math.Sqrt2) {
		t.Errorf("First node Fx expected to be %f, but was %f", 5.0/math.Sqrt2, fx)
	}
	if fy := slicedEl.Nodes[1].LocalFy(); !nums.FuzzyEqual(fy, 5.0/math.Sqrt2) {
		t.Errorf("First node Fy expected to be %f, but was %f", 5.0/math.Sqrt2, fy)
	}
	if mz := slicedEl.Nodes[1].LocalMz(); mz != 0.0 {
		t.Errorf("First node Mz expected to be 0.0, but was %f", mz)
	}
}

/* Utils */
func makeElementWithLoads(loads []load.Load) *structure.Element {
	return structure.MakeElement(
		1,
		structure.MakeFreeNodeFromProjs(1, 1.0, 2.0),
		structure.MakeFreeNodeFromProjs(2, 3.0, 4.0),
		structure.MakeDispConstraint(),
		structure.MakeDispConstraint(),
		*structure.MakeUnitMaterial(),
		*structure.MakeUnitSection(),
		loads)
}
