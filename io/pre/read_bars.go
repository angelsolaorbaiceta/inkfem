package pre

import (
	"regexp"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	iodef "github.com/angelsolaorbaiceta/inkfem/io/def"
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

const (
	linesPerNode         = 6
	tPosGroupName        = "t"
	xPosGroupName        = "x"
	yPosGroupName        = "y"
	extTorsorGroupName   = "ext"
	leftTorsorGroupName  = "left"
	rightTorsorGroupName = "right"
	netTorsorGroupName   = "net"
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

	externalLoadPattern = regexp.MustCompile(
		"^" + "ext" + inkio.OptionalSpaceExpr + ":" + inkio.OptionalSpaceExpr +
			inkio.TorsorGroupExpr(extTorsorGroupName) + "$",
	)
	leftLoadPattern = regexp.MustCompile(
		"^" + "left" + inkio.OptionalSpaceExpr + ":" + inkio.OptionalSpaceExpr +
			inkio.TorsorGroupExpr(leftTorsorGroupName) + "$",
	)
	rightLoadPattern = regexp.MustCompile(
		"^" + "right" + inkio.OptionalSpaceExpr + ":" + inkio.OptionalSpaceExpr +
			inkio.TorsorGroupExpr(rightTorsorGroupName) + "$",
	)
	netLoadPattern = regexp.MustCompile(
		"^" + "net" + inkio.OptionalSpaceExpr + ":" + inkio.OptionalSpaceExpr +
			inkio.TorsorGroupExpr(netTorsorGroupName) + "$",
	)

	dofPattern = regexp.MustCompile(
		"^" + "dof" + inkio.OptionalSpaceExpr + ":" + inkio.OptionalSpaceExpr +
			inkio.DofGrpExpr + "$",
	)
)

func DeserializeBar(
	linesReader *inkio.LinesReader,
	data *structure.StructureData,
) *preprocess.Element {
	var (
		originalBarDTO, nNodes = iodef.DeserializeBar(linesReader.GetNextLine())
		originalBar            = iodef.BarFromDeserialization(originalBarDTO, data)
		nLines                 = nNodes * linesPerNode
		lines                  = linesReader.GetNextLines(nLines)
		nodes                  = deserializeBarNodes(nNodes, lines)
	)

	return preprocess.MakeElement(originalBar, nodes)
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
	var (
		t, pos           = parsePosition(lines[0])
		extLoad          = parseExternalLoad(lines[1])
		leftLoad         = parseLeftLoad(lines[2])
		rightLoad        = parseRightLoad(lines[3])
		netLoad          = parseNetLoad(lines[4])
		dof1, dof2, dof3 = parseDof(lines[5])
		node             = preprocess.MakeNode(t, pos, extLoad.Fx(), extLoad.Fy(), extLoad.Mz())
	)

	node.SetDegreesOfFreedomNum(dof1, dof2, dof3)
	node.AddLocalLeftLoad(leftLoad.Fx(), leftLoad.Fy(), leftLoad.Mz())
	node.AddLocalRightLoad(rightLoad.Fx(), rightLoad.Fy(), rightLoad.Mz())

	// Net load is added as a checksum. Ensure it checks out or panic.
	if !node.NetLocalLoadTorsor().Equals(netLoad) {
		panic("Expected net load doesn't match the read one")
	}

	return node
}

func parsePosition(line string) (nums.TParam, *g2d.Point) {
	var (
		groups = inkio.ExtractNamedGroups(positionPattern, line)
		t      = inkio.EnsureParseFloat(groups[tPosGroupName], "t position")
		x      = inkio.EnsureParseFloat(groups[xPosGroupName], "x position")
		y      = inkio.EnsureParseFloat(groups[yPosGroupName], "y position")
	)

	return nums.MakeTParam(t), g2d.MakePoint(x, y)
}

func parseExternalLoad(line string) *math.Torsor {
	groups := inkio.ExtractNamedGroups(externalLoadPattern, line)
	return inkio.EnsureParseTorsor(groups[extTorsorGroupName], "external load")
}

func parseLeftLoad(line string) *math.Torsor {
	groups := inkio.ExtractNamedGroups(leftLoadPattern, line)
	return inkio.EnsureParseTorsor(groups[leftTorsorGroupName], "left load")
}

func parseRightLoad(line string) *math.Torsor {
	groups := inkio.ExtractNamedGroups(rightLoadPattern, line)
	return inkio.EnsureParseTorsor(groups[rightTorsorGroupName], "right load")
}

func parseNetLoad(line string) *math.Torsor {
	groups := inkio.ExtractNamedGroups(netLoadPattern, line)
	return inkio.EnsureParseTorsor(groups[netTorsorGroupName], "net load")
}

func parseDof(line string) (int, int, int) {
	var (
		groups           = inkio.ExtractNamedGroups(dofPattern, line)
		dof1, dof2, dof3 = inkio.EnsureParseDOF(groups[inkio.DofGrpName], "node")
	)

	return dof1, dof2, dof3
}
