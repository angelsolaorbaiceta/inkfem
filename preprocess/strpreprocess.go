package preprocess

import (
	"sort"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

/*
DoStructure preprocesses the structure by concurrently slicing each of the structural members.
*/
func DoStructure(s structure.Structure) Structure {
	channel := make(chan Element)

	for _, element := range s.Elements {
		go DoElement(element, channel)
	}

	var slicedElements []Element
	for i := 0; i < len(s.Elements); i++ {
		slicedElements = append(slicedElements, <-channel)
	}

	str := Structure{Metadata: s.Metadata, Nodes: s.Nodes, Elements: slicedElements}
	dofsCount := assignDof(&str)
	str.DofsCount = dofsCount

	return str
}

/*
Assings degrees of freedom numbers to all nodes on sliced elements.

Structural nodes are given degrees of freedom to help in the correct assignment of DOF numbers
to the elements that meet in the node. Structural elements are first sorted by their geometry
positions, so the degrees of freedom numbers follow a logical sequence.

The method returns the number degrees of freedom assigned.
*/
func assignDof(s *Structure) (dofsCount int) {
	sort.Sort(ByGeometryPos(s.Elements))

	var (
		startNode, endNode structure.Node
		startLink, endLink *structure.Constraint
		nodesCount         int
		dof                = 0
	)

	updateStructuralNodeDof := func(n *structure.Node) {
		if !n.HasDegreesOfFreedomNum() {
			n.SetDegreesOfFreedomNum(dof, dof+1, dof+2)
			s.Nodes[n.Id] = *n
			dof += 3
		}
	}

	endNodesDof := func(link *structure.Constraint, node structure.Node) (dxDof, dyDof, rzDof int) {
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

	for _, element := range s.Elements {
		startNode = s.Nodes[element.StartNodeID()]
		endNode = s.Nodes[element.EndNodeID()]
		startLink = element.StartLink()
		endLink = element.EndLink()
		nodesCount = len(element.Nodes)

		/* First Node */
		updateStructuralNodeDof(&startNode)
		element.Nodes[0].SetDegreesOfFreedomNum(
			endNodesDof(startLink, startNode),
		)

		/* Middle Nodes */
		for i := 1; i < nodesCount-1; i++ {
			element.Nodes[i].SetDegreesOfFreedomNum(dof, dof+1, dof+2)
			dof += 3
		}

		/* Last Node */
		updateStructuralNodeDof(&endNode)
		element.Nodes[nodesCount-1].SetDegreesOfFreedomNum(
			endNodesDof(endLink, endNode),
		)
	}

	dofsCount = dof
	return
}
