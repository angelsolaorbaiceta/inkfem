package load

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom"
)

func TestLoadIsConcentrated(t *testing.T) {
	if load := MakeConcentrated(FX, true, inkgeom.MakeTParam(0.25), 45.0); !load.IsConcentrated() {
		t.Error("Expected 'concentrated' load")
	}
}

func TestLoadIsDistributed(t *testing.T) {
	if load := MakeDistributed(FX, true, inkgeom.MakeTParam(0.25), 45.0, inkgeom.MakeTParam(0.75), 67.0); !load.IsDistributed() {
		t.Error("Expected 'distributed' load")
	}
}

func TestLoadIsNodal(t *testing.T) {
	if load := MakeConcentrated(FX, true, inkgeom.MIN_T, 45.0); !load.IsNodal() {
		t.Error("Expected load to be nodal (t = 0.0)")
	}

	if load := MakeConcentrated(FX, true, inkgeom.MAX_T, 45.0); !load.IsNodal() {
		t.Error("Expected load to be nodal (t = 1.0)")
	}
}
