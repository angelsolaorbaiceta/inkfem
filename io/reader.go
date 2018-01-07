package io

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
)

/* Constants & Regexps */
const (
	commentDeclaration = "#"
	dispX              = "dx"
	dispY              = "dy"
	rotZ               = "rz"
)

var (
	versionRegex = regexp.MustCompile(`(?:inkfem\s+v)(\d+)(?:[.])(\d+)`)

	// <id> -> <xCoord> <yCoord> {[dx dy dz]}
	nodesHeaderRegex    = regexp.MustCompile(`(?:\|nodes\|\s*)(\d+)`)
	nodeDefinitionRegex = regexp.MustCompile(`(?P<id>\d+)(?:\s*->\s*)(?P<x>\d+\.*\d*)(?:\s+)(?P<y>\d+\.*\d*)(?:\s+)(?P<constraints>{.*})`)

	// <name> -> <density> <young> <shear> <poisson> <yield> <ultimate>
	materialsHeaderRegex    = regexp.MustCompile(`(?:\|materials\|\s*)(\d+)`)
	materialDefinitionRegex = regexp.MustCompile(`(?P<name>'\w+')(?:\s*->\s*)(?P<density>\d+\.*\d*)(?:\s+)(?P<young>\d+\.+\d+)(?:\s+)(?P<shear>\d+\.+\d+)(?:\s+)(?P<poisson>\d+\.+\d+)(?:\s+)(?P<yield>\d+\.+\d+)(?:\s+)(?P<ultimate>\d+\.+\d+)`)

	// <name> -> <area> <iStrong> <iWeak> <sStrong> <sWeak>
	sectionsHeaderRegex    = regexp.MustCompile(`(?:\|sections\|\s*)(\d+)`)
	sectionDefinitionRegex = regexp.MustCompile(`(?P<name>'\w+')(?:\s*->\s*)(?P<area>\d+\.*\d*)(?:\s+)(?P<istrong>\d+\.+\d+)(?:\s+)(?P<iweak>\d+\.+\d+)(?:\s+)(?P<sstrong>\d+\.+\d+)(?:\s+)(?P<sweak>\d+\.+\d+)`)

	loadsHeaderRegex        = regexp.MustCompile(`(?:\|loads\|\s*)(\d+)`)
	distLoadDefinitionRegex = regexp.MustCompile(`(?P<term>[fm]{1}[xyz]{1})(?:\s+)(?P<ref>[lg]{1})(?:d{1})(?:\s+)(?P<element>\d+)(?:\s+)(?P<t_start>\d+\.*\d*)(?:\s+)(?P<val_start>-*\d+\.*\d*)(?:\s+)(?P<t_end>\d+\.*\d*)(?:\s+)(?P<val_end>-*\d+\.*\d*)`)
	concLoadDefinitionRegex = regexp.MustCompile(`(?P<term>[fm]{1}[xyz]{1})(?:\s+)(?P<ref>[lg]{1})(?:c{1})(?:\s+)(?P<element>\d+)(?:\s+)(?P<t>\d+\.*\d*)(?:\s+)(?P<val>-*\d+\.*\d*)`)

	// <id> -> <s_node> {[dx dy rz]} <e_node> {[dx dy rz]} <material> <section>
	elementsHeaderRegex    = regexp.MustCompile(`(?:\|elements\|\s*)(\d+)`)
	elementDefinitionRegex = regexp.MustCompile(`(?P<id>\d+)(?:\s*->\s*)(?P<start_node>\d+)(?:\s*)(?P<start_link>{.*})(?:\s+)(?P<end_node>\d+)(?:\s*)(?P<end_link>{.*})(?:\s+)(?P<material>'[A-Za-z0-9_ ]+')(?:\s+)(?P<section>'[A-Za-z0-9_ ]+')`)
)

