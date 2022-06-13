package def

import (
	"fmt"
	"regexp"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// '<name>' -> <density> <young> <shear> <poisson> <yield> <ultimate>
var materialDefinitionRegex = regexp.MustCompile(
	"^" + inkio.NameGrpExpr + inkio.ArrowExpr +
		inkio.FloatGroupExpr("density") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("young") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("shear") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("poisson") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("yield") + inkio.SpaceExpr +
		inkio.FloatGroupExpr("ultimate") + inkio.OptionalSpaceExpr + "$")

func DeserializeMaterial(definition string) *structure.Material {
	if !materialDefinitionRegex.MatchString(definition) {
		panic(fmt.Sprintf("Found material with wrong format: '%s'", definition))
	}

	groups := materialDefinitionRegex.FindStringSubmatch(definition)

	name := groups[1]
	density := inkio.EnsureParseFloat(groups[2], "material density")
	youngMod := inkio.EnsureParseFloat(groups[3], "material Young modulus")
	shearMod := inkio.EnsureParseFloat(groups[4], "material shear modulus")
	possonRatio := inkio.EnsureParseFloat(groups[5], "material poisson ratio")
	yieldStrength := inkio.EnsureParseFloat(groups[6], "material yield strength")
	ultimateStrength := inkio.EnsureParseFloat(groups[7], "material ultimate strength")

	return &structure.Material{
		Name:             name,
		Density:          density,
		YoungMod:         youngMod,
		ShearMod:         shearMod,
		PoissonRatio:     possonRatio,
		YieldStrength:    yieldStrength,
		UltimateStrength: ultimateStrength}
}
