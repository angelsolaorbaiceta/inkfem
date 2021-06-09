package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func applyDistributedLoadsToNodes(nodes []*Node, loads []*load.DistributedLoad) {
	var trailNode, leadNode *Node

	for i, j := 0, 1; j < len(nodes); i, j = i+1, j+1 {
		trailNode, leadNode = nodes[i], nodes[j]

		for _, load := range loads {
			applyDistributedLoadToNodes(load, trailNode, leadNode)
		}
	}
}

// TODO: distribute Mz loads
// Applies a distribute load to the trailing and leading nodes in a finite element.
func applyDistributedLoadToNodes(load *load.DistributedLoad, trailNode, leadNode *Node) {
	var (
		startLoad, endLoad = forceTorsorInLocalCoords(load, trailNode, leadNode)
		length             = trailNode.DistanceTo(leadNode)
		halfLength         = 0.5 * length
		length2            = length * length
		length3            = length2 * length
		loadSlopes         = computeLoadSlopes(startLoad, endLoad, length)
	)

	var (
		trailFx       = (startLoad.Fx() * halfLength) + (length2 * loadSlopes.Fx() / 6.0)
		trailFy       = (startLoad.Fy() * halfLength) + (3.0 * length2 * loadSlopes.Fy() / 20.0)
		trailFyMoment = (startLoad.Fy() * length2 / 12.0) + (length3 * loadSlopes.Fy() / 30.0)
	)
	trailNode.AddLocalLeftLoad(
		trailFx,
		trailFy,
		(startLoad.Mz()*halfLength)+trailFyMoment,
	)

	var (
		leadFx       = (startLoad.Fx() * halfLength) + (length2 * loadSlopes.Fx() / 3.0)
		leadFy       = (startLoad.Fy() * halfLength) + (7.0 * length2 * loadSlopes.Fy() / 20.0)
		leadFyMoment = -(startLoad.Fy() * length2 / 12.0) - (length3 * loadSlopes.Fy() / 20.0)
	)
	leadNode.AddLocalRightLoad(
		leadFx,
		leadFy,
		(startLoad.Mz()*halfLength)+leadFyMoment,
	)
}

func forceTorsorInLocalCoords(
	load *load.DistributedLoad,
	trailNode, leadNode *Node,
) (startLoad, endLoad *math.Torsor) {
	if load.IsInLocalCoords {
		startLoad = load.AsVectorAt(trailNode.T)
		endLoad = load.AsVectorAt(leadNode.T)
	} else {
		elementReferenceFrame := g2d.MakeRefFrameWithIVersor(
			g2d.MakeVectorFromTo(trailNode.Position, leadNode.Position),
		)

		startLoad = load.ProjectedVectorAt(trailNode.T, elementReferenceFrame)
		endLoad = load.ProjectedVectorAt(leadNode.T, elementReferenceFrame)
	}

	return
}

func computeLoadSlopes(startLoad, endLoad *math.Torsor, length float64) *math.Torsor {
	return math.MakeTorsor(
		(endLoad.Fx()-startLoad.Fx())/length,
		(endLoad.Fy()-startLoad.Fy())/length,
		(endLoad.Mz()-startLoad.Mz())/length,
	)
}
