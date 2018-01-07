package structure

// Section of a resistant element
type Section struct {
	Name           string
	Area           float64
	IStrong, IWeak float64 // Moments of Inertia
	SStrong, SWeak float64 // Section Moduli
}

/* Construction */

// MakeUnitSection returns a section with all properties set to 1.0.
func MakeUnitSection() Section {
	return Section{"unit_section", 1.0, 1.0, 1.0, 1.0, 1.0}
}
