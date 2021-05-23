package io

import (
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	dispX = "dx"
	dispY = "dy"
	rotZ  = "rz"
)

func constraintFromString(str string) structure.Constraint {
	var (
		dxConst = strings.Contains(str, dispX)
		dyConst = strings.Contains(str, dispY)
		rzConst = strings.Contains(str, rotZ)
	)

	return structure.MakeConstraint(dxConst, dyConst, rzConst)
}
