package structure

type Section struct {
	Name string
	Area float64
	IStrong, IWeak float64 // Moments of Inertia
	SStrong, SWeak float64 // Section Moduli
}

/* Construction */
func MakeUnitSection() Section {
	return Section{"unit_section", 1.0, 1.0, 1.0, 1.0, 1.0}
}
