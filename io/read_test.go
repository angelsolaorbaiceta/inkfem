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
		got                = deserializeNode("1 -> 10.1 20.2 { dx dy rz }")
		expectedPosition   = g2d.MakePoint(10.1, 20.2)
		expectedConstraint = structure.FullConstraint
	)

	if got.Id != 1 {
		t.Errorf("Expected id 1, got %d", got.Id)
	}
	if !got.Position.Equals(expectedPosition) {
		t.Errorf("Expected position of (10.1, 20.2), got %v", got.Position)
	}
	if got.ExternalConstraint != expectedConstraint {
		t.Errorf("Expected constraint of { dx dy rz }, got %v", got.ExternalConstraint)
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
	got := deserializeMaterial("'mat steel' -> 1.1 2.2 3.3 4.4 5.5 6.6")

	if got.Name != "mat steel" {
		t.Errorf("Expected name 'mat steel', got '%s'", got.Name)
	}
	if got.Density != 1.1 {
		t.Errorf("Expected density of 1.1, got %f", got.Density)
	}
	if got.YoungMod != 2.2 {
		t.Errorf("Expected YoungMod of 2.2, got %f", got.YoungMod)
	}
	if got.ShearMod != 3.3 {
		t.Errorf("Expected ShearMod of 3.3, got %f", got.ShearMod)
	}
	if got.PoissonRatio != 4.4 {
		t.Errorf("Expected PoissonRatio of 4.4, got %f", got.PoissonRatio)
	}
	if got.YieldStrength != 5.5 {
		t.Errorf("Expected YieldStrength of 5.5, got %f", got.YieldStrength)
	}
	if got.UltimateStrength != 6.6 {
		t.Errorf("Expected UltimateStrength of 6.6, got %f", got.UltimateStrength)
	}
}

func TestReadSection(t *testing.T) {
	got := deserializeSection("'IPE 100' -> 1.1 2.2 3.3 4.4 5.5")

	if got.Name != "IPE 100" {
		t.Errorf("Expected name 'IPE 100', got '%s'", got.Name)
	}
	if got.Area != 1.1 {
		t.Errorf("Expected area of 1.1, got %f", got.Area)
	}
	if got.IStrong != 2.2 {
		t.Errorf("Expected IStrong of 2.2, got %f", got.IStrong)
	}
	if got.IWeak != 3.3 {
		t.Errorf("Expected IWeak of 3.3, got %f", got.IWeak)
	}
	if got.SStrong != 4.4 {
		t.Errorf("Expected SStrong of 4.4, got %f", got.SStrong)
	}
	if got.SWeak != 5.5 {
		t.Errorf("Expected SWeak of 5.5, got %f", got.SWeak)
	}
}

func TestReadDistributedLoad(t *testing.T) {
	barID, gotLoad := deserializeDistributedLoad("fx ld 34 0.1 -50.2 0.9 -65.5")
	var (
		startT = inkgeom.MakeTParam(0.1)
		endT   = inkgeom.MakeTParam(0.9)
	)

	if barID != 34 {
		t.Errorf("Expected bar id 34, got %d", barID)
	}

	if gotLoad.Term != load.FX {
		t.Errorf("Expected load term fx, got %s", gotLoad.Term)
	}
	if !gotLoad.IsInLocalCoords {
		t.Error("Expected load in local coords")
	}
	if gotLoad.StartT() != startT {
		t.Errorf("Expected load start t = 0.1, got %f", gotLoad.StartT())
	}
	if val := gotLoad.ValueAt(startT); val != -50.2 {
		t.Errorf("Expected load start value = -50.2, got %f", val)
	}
	if gotLoad.EndT() != endT {
		t.Errorf("Expected load end t = 0.9, got %f", gotLoad.EndT())
	}
	if val := gotLoad.ValueAt(endT); val != -65.5 {
		t.Errorf("Expected load end value = -65.5, got %f", val)
	}
}

func TestReadConcentratedLoad(t *testing.T) {
	barID, gotLoad := deserializeConcentratedLoad("fy gc 45 0.5 -70.5")

	if barID != 45 {
		t.Errorf("Expected bar id 45, got %d", barID)
	}

	if gotLoad.Term != load.FY {
		t.Errorf("Expected load term fy, got %s", gotLoad.Term)
	}
	if gotLoad.IsInLocalCoords {
		t.Error("Expected load in global coords")
	}
	if gotLoad.T() != inkgeom.HalfT {
		t.Errorf("Expected load t = 0.5, got %f", gotLoad.StartT())
	}
	if val := gotLoad.Value(); val != -70.5 {
		t.Errorf("Expected load start value = -70.5, got %f", val)
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
