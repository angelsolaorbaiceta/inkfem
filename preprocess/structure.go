package preprocess

import (
	"sort"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// Structure result of preprocessing original structure, ready to be solved.
// The elements of a preprocessed structure are already sliced.
type Structure struct {
	Metadata structure.StrMetadata
	structure.NodesById
	Elements  []*Element
	dofsCount int
}

// GetElementNodes returns the element's start and end nodes.
func (s *Structure) GetElementNodes(element *Element) (*structure.Node, *structure.Node) {
	return s.GetNodeById(element.StartNodeID()), s.GetNodeById(element.EndNodeID())
}

// ElementsCount returns the number of elements in the original structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}

// DofsCount is the number of degrees of freedom in the preprocessed structure.
func (s *Structure) DofsCount() int {
	return s.dofsCount
}

// Assings degrees of freedom numbers to all nodes on sliced elements.
//
// Structural nodes are given degrees of freedom to help in the correct assignment of DOF numbers
// to the elements that meet in the node. Structural elements are first sorted by their geometry
// positions, so the degrees of freedom numbers follow a logical sequence.
func (str *Structure) assignDof() {
	sort.Sort(ByGeometryPos(str.Elements))

	var (
		startNode, endNode *structure.Node
		startLink, endLink *structure.Constraint
		nodesCount         int
		dof                = 0
	)

	assignNodeDof := func(node *structure.Node) {
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
		assignNodeDof(startNode)
		element.Nodes[0].SetDegreesOfFreedomNum(
			endNodesDof(startLink, startNode),
		)

		/* Middle Nodes */
		for i := 1; i < nodesCount-1; i++ {
			element.Nodes[i].SetDegreesOfFreedomNum(dof, dof+1, dof+2)
			dof += 3
		}

		/* Last Node */
		assignNodeDof(endNode)
		element.Nodes[nodesCount-1].SetDegreesOfFreedomNum(
			endNodesDof(endLink, endNode),
		)
	}

	str.dofsCount = dof
}

// AddDispConstraints sets the node's external constraints in the system of equations
// matrix and vector.
//
// A constrained degree of freedom is enforced by setting the corresponding matrix row as the
// identity, and the associated free value as zero. This yields a trivial equation of the form
// x = 0, where x is the constrained degree of freedom.
func (s *Structure) AddDispConstraints(matrix mat.MutableMatrix, vector vec.MutableVector) {
	var (
		constraint *structure.Constraint
		dofs       [3]int
	)

	addConstraintAtDof := func(dof int) {
		matrix.SetZeroCol(dof)
		matrix.SetIdentityRow(dof)
		vector.SetZero(dof)
	}

	for _, node := range s.GetAllNodes() {
		if node.IsExternallyConstrained() {
			constraint = node.ExternalConstraint
			dofs = node.DegreesOfFreedomNum()

			if !constraint.AllowsDispX() {
				addConstraintAtDof(dofs[0])
			}
			if !constraint.AllowsDispY() {
				addConstraintAtDof(dofs[1])
			}
			if !constraint.AllowsRotation() {
				addConstraintAtDof(dofs[2])
			}
		}
	}
}
