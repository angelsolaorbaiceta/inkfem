package preprocess

import (
	"sort"
	"sync"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

/*
DoStructure preprocesses the structure by concurrently slicing each of the structural members.
*/
func DoStructure(s structure.Structure, wg *sync.WaitGroup) Structure {
	channel := make(chan Element, len(s.Elements))

	for _, element := range s.Elements {
		wg.Add(1)
		go DoElement(element, channel, wg)
	}
	wg.Wait()
	close(channel)

	var slicedElements []Element
	for slicedEl := range channel {
		slicedElements = append(slicedElements, slicedEl)
	}

	str := Structure{s.Metadata, s.Nodes, slicedElements}
	assignDof(&str)

	return str
}

/*
Assings degrees of freedom numbers to all nodes on sliced elements.

Structural nodes are given degrees of freedom to help in the correct assignment of DOF numbers
to the elements that meet in the node.
*/
func assignDof(s *Structure) {
	sort.Sort(ByGeometryPos(s.Elements))

	var (
		startNode/*, endNode*/ structure.Node
		startLink/*, endLink*/ *structure.Constraint
		dxDof, dyDof, rzDof int
		dof                 = 0
	)

	updateStructuralNodeDof := func(n *structure.Node) {
		if !n.HasDegreesOfFreedomNum() {
			n.SetDegreesOfFreedomNum(dof, dof+1, dof+2)
			s.Nodes[n.Id] = *n
			dof += 3
		}
	}

	for _, element := range s.Elements {
		startNode = s.Nodes[element.StartNodeID()]
		// endNode = s.Nodes[element.EndNodeID()]
		startLink = element.StartLink()
		// endLink = element.OriginalElement.EndLink

		// First Node
		updateStructuralNodeDof(&startNode)

		if startLink.AllowsDispX() {
			dxDof = dof
			dof++
		} else {
			dxDof = startNode.DegreesOfFreedomNum()[0]
		}

		if startLink.AllowsDispY() {
			dyDof = dof
			dof++
		} else {
			dyDof = startNode.DegreesOfFreedomNum()[1]
		}

		if startLink.AllowsRotation() {
			rzDof = dof
			dof++
		} else {
			rzDof = startNode.DegreesOfFreedomNum()[2]
		}

		element.Nodes[0].SetDegreesOfFreedomNum(dxDof, dyDof, rzDof)

		// Last Node
		// updateStructuralNodeDof(endNode)
	}
}
