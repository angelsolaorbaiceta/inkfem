package def

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestDeserializeDistributedLoad(t *testing.T) {
	barID, gotLoad := deserializeDistributedLoad("fx ld 34 0.1 -50.2 0.9 -65.5")
	var (
		startT = nums.MakeTParam(0.1)
		endT   = nums.MakeTParam(0.9)
		want   = load.MakeDistributed(load.FX, true, startT, -50.2, endT, -65.5)
	)

	if barID != "34" {
		t.Errorf("Expected bar id 34, got %s", barID)
	}
	if !gotLoad.Equals(want) {
		t.Errorf("Expected load %v, got %v", want, gotLoad)
	}
}

func TestDeserializeConcentratedLoad(t *testing.T) {
	barID, gotLoad := deserializeConcentratedLoad("fy gc 45 0.5 -70.5")
	want := load.MakeConcentrated(load.FY, false, nums.HalfT, -70.5)

	if barID != "45" {
		t.Errorf("Expected bar id 45, got %s", barID)
	}

	if !gotLoad.Equals(want) {
		t.Errorf("Expected load %v, got %v", want, gotLoad)
	}
}

func TestDeserializeLoads(t *testing.T) {
	var (
		lines []string = []string{
			"fx ld 34 0.1 -50.2 0.9 -65.5",
			"fy gc 34 0.1 -70.5",
		}
		allConcentrated, allDistributed = deserializeLoadsByElementID(lines)
		concentrated                    = allConcentrated["34"]
		distributed                     = allDistributed["34"]

		startT  = nums.MakeTParam(0.1)
		endT    = nums.MakeTParam(0.9)
		loadOne = load.MakeDistributed(load.FX, true, startT, -50.2, endT, -65.5)
		loadTwo = load.MakeConcentrated(load.FY, false, startT, -70.5)
	)

	if numberOfConcentratedLoads := len(concentrated); numberOfConcentratedLoads != 1 {
		t.Errorf("Expected 2 concentrated loads, got %d", numberOfConcentratedLoads)
	}
	if numberOfDistributedLoads := len(distributed); numberOfDistributedLoads != 1 {
		t.Errorf("Expected 2 distributed loads, got %d", numberOfDistributedLoads)
	}

	if got := distributed[0]; !got.Equals(loadOne) {
		t.Errorf("Expected distributed load %v, but got %v", loadOne, got)
	}
	if got := concentrated[0]; !got.Equals(loadTwo) {
		t.Errorf("Expected concentrated load %v, but got %v", loadTwo, got)
	}
}
