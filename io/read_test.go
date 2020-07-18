/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestReadNode(t *testing.T) {
	var (
		got  = deserializeNode("1 -> 10.1 20.2 { dx dy rz }")
		want = structure.MakeNode(1, g2d.MakePoint(10.1, 20.2), structure.FullConstraint)
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

		nodeOne   = structure.MakeNode(1, g2d.MakePoint(10.1, 20.2), structure.FullConstraint)
		nodeTwo   = structure.MakeNode(2, g2d.MakePoint(40.1, 50.2), structure.DispConstraint)
		nodeThree = structure.MakeNode(3, g2d.MakePoint(70.1, 80.2), structure.NilConstraint)
	)

	if size := len(*nodes); size != 3 {
		t.Errorf("Expected 3 nodes, but got %d", size)
	}
	if got := (*nodes)[1]; !got.Equals(nodeOne) {
		t.Errorf("Expected node %v, but got %v", nodeOne, got)
	}
	if got := (*nodes)[2]; !got.Equals(nodeTwo) {
		t.Errorf("Expected node %v, but got %v", nodeTwo, got)
	}
	if got := (*nodes)[3]; !got.Equals(nodeThree) {
		t.Errorf("Expected node %v, but got %v", nodeThree, got)
	}
}

func TestReadMaterial(t *testing.T) {
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

func TestReadSection(t *testing.T) {
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

func TestReadDistributedLoad(t *testing.T) {
	barID, gotLoad := deserializeDistributedLoad("fx ld 34 0.1 -50.2 0.9 -65.5")
	var (
		startT = inkgeom.MakeTParam(0.1)
		endT   = inkgeom.MakeTParam(0.9)
		want   = load.MakeDistributed(load.FX, true, startT, -50.2, endT, -65.5)
	)

	if barID != 34 {
		t.Errorf("Expected bar id 34, got %d", barID)
	}
	if !gotLoad.Equals(want) {
		t.Errorf("Expected load %v, got %v", want, gotLoad)
	}
}

func TestReadConcentratedLoad(t *testing.T) {
	barID, gotLoad := deserializeConcentratedLoad("fy gc 45 0.5 -70.5")
	want := load.MakeConcentrated(load.FY, false, inkgeom.HalfT, -70.5)

	if barID != 45 {
		t.Errorf("Expected bar id 45, got %d", barID)
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
		loads = deserializeLoadsByElementID(lines)[34]

		startT  = inkgeom.MakeTParam(0.1)
		endT    = inkgeom.MakeTParam(0.9)
		loadOne = load.MakeDistributed(load.FX, true, startT, -50.2, endT, -65.5)
		loadTwo = load.MakeConcentrated(load.FY, false, startT, -70.5)
	)

	if numberOfLoads := len(loads); numberOfLoads != 2 {
		t.Errorf("Expected 2 loads, got %d", numberOfLoads)
	}
	if got := loads[0]; !got.Equals(loadOne) {
		t.Errorf("Expected load %v, but got %v", loadOne, got)
	}
	if got := loads[1]; !got.Equals(loadTwo) {
		t.Errorf("Expected load %v, but got %v", loadTwo, got)
	}
}
