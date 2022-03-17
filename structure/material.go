package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkmath/nums"
)

// Material represents a the substance from which resistant elements are made of.
// Materials have the properties of:
// - Density
// - Young Modulus
// - Shear Modulus
// - Poisson Ratio
// - Yield Strength
// - Ultimate Strength
type Material struct {
	Name                             string
	Density                          float64
	YoungMod, ShearMod, PoissonRatio float64
	YieldStrength, UltimateStrength  float64
}

// MakeMaterial creates a material with the given properties.
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

// MakeUnitMaterial creates a material with all properties set to 1.0.
func MakeUnitMaterial() *Material {
	return &Material{"unit_material", 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
}

// String representation of the material.
// This method is used for serialization, thus if the format is changed, the definition,
// preprocessed and solution file formats are affected.
func (m *Material) String() string {
	return fmt.Sprintf(
		"%s -> %f %f %f %f %f %f",
		m.Name,
		m.Density,
		m.YoungMod,
		m.ShearMod,
		m.PoissonRatio,
		m.YieldStrength,
		m.UltimateStrength,
	)
}

// Equals tests whether this and other materials are equal.
//
// Materials are equal if all its numerical properties are equal. the name isn't considered for the
// equality check.
func (m *Material) Equals(other *Material) bool {
	return nums.FuzzyEqual(m.Density, other.Density) &&
		nums.FuzzyEqual(m.YoungMod, other.YoungMod) &&
		nums.FuzzyEqual(m.ShearMod, other.ShearMod) &&
		nums.FuzzyEqual(m.PoissonRatio, other.PoissonRatio) &&
		nums.FuzzyEqual(m.YieldStrength, other.YieldStrength) &&
		nums.FuzzyEqual(m.UltimateStrength, other.UltimateStrength)
}
