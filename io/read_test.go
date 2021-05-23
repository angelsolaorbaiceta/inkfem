package io

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestDeserializeNode(t *testing.T) {
	var (
		got  = deserializeNode("1 -> 10.1 20.2 { dx dy rz }")
		want = structure.MakeNode("1", g2d.MakePoint(10.1, 20.2), structure.FullConstraint)
	)

	if !got.Equals(want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestDeserializeNodes(t *testing.T) {
	var (
		lines []string = []string{
			"1 -> 10.1 20.2 { dx dy rz }",
			"2 -> 40.1 50.2 { dx dy }",
			"3 -> 70.1 80.2 { }",
		}
		nodes = deserializeNodesByID(lines)

		nodeOne   = structure.MakeNode("1", g2d.MakePoint(10.1, 20.2), structure.FullConstraint)
		nodeTwo   = structure.MakeNode("2", g2d.MakePoint(40.1, 50.2), structure.DispConstraint)
		nodeThree = structure.MakeNode("3", g2d.MakePoint(70.1, 80.2), structure.NilConstraint)
	)

	if size := len(*nodes); size != 3 {
		t.Errorf("Expected 3 nodes, but got %d", size)
	}
	if got := (*nodes)["1"]; !got.Equals(nodeOne) {
		t.Errorf("Expected node %v, but got %v", nodeOne, got)
	}
	if got := (*nodes)["2"]; !got.Equals(nodeTwo) {
		t.Errorf("Expected node %v, but got %v", nodeTwo, got)
	}
	if got := (*nodes)["3"]; !got.Equals(nodeThree) {
		t.Errorf("Expected node %v, but got %v", nodeThree, got)
	}
}

func TestDeserializeMaterial(t *testing.T) {
	var (
		got      = deserializeMaterial("'mat steel' -> 1.1 2.2 3.3 4.4 5.5 6.6")
		wantName = "mat steel"
		want     = structure.MakeMaterial(wantName, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6)
	)

	if got.Name != wantName {
		t.Errorf("Expected name %s, got '%s'", wantName, got.Name)
	}
	if !got.Equals(want) {
		t.Errorf("Wrong material. Want %v, got %v", want, got)
	}
}

func TestDeserializeMaterials(t *testing.T) {
	var (
		lines []string = []string{
			"'mat one' -> 1.1 2.2 3.3 4.4 5.5 6.6",
			"'mat two' -> 10.1 20.2 30.3 40.4 50.5 60.6",
		}
		materialsByName = deserializeMaterialsByName(lines)

		matOneName = "mat one"
		wantMatOne = structure.MakeMaterial(matOneName, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6)
		matTwoName = "mat two"
		wantMatTwo = structure.MakeMaterial(matTwoName, 10.1, 20.2, 30.3, 40.4, 50.5, 60.6)
	)

	if got := (*materialsByName)[matOneName]; !got.Equals(wantMatOne) {
		t.Errorf("Want material %v, got %v", wantMatOne, got)
	}
	if got := (*materialsByName)[matTwoName]; !got.Equals(wantMatTwo) {
		t.Errorf("Want material %v, got %v", wantMatTwo, got)
	}
}

func TestDeserializeSection(t *testing.T) {
	var (
		got      = deserializeSection("'IPE 100' -> 1.1 2.2 3.3 4.4 5.5")
		wantName = "IPE 100"
		want     = structure.MakeSection(wantName, 1.1, 2.2, 3.3, 4.4, 5.5)
	)

	if got.Name != wantName {
		t.Errorf("Expected name '%s', got '%s'", wantName, got.Name)
	}
	if !got.Equals(want) {
		t.Errorf("Expected section %v, got %v", want, got)
	}
}

func TestDeserializeSections(t *testing.T) {
	var (
		lines = []string{
			"'IPE 100' -> 1.1 2.2 3.3 4.4 5.5",
			"'IPE 200' -> 10.1 20.2 30.3 40.4 50.5",
		}
		sectionsByName = deserializeSectionsByName(lines)

		secOneName = "IPE 100"
		wantSecOne = structure.MakeSection(secOneName, 1.1, 2.2, 3.3, 4.4, 5.5)
		secTwoName = "IPE 200"
		wantSecTwo = structure.MakeSection(secTwoName, 10.1, 20.2, 30.3, 40.4, 50.5)
	)

	if got := (*sectionsByName)[secOneName]; !got.Equals(wantSecOne) {
		t.Errorf("Expected section %v, got %v", wantSecOne, got)
	}
	if got := (*sectionsByName)[secTwoName]; !got.Equals(wantSecTwo) {
		t.Errorf("Expected section %v, got %v", wantSecTwo, got)
	}
}

func TestDeserializeDistributedLoad(t *testing.T) {
	barID, gotLoad := deserializeDistributedLoad("fx ld 34 0.1 -50.2 0.9 -65.5")
	var (
		startT = inkgeom.MakeTParam(0.1)
		endT   = inkgeom.MakeTParam(0.9)
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
	want := load.MakeConcentrated(load.FY, false, inkgeom.HalfT, -70.5)

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

		startT  = inkgeom.MakeTParam(0.1)
		endT    = inkgeom.MakeTParam(0.9)
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

func TestDeserializeElements(t *testing.T) {
	var (
		lines = []string{
			"1 -> 1{ dx dy rz } 2{ dx dy } 'mat' 'sec'",
			"2 -> 1{ dx dy } 3{ dx dy rz } 'mat' 'sec'",
		}
		nodes = map[contracts.StrID]*structure.Node{
			"1": structure.MakeFreeNodeAtPosition("1", 100.0, 200.0),
			"2": structure.MakeFreeNodeAtPosition("2", 300.0, 400.0),
			"3": structure.MakeFreeNodeAtPosition("3", 500.0, 600.0),
		}
		materials = map[string]*structure.Material{
			"mat": structure.MakeMaterial("mat", 1.1, 2.2, 3.3, 4.4, 5.5, 6.6),
		}
		sections = map[string]*structure.Section{
			"sec": structure.MakeSection("sec", 10.1, 20.2, 30.3, 40.4, 50.5),
		}
		concentratedLoads = ConcLoadsById{
			"1": {load.MakeConcentrated(load.FY, true, inkgeom.MinT, -50)},
			"2": {load.MakeConcentrated(load.MZ, true, inkgeom.MaxT, -30)},
		}
		distributedLoads = DistLoadsById{}

		elements = deserializeElements(lines, &nodes, &materials, &sections, &concentratedLoads, &distributedLoads)

		wantElOne = structure.MakeElement(
			"1", nodes["1"], nodes["2"],
			structure.FullConstraint, structure.DispConstraint,
			materials["mat"], sections["sec"],
			concentratedLoads["1"],
			distributedLoads["1"],
		)
		wantElTwo = structure.MakeElement(
			"2", nodes["1"], nodes["3"],
			structure.DispConstraint, structure.FullConstraint,
			materials["mat"], sections["sec"],
			concentratedLoads["2"],
			distributedLoads["2"],
		)
	)

	if got := (*elements)[0]; !got.Equals(wantElOne) {
		t.Errorf("Expected element %v, got %v", wantElOne, got)
	}
	if got := (*elements)[1]; !got.Equals(wantElTwo) {
		t.Errorf("Expected element %v, got %v", wantElTwo, got)
	}
}
