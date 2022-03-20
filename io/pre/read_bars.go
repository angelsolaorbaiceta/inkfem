package pre

import (
	"regexp"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

const (
	linesPerNode  = 5
	tPosGroupName = "t"
	xPosGroupName = "x"
	yPosGroupName = "y"
)

var (
	positionPattern = regexp.MustCompile(
		"^" +
			inkio.FloatGroupExpr(tPosGroupName) +
			inkio.OptionalSpaceExpr + ":" + inkio.OptionalSpaceExpr +
			inkio.FloatGroupExpr(xPosGroupName) +
			inkio.SpaceExpr +
			inkio.FloatGroupExpr(yPosGroupName) +
			"$",
	)
)

func readBars(
	linesReader *inkio.LinesReader,
	count int,
	data *structure.StructureData,
) []*preprocess.Element {
	bars := make([]*preprocess.Element, count)

	for i := 0; i < count; i++ {
		bars[i] = deserializeBar(linesReader, data)
	}

	return bars
}

func deserializeBar(
	linesReader *inkio.LinesReader,
	data *structure.StructureData,
) *preprocess.Element {
	linesReader.ReadNext()

	var (
		readOps                 = inkio.ReaderOptions{ShouldIncludeOwnWeight: false}
		originalElement, nNodes = inkio.DeserializeBar(linesReader.GetNextLine(), data, readOps)
		nLines                  = nNodes * linesPerNode
		lines                   = linesReader.GetNextLines(nLines)
		nodes                   = deserializeBarNodes(nNodes, lines)
	)

	return preprocess.MakeElement(originalElement, nodes)
}

func deserializeBarNodes(nNodes int, lines []string) []*preprocess.Node {
	var (
		nodes            = make([]*preprocess.Node, nNodes)
		idxStart, idxEnd int
	)

	for i := 0; i < nNodes; i++ {
		idxStart = i * linesPerNode
		idxEnd = idxStart + linesPerNode
		nodes[i] = deserializeNode(lines[idxStart:idxEnd])
	}

	return nodes
}

func deserializeNode(lines []string) *preprocess.Node {
	// 0.000000 : 0.000000 0.000000
	//  ext   : {}
	//  left  : {0.000000 -970.000000 -1600.000000}
	//  right : {0.000000 0.000000 0.000000}
	//  net   : {0.000000 -970.000000 -1600.000000}
	//  dof   : [0 1 2]
	var (
		t, pos = parsePosition(lines[0])
	)

	node := preprocess.MakeNode(t, pos, 0, 0, 0)

	return node
}

func parsePosition(line string) (nums.TParam, *g2d.Point) {
	// var (
	// 	groups = inkio.ExtractNamedGroups(positionPattern, line)
	// )

	return nums.MinT, g2d.MakePoint(1, 2)
}
