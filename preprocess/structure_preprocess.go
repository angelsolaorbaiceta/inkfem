package preprocess

import (
	"sort"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// StructureModel preprocesses the structure by concurrently slicing each of the structural members.
// The resulting sliced structure includes the degrees of freedom numbering.
func StructureModel(str *structure.Structure) *Structure {
	var (
		channel        = make(chan *Element)
		slicedElements []*Element
	)

	for _, element := range str.Elements {
		go sliceElement(element, channel)
	}

	for i := 0; i < str.ElementsCount(); i++ {
		slicedElements = append(slicedElements, <-channel)
	}

	slicedStr := &Structure{
		Metadata: str.Metadata,
		nodes:    str.NodesById(),
		Elements: slicedElements,
	}
	assignDof(slicedStr)

	return slicedStr
}

// Assings degrees of freedom numbers to all nodes on sliced elements.
//
// Structural nodes are given degrees of freedom to help in the correct assignment of DOF numbers
// to the elements that meet in the node. Structural elements are first sorted by their geometry
// positions, so the degrees of freedom numbers follow a logical sequence.
func assignDof(str *Structure) {
	sort.Sort(ByGeometryPos(str.Elements))

	var (
		startNode, endNode *structure.Node
		startLink, endLink *structure.Constraint
		nodesCount         int
		dof                = 0
	)

	updateStructuralNodeDof := func(node *structure.Node) {
		if !node.HasDegreesOfFreedomNum() {
			node.SetDegreesOfFreedomNum(dof, dof+1, dof+2)
			dof += 3
		}
	}

	endNodesDof := func(
		link *structure.Constraint,
		node *structure.Node,
	) (dxDof, dyDof, rzDof int) {
		if link.AllowsDispX() {
			dxDof = dof
			dof++
		} else {
			dxDof = node.DegreesOfFreedomNum()[0]
		}

		if link.AllowsDispY() {
			dyDof = dof
			dof++
		} else {
			dyDof = node.DegreesOfFreedomNum()[1]
		}

		if link.AllowsRotation() {
			rzDof = dof
			dof++
		} else {
			rzDof = node.DegreesOfFreedomNum()[2]
		}

		return
	}

	for _, element := range str.Elements {
		startNode, endNode = str.GetElementNodes(element)
		startLink = element.StartLink()
		endLink = element.EndLink()
		nodesCount = len(element.Nodes)

		/* First Node */
		updateStructuralNodeDof(startNode)
		element.Nodes[0].SetDegreesOfFreedomNum(
			endNodesDof(startLink, startNode),
		)

		/* Middle Nodes */
		for i := 1; i < nodesCount-1; i++ {
			element.Nodes[i].SetDegreesOfFreedomNum(dof, dof+1, dof+2)
			dof += 3
		}

		/* Last Node */
		updateStructuralNodeDof(endNode)
		element.Nodes[nodesCount-1].SetDegreesOfFreedomNum(
			endNodesDof(endLink, endNode),
		)
	}

	str.DofsCount = dof
}
