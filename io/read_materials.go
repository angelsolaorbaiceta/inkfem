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
	"bufio"
	"fmt"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// <name> -> <density> <young> <shear> <poisson> <yield> <ultimate>
var materialDefinitionRegex = regexp.MustCompile(
	`(?P<name>'\w+')(?:\s*->\s*)` +
		`(?P<density>\d+\.*\d*)(?:\s+)` +
		`(?P<young>\d+\.+\d+)(?:\s+)` +
		`(?P<shear>\d+\.+\d+)(?:\s+)` +
		`(?P<poisson>\d+\.+\d+)(?:\s+)` +
		`(?P<yield>\d+\.+\d+)(?:\s+)` +
		`(?P<ultimate>\d+\.+\d+)`)

func readMaterials(scanner *bufio.Scanner, count int) *map[string]*structure.Material {
	var (
		name                            string
		density, youngMod, shearMod     float64
		possonRatio                     float64
		yieldStrength, ultimateStrength float64
		materials                       = make(map[string]*structure.Material)
	)

	for _, line := range definitionLines(scanner, count) {
		if !materialDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found material with wrong format: '%s'", line))
		}

		groups := materialDefinitionRegex.FindStringSubmatch(line)

		name = groups[1]
		density, _ = strconv.ParseFloat(groups[2], 64)
		youngMod, _ = strconv.ParseFloat(groups[3], 64)
		shearMod, _ = strconv.ParseFloat(groups[4], 64)
		possonRatio, _ = strconv.ParseFloat(groups[5], 64)
		yieldStrength, _ = strconv.ParseFloat(groups[6], 64)
		ultimateStrength, _ = strconv.ParseFloat(groups[7], 64)

		materials[name] = &structure.Material{
			Name:             name,
			Density:          density,
			YoungMod:         youngMod,
			ShearMod:         shearMod,
			PoissonRatio:     possonRatio,
			YieldStrength:    yieldStrength,
			UltimateStrength: ultimateStrength}
	}

	return &materials
}
