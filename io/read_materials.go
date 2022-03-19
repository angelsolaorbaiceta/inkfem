package io

import (
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// '<name>' -> <density> <young> <shear> <poisson> <yield> <ultimate>
var materialDefinitionRegex = regexp.MustCompile(
	"^" + NameGrpExpr + ArrowExpr +
		FloatGroupExpr("density") + SpaceExpr +
		FloatGroupExpr("young") + SpaceExpr +
		FloatGroupExpr("shear") + SpaceExpr +
		FloatGroupExpr("poisson") + SpaceExpr +
		FloatGroupExpr("yield") + SpaceExpr +
		FloatGroupExpr("ultimate") + OptionalSpaceExpr + "$")

func ReadMaterials(linesReader *LinesReader, count int) map[string]*structure.Material {
	lines := linesReader.GetNextLines(count)
	return deserializeMaterialsByName(lines)
}

func deserializeMaterialsByName(lines []string) map[string]*structure.Material {
	var (
		material  *structure.Material
		materials = make(map[string]*structure.Material)
	)

	for _, line := range lines {
		material = deserializeMaterial(line)
		materials[material.Name] = material
	}

	return materials
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
