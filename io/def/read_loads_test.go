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
