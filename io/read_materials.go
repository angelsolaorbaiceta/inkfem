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
	"bufio"
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// '<name>' -> <density> <young> <shear> <poisson> <yield> <ultimate>
var materialDefinitionRegex = regexp.MustCompile(
	"^" + nameGrpExpr + arrowExpr +
		floatGroupAndSpaceExpr("density") +
		floatGroupAndSpaceExpr("young") +
		floatGroupAndSpaceExpr("shear") +
		floatGroupAndSpaceExpr("poisson") +
		floatGroupAndSpaceExpr("yield") +
		floatGroupAndOptinalSpaceExpr("ultimate") + "$")

func readMaterials(scanner *bufio.Scanner, count int) *map[string]*structure.Material {
	var (
		material  *structure.Material
		materials = make(map[string]*structure.Material)
	)

	for _, line := range definitionLines(scanner, count) {
		material = deserializeMaterial(line)
		materials[material.Name] = material
	}

	return &materials
}

func deserializeMaterial(definition string) *structure.Material {
	if !materialDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found material with wrong format: '%s'", definition))
	}

	groups := materialDefinitionRegex.FindStringSubmatch(definition)

	name := groups[1]
	density := ensureParseFloat(groups[2], "material density")
	youngMod := ensureParseFloat(groups[3], "material Young modulus")
	shearMod := ensureParseFloat(groups[4], "material shear modulus")
	possonRatio := ensureParseFloat(groups[5], "material poisson ratio")
	yieldStrength := ensureParseFloat(groups[6], "material yield strength")
	ultimateStrength := ensureParseFloat(groups[7], "material ultimate strength")

	return &structure.Material{
		Name:             name,
		Density:          density,
		YoungMod:         youngMod,
		ShearMod:         shearMod,
		PoissonRatio:     possonRatio,
		YieldStrength:    yieldStrength,
		UltimateStrength: ultimateStrength}
}
