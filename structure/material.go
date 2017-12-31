package structure

type Material struct {
	Name string
	Density, YoungMod, ShearMod, PoissonRatio, YieldStrength, UltimateStrength float64
}

/* Construction */
func MakeUnitMaterial() Material {
	return Material{"unit_material", 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
}
