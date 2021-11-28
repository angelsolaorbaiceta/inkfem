package preprocess

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// TODO: separate into shorter and more focused test files

/* <-- Axial Member --> */

func TestSliceAxialMemberNodePositions(t *testing.T) {
	var (
		element  = makeElementWithoutLoads()
		slicedEl = sliceAxialElement(element)
	)

	if slicedEl.NodesCount() != 2 {
		t.Error("Expected element to have two nodes")
	}

	if pos := slicedEl.Nodes[0].Position; !pos.Equals(element.StartPoint()) {
		t.Errorf("Wrong node position. Want %v, got %v", element.StartPoint(), pos)
	}
	if pos := slicedEl.Nodes[1].Position; !pos.Equals(element.EndPoint()) {
		t.Errorf("Wrong node position. Want %v, got %v", element.EndPoint(), pos)
	}
}

func TestSliceAxialMemberNodeLoads(t *testing.T) {
	var (
		loads = []*load.ConcentratedLoad{
			load.MakeConcentrated(load.FX, true, nums.MinT, 50.0),
			load.MakeConcentrated(load.FY, true, nums.MinT, 75.0),
			load.MakeConcentrated(load.FX, true, nums.MaxT, 100.0),
			load.MakeConcentrated(load.FY, true, nums.MaxT, 200.0),
		}
		element  = makeElementWithLoads(loads)
		slicedEl = sliceAxialElement(element)
	)

	t.Run("Start node loads", func(t *testing.T) {
		if slicedEl.Nodes[0].NetLocalFx() != 50.0 {
			t.Error("Node Fx value not as expected")
		}
		if slicedEl.Nodes[0].NetLocalFy() != 75.0 {
			t.Error("Node Fy value not as expected")
		}
	})

	t.Run("End node loads", func(t *testing.T) {
		if slicedEl.Nodes[1].NetLocalFx() != 100.0 {
			t.Error("Node Fx value not as expected")
		}
		if slicedEl.Nodes[1].NetLocalFy() != 200.0 {
			t.Error("Node Fy value not as expected")
		}
	})
}

func TestSliceAxialMemberGlobalLoadProjected(t *testing.T) {
	var (
		loads = []*load.ConcentratedLoad{
			load.MakeConcentrated(load.FY, false, nums.MinT, 100.0),
		}
		element           = makeElementWithLoads(loads)
		slicedEl          = sliceAxialElement(element)
		expectedProjLoadX = g2d.MakeVector(0, 100).DotTimes(element.DirectionVersor())
		expectedProjLoadY = g2d.MakeVector(0, 100).DotTimes(element.NormalVersor())
	)

	if fx := slicedEl.Nodes[0].NetLocalFx(); !nums.FloatsEqual(fx, expectedProjLoadX) {
		t.Error("Node projected Fx value was not as expected")
	}
	if fy := slicedEl.Nodes[0].NetLocalFy(); !nums.FloatsEqual(fy, expectedProjLoadY) {
		t.Error("Node projected Fy value was not as expected")
	}
}

/* <-- Non Axial Member : Unloaded --> */

func TestSliceNonAxialUnloadedMemberNodePositions(t *testing.T) {
	var (
		element  = makeElementWithoutLoads()
		slicedEl = sliceElementWithoutLoads(element, 2)
	)

	if slicedEl.NodesCount() != 3 {
		t.Error("Expected element to have three nodes")
	}

	var (
		wantStartPoint = element.StartPoint()
		wantMidPoint   = element.PointAt(nums.HalfT)
		wantEndPoint   = element.EndPoint()
	)
	if pos := slicedEl.Nodes[0].Position; !pos.Equals(wantStartPoint) {
		t.Error("First node's position was not as expected")
	}
	if pos := slicedEl.Nodes[1].Position; !pos.Equals(wantMidPoint) {
		t.Error("Middle node's position was not as expected")
	}
	if pos := slicedEl.Nodes[2].Position; !pos.Equals(wantEndPoint) {
		t.Error("Last node's position was not as expected")
	}
}

/* <-- Non Axial Member : Loaded -> slicing --> */

func TestDistributedLoadInEntireLengthAddsNoPositions(t *testing.T) {
	loads := []*load.DistributedLoad{
		load.MakeDistributed(load.FY, true, nums.MinT, 45.0, nums.MaxT, 55.0),
	}
	tPos := sliceLoadedElementPositions([]*load.ConcentratedLoad{}, loads, 2)

	if posCount := len(tPos); posCount != 3 {
		t.Errorf("Expected 3 positions, got %d", posCount)
	}
}

func TestConcentratedLoadAddsPosition(t *testing.T) {
	loads := []*load.ConcentratedLoad{
		load.MakeConcentrated(load.FY, true, nums.MakeTParam(0.75), 45.0),
	}
	tPos := sliceLoadedElementPositions(loads, []*load.DistributedLoad{}, 2)

	if posCount := len(tPos); posCount != 4 {
		t.Errorf("Expected 4 positions, got %d", posCount)
	}
	if loadPos := tPos[2]; loadPos.Value() != 0.75 {
		t.Errorf("Expected load position to be at 0.75, fount it at %f", loadPos)
	}
}

