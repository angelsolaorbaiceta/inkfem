package io

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	nodesHeaderRegex     = regexp.MustCompile(`(?:\|nodes\|\s*)(\d+)`)
	materialsHeaderRegex = regexp.MustCompile(`(?:\|materials\|\s*)(\d+)`)
	sectionsHeaderRegex  = regexp.MustCompile(`(?:\|sections\|\s*)(\d+)`)
	loadsHeaderRegex     = regexp.MustCompile(`(?:\|loads\|\s*)(\d+)`)
	barsHeaderRegex      = regexp.MustCompile(`(?:\|bars\|\s*)(\d+)`)
)

func IsNodesHeader(line string) bool {
	return nodesHeaderRegex.MatchString(line)
}

func ExtractNodesCount(nodesHeader string) int {
	count, err := strconv.Atoi(nodesHeaderRegex.FindStringSubmatch(nodesHeader)[1])
	if err != nil {
		panic(fmt.Sprintf("Can't read number of nodes from '%s'", nodesHeader))
	}

	return count
}

func IsMaterialsHeader(line string) bool {
	return materialsHeaderRegex.MatchString(line)
}

func ExtractMaterialsCount(materialsHeader string) int {
	count, err := strconv.Atoi(materialsHeaderRegex.FindStringSubmatch(materialsHeader)[1])
	if err != nil {
		panic(fmt.Sprintf("Can't read number of materials from '%s'", materialsHeader))
	}

	return count
}

func IsSectionsHeader(line string) bool {
	return sectionsHeaderRegex.MatchString(line)
}

func ExtractSectionsCount(sectionsHeader string) int {
	count, err := strconv.Atoi(sectionsHeaderRegex.FindStringSubmatch(sectionsHeader)[1])
	if err != nil {
		panic(fmt.Sprintf("Can't read number of sections from '%s'", sectionsHeader))
	}

	return count
}

func IsLoadsHeader(line string) bool {
	return loadsHeaderRegex.MatchString(line)
}

func ExtractLoadsCount(loadsHeader string) int {
	count, err := strconv.Atoi(loadsHeaderRegex.FindStringSubmatch(loadsHeader)[1])
	if err != nil {
		panic(fmt.Sprintf("Can't read number of loads from '%s'", loadsHeader))
	}

	return count
}

func IsBarsHeader(line string) bool {
	return barsHeaderRegex.MatchString(line)
}

func ExtractBarsCount(barsHeader string) int {
	count, err := strconv.Atoi(barsHeaderRegex.FindStringSubmatch(barsHeader)[1])
	if err != nil {
		panic(fmt.Sprintf("Can't read number of loads from '%s'", barsHeader))
	}

	return count
}
