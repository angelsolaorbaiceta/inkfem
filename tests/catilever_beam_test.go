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

package tests

import (
	"testing"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom"
)

var (
	material = &structure.Material{
		Name: "steel", 
		Density: 0, 
		YoungMod: 20E6, 
		ShearMod: 0, 
		PoissonRatio: 1, 
		YieldStrength: 0, 
		UltimateStrength: 0,
	}
	section  = &structure.Section{
		Name: "IPE 120", 
		Area: 14, 
		IStrong: 318, 
		IWeak: 28, 
		SStrong: 53, 
		SWeak: 9,
	}
	displError = 1E-5
	solveOptions = process.SolveOptions{false, "", true, displError}
)

func TestCantileverBeamWithConcentratedLoadAtEnd(t *testing.T) {
	var (
		l = load.MakeConcentrated(load.FY, true, inkgeom.MaxT, -2000)
		str = makeBeamStructure([]load.Load{l})
		pre = preprocess.DoStructure(str)
		sol = process.Solve(pre, solveOptions)
		solutionElement = sol.Elements[0]
	)

	t.Run("global Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)
		maxDispl := -200.0 / 1908.0

		if got := solutionElement.GlobalYDispl[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.GlobalYDispl[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxDispl, displError) {
			t.Errorf("expected max displacement of %f, but got %f", maxDispl, got)
		}
	})
}

func makeBeamStructure(loads []load.Load) *structure.Structure {
	var(
		nodeOne = structure.MakeNode("fixed-node", g2d.MakePoint(0, 0), structure.FullConstraint)
		nodeTwo = structure.MakeNode("free-node", g2d.MakePoint(100, 0), structure.NilConstraint)
		beam = structure.MakeElement(
			"beam", 
			nodeOne, 
			nodeTwo, 
			structure.FullConstraint, 
			structure.FullConstraint, 
			material, 
			section, 
			loads,
		)
	)

	return &structure.Structure{
		structure.StrMetadata{1, 0},
		map[contracts.StrID]*structure.Node {
			nodeOne.Id: nodeOne,
			nodeTwo.Id: nodeTwo,
		},
		[]*structure.Element{beam},
	}
}