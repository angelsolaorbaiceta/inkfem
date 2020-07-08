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
MakeUnitMaterial returns a material with all properties set to 1.0.
*/
func MakeUnitMaterial() *Material {
	return &Material{"unit_material", 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
}
