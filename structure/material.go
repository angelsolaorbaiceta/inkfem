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

package structure

import "github.com/angelsolaorbaiceta/inkmath/nums"

/*
Material represents a the substance from which resistant elements are made of.
Materials have the properties of:
	- Density
	- Young Modulus
	- Shear Modulus
	- Poisson Ratio
	- Yield Strength
	- Ultimate Strength
*/
type Material struct {
	Name                             string
	Density                          float64
	YoungMod, ShearMod, PoissonRatio float64
	YieldStrength, UltimateStrength  float64
}

/* <-- Construction --> */

/*
MakeMaterial creates a material with the given properties.
*/
func MakeMaterial(name string, density, young, shear, poisson, yield, ultimate float64) *Material {
	return &Material{
		Name:             name,
		Density:          density,
		YoungMod:         young,
		ShearMod:         shear,
		PoissonRatio:     poisson,
		YieldStrength:    yield,
		UltimateStrength: ultimate,
	}
}

/*
MakeUnitMaterial creates a material with all properties set to 1.0.
*/
func MakeUnitMaterial() *Material {
	return &Material{"unit_material", 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
}

/* <-- Methods --> */

/*
Equals tests whether this and other material are equal.

Materials are equal if all its numerical properties are equal. the name isn't considered
for the equality check.
*/
func (m *Material) Equals(other *Material) bool {
	return nums.FuzzyEqual(m.Density, other.Density) &&
		nums.FuzzyEqual(m.YoungMod, other.YoungMod) &&
		nums.FuzzyEqual(m.ShearMod, other.ShearMod) &&
		nums.FuzzyEqual(m.PoissonRatio, other.PoissonRatio) &&
		nums.FuzzyEqual(m.YieldStrength, other.YieldStrength) &&
		nums.FuzzyEqual(m.UltimateStrength, other.UltimateStrength)
}