func TestDistributedLoadAddsTwoPositions(t *testing.T) {
	loads := []*load.DistributedLoad{
		load.MakeDistributed(load.FY, true, nums.MakeTParam(0.25), 45.0, nums.MakeTParam(0.75), 55.0),
	}
	tPos := sliceLoadedElementPositions([]*load.ConcentratedLoad{}, loads, 2)

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
	var (
		concentratedLoads = []*load.ConcentratedLoad{
			load.MakeConcentrated(load.FY, true, nums.MakeTParam(0.75), 45.0),
		}
		distributedLoads = []*load.DistributedLoad{
			load.MakeDistributed(load.FY, true, nums.MakeTParam(0.25), 45.0, nums.MakeTParam(0.75), 55.0),
		}
		tPos = sliceLoadedElementPositions(concentratedLoads, distributedLoads, 2)
	)

	if posCount := len(tPos); posCount != 5 {
		t.Errorf("Expected 5 positions, got %d", posCount)
	}
}

/* <-- Non Axial Member : Loaded -> loads --> */

func TestDistributedLocalLoadDistribution(t *testing.T) {
	element := structure.MakeElementBuilder(
		"1",
	).WithStartNode(
		structure.MakeFreeNodeAtPosition("1", 0.0, 0.0), &structure.DispConstraint,
	).WithEndNode(
		structure.MakeFreeNodeAtPosition("2", 4.0, 0.0), &structure.DispConstraint,
	).WithSection(
		structure.MakeUnitSection(),
	).WithMaterial(
		structure.MakeUnitMaterial(),
	).AddDistributedLoads(
		[]*load.DistributedLoad{
			load.MakeDistributed(load.FY, true, nums.MinT, 5.0, nums.MaxT, 5.0),
		},
	).Build()

	slicedEl := sliceLoadedElement(element, 2)

	// First Node
	if fx := slicedEl.Nodes[0].NetLocalFx(); fx != 0.0 {
		t.Errorf("First node Fx expected to be 0.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[0].NetLocalFy(); fy != 5.0 {
		t.Errorf("First node Fy expected to be 5.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[0].NetLocalMz(); !nums.FloatsEqual(mz, 5.0/3.0) {
		t.Errorf("First node Mz expected to be %f, but was %f", 5.0/3.0, mz)
	}

	// Second Node
	if fx := slicedEl.Nodes[1].NetLocalFx(); fx != 0.0 {
		t.Errorf("Second node Fx expected to be 0.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[1].NetLocalFy(); fy != 10.0 {
		t.Errorf("Second node Fy expected to be 10.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[1].NetLocalMz(); mz != 0.0 {
		t.Errorf("Second node Mz expected to be 0.0, but was %f", mz)
	}

	// Third Node
	if fx := slicedEl.Nodes[2].NetLocalFx(); fx != 0.0 {
		t.Errorf("Third node Fx expected to be 0.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[2].NetLocalFy(); fy != 5.0 {
		t.Errorf("Third node Fy expected to be 5.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[2].NetLocalMz(); !nums.FloatsEqual(mz, -5.0/3.0) {
		t.Errorf("Third node Mz expected to be %f, but was %f", -5.0/3.0, mz)
	}
}

func TestDistributedGlobalLoadDistribution(t *testing.T) {
	element := structure.MakeElementBuilder(
		"1",
	).WithStartNode(
		structure.MakeFreeNodeAtPosition("1", 0.0, 0.0), &structure.DispConstraint,
	).WithEndNode(
		structure.MakeFreeNodeAtPosition("2", 4.0, 4.0), &structure.DispConstraint,
	).WithSection(
		structure.MakeUnitSection(),
	).WithMaterial(
		structure.MakeUnitMaterial(),
	).AddDistributedLoads(
		[]*load.DistributedLoad{
			load.MakeDistributed(load.FY, false, nums.MinT, 5.0, nums.MaxT, 5.0),
		},
	).Build()

	slicedEl := sliceLoadedElement(element, 2)

	// First Node
	if fx := slicedEl.Nodes[0].NetLocalFx(); !nums.FloatsEqual(fx, 5.0) {
		t.Errorf("First node Fx expected to be 5.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[0].NetLocalFy(); !nums.FloatsEqual(fy, 5.0) {
		t.Errorf("First node Fy expected to be 5.0, but was %f", fy)
	}
	if mz, expected := slicedEl.Nodes[0].NetLocalMz(), 10.0/math.Sqrt(18.0); !nums.FloatsEqual(mz, expected) {
		t.Errorf("First node Mz expected to be %f, but was %f", expected, mz)
	}

	// Second Node
	if fx := slicedEl.Nodes[1].NetLocalFx(); !nums.FloatsEqual(fx, 10.0) {
		t.Errorf("Second node Fx expected to be 10.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[1].NetLocalFy(); !nums.FloatsEqual(fy, 10.0) {
		t.Errorf("Second node Fy expected to be 10.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[1].NetLocalMz(); mz != 0.0 {
		t.Errorf("Second node Mz expected to be 0.0, but was %f", mz)
	}

	// Third Node
	if fx := slicedEl.Nodes[2].NetLocalFx(); !nums.FloatsEqual(fx, 5.0) {
		t.Errorf("Third node Fx expected to be 5.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[2].NetLocalFy(); !nums.FloatsEqual(fy, 5.0) {
		t.Errorf("Third node Fy expected to be 5.0, but was %f", fy)
	}
	if mz, expected := slicedEl.Nodes[2].NetLocalMz(), -10.0/math.Sqrt(18.0); !nums.FloatsEqual(mz, expected) {
		t.Errorf("Third node Mz expected to be %f, but was %f", expected, mz)
	}
}

func TestConcentratedLocalLoadDistribution(t *testing.T) {
	element := structure.MakeElementBuilder(
		"1",
	).WithStartNode(
		structure.MakeFreeNodeAtPosition("1", 0.0, 0.0), &structure.DispConstraint,
	).WithEndNode(
		structure.MakeFreeNodeAtPosition("2", 4.0, 0.0), &structure.DispConstraint,
	).WithSection(
		structure.MakeUnitSection(),
	).WithMaterial(
		structure.MakeUnitMaterial(),
	).AddConcentratedLoads(
		[]*load.ConcentratedLoad{
			load.MakeConcentrated(load.FX, true, nums.MakeTParam(0.25), 3.0),
			load.MakeConcentrated(load.FY, true, nums.MakeTParam(0.25), 5.0),
			load.MakeConcentrated(load.MZ, true, nums.MakeTParam(0.25), 7.0),
		},
	).Build()

	slicedEl := sliceLoadedElement(element, 2)

	if fx := slicedEl.Nodes[1].NetLocalFx(); fx != 3.0 {
		t.Errorf("First node Fx expected to be 3.0, but was %f", fx)
	}
	if fy := slicedEl.Nodes[1].NetLocalFy(); fy != 5.0 {
		t.Errorf("First node Fy expected to be 5.0, but was %f", fy)
	}
	if mz := slicedEl.Nodes[1].NetLocalMz(); mz != 7.0 {
		t.Errorf("First node Mz expected to be 7.0, but was %f", mz)
	}
}

func TestConcentratedGlobalLoadDistribution(t *testing.T) {
	element := structure.MakeElementBuilder(
		"1",
	).WithStartNode(
		structure.MakeFreeNodeAtPosition("1", 0.0, 0.0), &structure.DispConstraint,
	).WithEndNode(
		structure.MakeFreeNodeAtPosition("2", 4.0, 4.0), &structure.DispConstraint,
	).WithSection(
		structure.MakeUnitSection(),
	).WithMaterial(
		structure.MakeUnitMaterial(),
	).AddConcentratedLoads(
		[]*load.ConcentratedLoad{
			load.MakeConcentrated(load.FY, false, nums.MakeTParam(0.25), 5.0),
		},
	).Build()

	slicedEl := sliceLoadedElement(element, 2)

	if fx := slicedEl.Nodes[1].NetLocalFx(); !nums.FloatsEqual(fx, 5.0/math.Sqrt2) {
		t.Errorf("First node Fx expected to be %f, but was %f", 5.0/math.Sqrt2, fx)
	}
	if fy := slicedEl.Nodes[1].NetLocalFy(); !nums.FloatsEqual(fy, 5.0/math.Sqrt2) {
		t.Errorf("First node Fy expected to be %f, but was %f", 5.0/math.Sqrt2, fy)
	}
	if mz := slicedEl.Nodes[1].NetLocalMz(); mz != 0.0 {
		t.Errorf("First node Mz expected to be 0.0, but was %f", mz)
	}
}

/* <-- Utils --> */

func makeElementWithLoads(loads []*load.ConcentratedLoad) *structure.Element {
	return structure.MakeElementBuilder(
		"1",
	).WithStartNode(
		structure.MakeFreeNodeAtPosition("1", 1.0, 2.0), &structure.DispConstraint,
	).WithEndNode(
		structure.MakeFreeNodeAtPosition("2", 3.0, 4.0), &structure.DispConstraint,
	).WithMaterial(
		structure.MakeUnitMaterial(),
	).WithSection(
		structure.MakeUnitSection(),
	).AddConcentratedLoads(
		loads,
	).Build()
}

func makeElementWithoutLoads() *structure.Element {
	return structure.MakeElementBuilder(
		"1",
	).WithStartNode(
		structure.MakeFreeNodeAtPosition("1", 1.0, 2.0), &structure.DispConstraint,
	).WithEndNode(
		structure.MakeFreeNodeAtPosition("2", 3.0, 4.0), &structure.DispConstraint,
	).WithMaterial(
		structure.MakeUnitMaterial(),
	).WithSection(
		structure.MakeUnitSection(),
	).Build()
}
