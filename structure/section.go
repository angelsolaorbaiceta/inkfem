package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkmath/nums"
)

// A Section of a resistant element.
type Section struct {
	Name           string
	Area           float64
	IStrong, IWeak float64 // Moments of Inertia
	SStrong, SWeak float64 // Section Moduli
}

// MakeUnitSection creates a section with all properties set to 1.0.
func MakeUnitSection() *Section {
	return &Section{"unit_section", 1.0, 1.0, 1.0, 1.0, 1.0}
}

// MakeSection creates a section with the given properties.
func MakeSection(name string, area, iStrong, iWeak, sStrong, sWeak float64) *Section {
	return &Section{
		Name:    name,
		Area:    area,
		IStrong: iStrong,
		IWeak:   iWeak,
		SStrong: sStrong,
		SWeak:   sWeak,
	}
}

// String representation of the section.
func (s *Section) String() string {
	return fmt.Sprintf(
		"'%s': %f %f %f %f %f",
		s.Name,
		s.Area,
		s.IStrong,
		s.IWeak,
		s.SStrong,
		s.SWeak,
	)
}

// Equals tests whether this and other sections are equal.
//
// Sections are equal if all its numerical properties are equal. the name isn't considered for the
// equality check.
func (s *Section) Equals(other *Section) bool {
	return nums.FloatsEqual(s.Area, other.Area) &&
		nums.FloatsEqual(s.IStrong, other.IStrong) &&
		nums.FloatsEqual(s.IWeak, other.IWeak) &&
		nums.FloatsEqual(s.SStrong, other.SStrong) &&
		nums.FloatsEqual(s.SWeak, other.SWeak)
}
