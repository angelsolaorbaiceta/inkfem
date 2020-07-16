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

func TestReadMaterial(t *testing.T) {
	print(materialDefinitionRegex.String())
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
