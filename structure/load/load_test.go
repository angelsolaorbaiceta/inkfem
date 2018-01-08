package load

import (
	"fmt"
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath"
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

/* <---------- Avg Value ----------> */
func TestAvgValueAllCoveredByLoad(t *testing.T) {
	l := MakeDistributed(FY, true, inkgeom.MIN_T, 50.0, inkgeom.MAX_T, 50.0)

	if value := l.AvgValueBetween(inkgeom.MakeTParam(0.2), inkgeom.MakeTParam(0.7)); !inkmath.FuzzyEqual(value, 50.0) {
		t.Error(fmt.Sprintf("Average value not as expected: got %f, expected: %f", value, 50.0))
	}
}

func TestAvgValueNoneCoveredByLoad(t *testing.T) {
	l := MakeDistributed(FY, true, inkgeom.MakeTParam(0.2), 50.0, inkgeom.MakeTParam(0.3), 50.0)

	if value := l.AvgValueBetween(inkgeom.MakeTParam(0.4), inkgeom.MakeTParam(0.7)); !inkmath.FuzzyEqual(value, 0.0) {
		t.Error(fmt.Sprintf("Average value not as expected: got %f, expected: %f", value, 0.0))
	}
}

func TestAvgValuePartiallyCoveredByLoad(t *testing.T) {
	l := MakeDistributed(FY, true, inkgeom.MakeTParam(0.2), 100.0, inkgeom.MakeTParam(0.5), 100.0)

	if value := l.AvgValueBetween(inkgeom.MakeTParam(0.2), inkgeom.MakeTParam(0.8)); !inkmath.FuzzyEqual(value, 50.0) {
		t.Error(fmt.Sprintf("Average value not as expected: got %f, expected: %f", value, 50.0))
	}
	if value := l.AvgValueBetween(inkgeom.MakeTParam(0.1), inkgeom.MakeTParam(0.7)); !inkmath.FuzzyEqual(value, 50.0) {
		t.Error(fmt.Sprintf("Average value not as expected: got %f, expected: %f", value, 50.0))
	}
}