/*
StructureFromFile Reads the given .inkfem file and tries to parse a structure
from the data defined.

The first line in the file should be as follows: 'inkfem vM.m', where 'M' and 'm'
are the major and minor version numbers of inkfem used to produce the file or
required to compute the structure.
*/
func StructureFromFile(filePath string) structure.Structure {
	file, error := os.Open(filePath)
	if error != nil {
		log.Fatal(error)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	return parseStructure(scanner)
}

func parseStructure(scanner *bufio.Scanner) structure.Structure {
	var (
		nodesDefined, materialsDefined, sectionsDefined, loadsDefined bool = false, false, false, false
		majorVersion, minorVersion                                    int
		nodes                                                         map[int]structure.Node
		materials                                                     map[string]structure.Material
		sections                                                      map[string]structure.Section
		loads                                                         map[int][]load.Load
		elements                                                      []structure.Element
	)

	// First line must be "inkfem vM.m"
	scanner.Scan()
	majorVersion, minorVersion = parseVersionNumbers(scanner.Text())

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if lineIsComment(line) || lineIsEmpty(line) {
			continue
		}

		switch {
		case nodesHeaderRegex.MatchString(line):
			{
				nodesCount, _ := strconv.Atoi(nodesHeaderRegex.FindStringSubmatch(line)[1])
				nodes = readNodes(scanner, nodesCount)
				nodesDefined = true
			}
		case materialsHeaderRegex.MatchString(line):
			{
				materialsCount, _ := strconv.Atoi(materialsHeaderRegex.FindStringSubmatch(line)[1])
				materials = readMaterials(scanner, materialsCount)
				materialsDefined = true
			}
		case sectionsHeaderRegex.MatchString(line):
			{
				sectionsCount, _ := strconv.Atoi(sectionsHeaderRegex.FindStringSubmatch(line)[1])
				sections = readSections(scanner, sectionsCount)
				sectionsDefined = true
			}
		case loadsHeaderRegex.MatchString(line):
			{
				loadsCount, _ := strconv.Atoi(loadsHeaderRegex.FindStringSubmatch(line)[1])
				loads = readLoads(scanner, loadsCount)
				loadsDefined = true
			}
		case elementsHeaderRegex.MatchString(line):
			{
				if !(nodesDefined && materialsDefined && sectionsDefined && loadsDefined) {
					panic("Cannot define elements if some of the following not already defined: nodes, materials, sections and loads")
				}

				elementsCount, _ := strconv.Atoi(elementsHeaderRegex.FindStringSubmatch(line)[1])
				elements = readElements(scanner, elementsCount, nodes, materials, sections, loads)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return structure.Structure{
		Metadata: structure.StrMetadata{
			MajorVersion: majorVersion,
			MinorVersion: minorVersion},
		Nodes:    nodes,
		Elements: elements}
}

/* <---------- READ : Version Numbers ----------> */
func parseVersionNumbers(firstLine string) (majorVersion, minorVersion int) {
	if foundMatch := versionRegex.MatchString(firstLine); !foundMatch {
		panic("Could not parse major and minor version numbers. Missing 'inkfem vM.m' in your file's first line?")
	}

	versions := versionRegex.FindStringSubmatch(firstLine)
	majorVersion, _ = strconv.Atoi(versions[1])
	minorVersion, _ = strconv.Atoi(versions[2])

	return
}

/* <---------- READ : Nodes ----------> */
func readNodes(scanner *bufio.Scanner, count int) map[int]structure.Node {
	var (
		id                 int
		x, y               float64
		externalConstraint string
		nodes              = make(map[int]structure.Node)
	)

	for _, line := range definitionLines(scanner, count) {
		if !nodeDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found node with wrong format: '%s'", line))
		}

		groups := nodeDefinitionRegex.FindStringSubmatch(line)

		id, _ = strconv.Atoi(groups[1])
		x, _ = strconv.ParseFloat(groups[2], 64)
		y, _ = strconv.ParseFloat(groups[3], 64)
		externalConstraint = groups[4]

		nodes[id] = structure.MakeNode(
			id,
			inkgeom.MakePoint(x, y),
			constraintFromString(externalConstraint))
	}

	return nodes
}

/* <---------- READ : Materials ----------> */
func readMaterials(scanner *bufio.Scanner, count int) map[string]structure.Material {
	var (
		name                                                                      string
		density, youngMod, shearMod, possonRatio, yieldStrength, ultimateStrength float64
		materials                                                                 = make(map[string]structure.Material)
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

		materials[name] = structure.Material{
			Name:             name,
			Density:          density,
			YoungMod:         youngMod,
			ShearMod:         shearMod,
			PoissonRatio:     possonRatio,
			YieldStrength:    yieldStrength,
			UltimateStrength: ultimateStrength}
	}

	return materials
}

/* <---------- READ : Sections ----------> */
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

/* <---------- READ : Loads ----------> */
func readLoads(scanner *bufio.Scanner, count int) map[int][]load.Load {
	var (
		elementNumber int
		_load         load.Load
		loads         = make(map[int][]load.Load)
	)

	for _, line := range definitionLines(scanner, count) {
		if !(distLoadDefinitionRegex.MatchString(line) || concLoadDefinitionRegex.MatchString(line)) {
			panic(fmt.Sprintf("Found load with wrong format: '%s'", line))
		}

		switch {
		case distLoadDefinitionRegex.MatchString(line):
			elementNumber, _load = distributedLoadFromString(line)

		case concLoadDefinitionRegex.MatchString(line):
			elementNumber, _load = concentratedLoadFromString(line)

		default:
			// shouldn't happen
			panic("Unknown type of load?")
		}

		loads[elementNumber] = append(loads[elementNumber], _load)
	}

	return loads
}

func distributedLoadFromString(line string) (int, load.Load) {
	groups := distLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.LoadTerm(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementNumber, _ := strconv.Atoi(groups[3])
	tStart, _ := strconv.ParseFloat(groups[4], 64)
	valStart, _ := strconv.ParseFloat(groups[5], 64)
	tEnd, _ := strconv.ParseFloat(groups[6], 64)
	valEnd, _ := strconv.ParseFloat(groups[7], 64)

	return elementNumber, load.MakeDistributed(term, isInLocalCoords, inkgeom.MakeTParam(tStart), valStart, inkgeom.MakeTParam(tEnd), valEnd)
}

func concentratedLoadFromString(line string) (int, load.Load) {
	groups := concLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.LoadTerm(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementNumber, _ := strconv.Atoi(groups[3])
	t, _ := strconv.ParseFloat(groups[4], 64)
	val, _ := strconv.ParseFloat(groups[5], 64)

	return elementNumber, load.MakeConcentrated(term, isInLocalCoords, inkgeom.MakeTParam(t), val)
}

/* <---------- READ : Elements ----------> */
func readElements(scanner *bufio.Scanner, count int, nodes map[int]structure.Node, materials map[string]structure.Material, sections map[string]structure.Section, loads map[int][]load.Load) []structure.Element {
	var (
		id, startNodeID, endNodeID int
		startNode, endNode         structure.Node
		startLink, endLink         string
		material                   structure.Material
		section                    structure.Section
		ok                         bool
		elements                   = make([]structure.Element, count)
	)

	for i, line := range definitionLines(scanner, count) {
		if !elementDefinitionRegex.MatchString(line) {
			panic(fmt.Sprintf("Found element with wrong format: '%s'", line))
		}

		groups := elementDefinitionRegex.FindStringSubmatch(line)

		id, _ = strconv.Atoi(groups[1])

		startNodeID, _ = strconv.Atoi(groups[2])
		startNode, ok = nodes[startNodeID]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown start node id: %d", id, startNodeID))
		}

		startLink = groups[3]

		endNodeID, _ = strconv.Atoi(groups[4])
		endNode, ok = nodes[endNodeID]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown end node id: %d", id, endNodeID))
		}

		endLink = groups[5]

		material, ok = materials[groups[6]]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown material name: %s", id, groups[6]))
		}

		section, ok = sections[groups[7]]
		if !ok {
			panic(fmt.Sprintf("Element %d with unknown section name: %s", id, groups[7]))
		}

		elements[i] = structure.MakeElement(
			id, startNode, endNode,
			constraintFromString(startLink),
			constraintFromString(endLink),
			material,
			section,
			loads[id])
	}

	return elements
}

/* Utils */
func lineIsComment(line string) bool {
	return strings.HasPrefix(line, commentDeclaration)
}

func lineIsEmpty(line string) bool {
	return len(line) < 1
}

func definitionLines(scanner *bufio.Scanner, count int) []string {
	var (
		line  string
		lines = make([]string, count)
	)

	for i := 0; i < count; {
		if !scanner.Scan() {
			panic("Couldn't read all expected lines")
		}

		line = scanner.Text()
		if lineIsComment(line) {
			continue
		}

		lines[i] = line
		i++
	}

	return lines
}

func constraintFromString(str string) structure.Constraint {
	dxConst := strings.Contains(str, dispX)
	dyConst := strings.Contains(str, dispY)
	rzConst := strings.Contains(str, rotZ)

	return structure.MakeConstraint(dxConst, dyConst, rzConst)
}
