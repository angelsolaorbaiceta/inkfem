package def

import (
	"fmt"
	"regexp"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

var (
	// <term> <reference-type> <elementId> <tStart> <valueStart> <tEnd> <valueEnd>
	distLoadDefinitionRegex = regexp.MustCompile(
		"^" + inkio.LoadTermExpr + inkio.DistributedLoadRefExpr +
			inkio.LoadElementID +
			inkio.FloatGroupExpr("t_start") + inkio.SpaceExpr +
			inkio.FloatGroupExpr("val_start") + inkio.SpaceExpr +
			inkio.FloatGroupExpr("t_end") + inkio.SpaceExpr +
			inkio.FloatGroupExpr("val_end") + inkio.OptionalSpaceExpr + "$",
	)

	// <term> <reference> <elementId> <t> <value>
	concLoadDefinitionRegex = regexp.MustCompile(
		"^" + inkio.LoadTermExpr + inkio.ConcentratedLoadRefExpr +
			inkio.LoadElementID +
			inkio.FloatGroupExpr("t") + inkio.SpaceExpr +
			inkio.FloatGroupExpr("val") + inkio.OptionalSpaceExpr + "$",
	)
)

func DeserializeLoad(line string) (contracts.StrID, *load.DistributedLoad, *load.ConcentratedLoad) {
	var (
		matchesDistributed  = distLoadDefinitionRegex.MatchString(line)
		matchesConcentrated = concLoadDefinitionRegex.MatchString(line)
	)

	if matchesDistributed {
		elementID, distributedLoad := deserializeDistributedLoad(line)
		return elementID, distributedLoad, nil
	}

	if matchesConcentrated {
		elementID, concentratedLoad := deserializeConcentratedLoad(line)
		return elementID, nil, concentratedLoad
	}

	panic(fmt.Sprintf("Found load with wrong format: '%s'", line))
}

func deserializeDistributedLoad(line string) (contracts.StrID, *load.DistributedLoad) {
	groups := distLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementID := groups[3]
	tStart := inkio.EnsureParseFloat(groups[4], "distributed load start T")
	valStart := inkio.EnsureParseFloat(groups[5], "distributed load start value")
	tEnd := inkio.EnsureParseFloat(groups[6], "distributed load end T")
	valEnd := inkio.EnsureParseFloat(groups[7], "distributed load end value")

	return elementID,
		load.MakeDistributed(
			term,
			isInLocalCoords,
			nums.MakeTParam(tStart),
			valStart,
			nums.MakeTParam(tEnd),
			valEnd,
		)
}

func deserializeConcentratedLoad(line string) (contracts.StrID, *load.ConcentratedLoad) {
	groups := concLoadDefinitionRegex.FindStringSubmatch(line)

	term := load.Term(groups[1])
	load.EnsureValidTerm(term)

	isInLocalCoords := groups[2] == "l"
	elementID := groups[3]
	t := inkio.EnsureParseFloat(groups[4], "concentrated load T")
	val := inkio.EnsureParseFloat(groups[5], "concentrated load value")

	return elementID, load.MakeConcentrated(term, isInLocalCoords, nums.MakeTParam(t), val)
}
