package io

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// <name> -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
var sectionDefinitionRegex = regexp.MustCompile(`(?P<name>'\w+')(?:\s*->\s*)(?P<area>\d+\.*\d*)(?:\s+)(?P<istrong>\d+\.+\d+)(?:\s+)(?P<iweak>\d+\.+\d+)(?:\s+)(?P<sstrong>\d+\.+\d+)(?:\s+)(?P<sweak>\d+\.+\d+)`)

func readSections(scanner *bufio.Scanner, count int) map[string]structure.Section {
	var (
		name                                 string
		area, iStrong, iWeak, sStrong, sWeak float64
		sections                             = make(map[string]structure.Section)
	)

	for _, line := range definitionLines(scanner, count) {
		if !sectionDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found section with wrong format: '%s'", line))
		}

		groups := sectionDefinitionRegex.FindStringSubmatch(line)

		name = groups[1]
		area, _ = strconv.ParseFloat(groups[2], 64)
		iStrong, _ = strconv.ParseFloat(groups[3], 64)
		iWeak, _ = strconv.ParseFloat(groups[4], 64)
		sStrong, _ = strconv.ParseFloat(groups[5], 64)
		sWeak, _ = strconv.ParseFloat(groups[6], 64)

		sections[name] = structure.Section{
			Name:    name,
			Area:    area,
			IStrong: iStrong,
			IWeak:   iWeak,
			SStrong: sStrong,
			SWeak:   sWeak}
	}

	return sections
}
