package preprocess

import (
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

/*
Applies a distribute load to the trailing and leading nodes in a finite element.

TODO: distribute Fx loads
TODO: distribute Mz loads
*/
func applyDistributedLoadToNodes(load *load.DistributedLoad, trailNode, leadNode *Node) {
	var (
		startLoad, endLoad = forceVectorInLocalCoords(load, trailNode, leadNode)
		length             = trailNode.DistanceTo(leadNode)
		halfLength         = 0.5 * length
		length2            = length * length
		length3            = length2 * length
		loadValueSlopes    = computeLoadValueSlopes(startLoad, endLoad, length)
	)

	var (
		trailFy       = (startLoad[1] * halfLength) + (3.0 * length2 * loadValueSlopes[1] / 20.0)
		trailFyMoment = (startLoad[1] * length2 / 12.0) + (length3 * loadValueSlopes[1] / 30.0)
	)
	trailNode.AddLocalLeftLoad(
		startLoad[0]*halfLength,
		trailFy,
		(startLoad[2]*halfLength)+trailFyMoment,
	)

	var (
		leadFy       = (startLoad[1] * halfLength) + (7.0 * length2 * loadValueSlopes[1] / 20.0)
		leadFyMoment = -(startLoad[1] * length2 / 12.0) - (length3 * loadValueSlopes[1] / 20.0)
	)
	leadNode.AddLocalRightLoad(
		startLoad[0]*halfLength,
		leadFy,
		(startLoad[2]*halfLength)+leadFyMoment,
	)
}

func forceVectorInLocalCoords(load *load.DistributedLoad, trailNode, leadNode *Node) (startLoad, endLoad [3]float64) {
	if load.IsInLocalCoords {
		startLoad = load.AsVectorAt(trailNode.T)
		endLoad = load.AsVectorAt(leadNode.T)
	} else {
		elementReferenceFrame := g2d.MakeRefFrameWithIVersor(g2d.MakeVectorFromTo(trailNode.Position, leadNode.Position))
		startLoad = load.ProjectedVectorAt(trailNode.T, elementReferenceFrame)
		endLoad = load.ProjectedVectorAt(leadNode.T, elementReferenceFrame)
	}

	return
}

func computeLoadValueSlopes(startLoad, endLoad [3]float64, length float64) [3]float64 {
	return [3]float64{
		(endLoad[0] - startLoad[0]) / length,
		(endLoad[1] - startLoad[1]) / length,
		(endLoad[2] - startLoad[2]) / length,
	}
}
