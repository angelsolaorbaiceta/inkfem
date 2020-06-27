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
