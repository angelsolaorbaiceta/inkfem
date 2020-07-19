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

package load

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/nums"
)

func TestLoadIsConcentrated(t *testing.T) {
	load := MakeConcentrated(FX, true, inkgeom.MinT, 45.0)

	if !load.IsConcentrated() {
		t.Error("Expected 'concentrated' load")
	}
}

func TestLoadIsDistributed(t *testing.T) {
	load := MakeDistributed(FX, true, inkgeom.MinT, 45.0, inkgeom.MaxT, 67.0)

	if !load.IsDistributed() {
		t.Error("Expected 'distributed' load")
	}
}

func TestLoadIsNodal(t *testing.T) {
	t.Run("distributed load isn't nodal", func(t *testing.T) {
		load := MakeDistributed(FX, false, inkgeom.MinT, 10, inkgeom.MaxT, 20)

		if load.IsNodal() {
			t.Error("Expected distributed load to not be nodal")
		}
	})

	t.Run("concentrated in the start position is nodal", func(t *testing.T) {
		load := MakeConcentrated(FX, true, inkgeom.MinT, 45.0)

		if !load.IsNodal() {
			t.Error("Expected load to be nodal (t = 0.0)")
		}
	})

	t.Run("concentrated in the end position is nodal", func(t *testing.T) {
		load := MakeConcentrated(FX, true, inkgeom.MaxT, 45.0)

		if !load.IsNodal() {
			t.Error("Expected load to be nodal (t = 1.0)")
		}
	})
}

/* <-- Avg Value --> */
func TestAvgValueAllCoveredByLoad(t *testing.T) {
	var (
		l     = MakeDistributed(FY, true, inkgeom.MinT, 50.0, inkgeom.MaxT, 50.0)
		value = l.AvgValueBetween(inkgeom.MakeTParam(0.2), inkgeom.MakeTParam(0.7))
	)

	if !nums.FuzzyEqual(value, 50.0) {
		t.Errorf("Average value not as expected: got %f, expected: %f", value, 50.0)
	}
}

func TestAvgValueNoneCoveredByLoad(t *testing.T) {
	var (
		startT = inkgeom.MakeTParam(0.2)
		endT   = inkgeom.MakeTParam(0.3)
		l      = MakeDistributed(FY, true, startT, 50.0, endT, 50.0)
		value  = l.AvgValueBetween(inkgeom.MakeTParam(0.4), inkgeom.MakeTParam(0.7))
	)

	if !nums.FuzzyEqual(value, 0.0) {
		t.Errorf("Average value not as expected: got %f, expected: %f", value, 0.0)
	}
}

func TestAvgValuePartiallyCoveredByLoad(t *testing.T) {
	var (
		startT = inkgeom.MakeTParam(0.2)
		endT   = inkgeom.MakeTParam(0.5)
		l      = MakeDistributed(FY, true, startT, 100.0, endT, 100.0)
	)

	value := l.AvgValueBetween(startT, inkgeom.MakeTParam(0.8))
	if !nums.FuzzyEqual(value, 50.0) {
		t.Errorf("Average value not as expected: got %f, expected: %f", value, 50.0)
	}

	value = l.AvgValueBetween(inkgeom.MakeTParam(0.1), inkgeom.MakeTParam(0.7))
	if !nums.FuzzyEqual(value, 50.0) {
		t.Errorf("Average value not as expected: got %f, expected: %f", value, 50.0)
	}
}
