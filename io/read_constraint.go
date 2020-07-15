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
