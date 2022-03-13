package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// Element after slicing the original structural element.
// Consists of a sequence of intermediate nodes with the element's loads applied to them.
type Element struct {
	*structure.Element
	nodes          []*Node
	globalStiffMat []mat.ReadOnlyMatrix
}

// MakeElement creates a new element given the original element and the nodes of the sliced result.
func MakeElement(originalElement *structure.Element, nodes []*Node) *Element {
	return &Element{
		Element:        originalElement,
		nodes:          nodes,
		globalStiffMat: make([]mat.ReadOnlyMatrix, len(nodes)-1),
	}
}

// NodesCount returns the number of nodes in the sliced element.
func (element Element) NodesCount() int {
	return len(element.nodes)
}

// Nodes is the slice of all nodes in the element.
func (element Element) Nodes() []*Node {
	return element.nodes
}

// NodeAt returns the node at a given index.
func (element Element) NodeAt(i int) *Node {
	return element.nodes[i]
}

// SetEquationTerms sets this element's stiffness and load terms into the global system of equations.
func (element *Element) SetEquationTerms(matrix mat.MutableMatrix, vector vec.MutableVector) {
	element.computeStiffnessMatrices()
	element.addTermsToStiffnessMatrix(matrix)
	element.addTermsToLoadVector(vector)
}

// computeStiffnessMatrices sets the global stiffness matrices for this element in place.
//
// Each element has a stiffness matrix between two contiguous nodes, so in total that makes n - 1
// matrices, where n is the number of nodes.
func (element *Element) computeStiffnessMatrices() {
	var trail, lead *Node

	for i := 1; i < len(element.nodes); i++ {
		trail = element.nodes[i-1]
		lead = element.nodes[i]
		element.globalStiffMat[i-1] = element.StiffnessGlobalMat(trail.T, lead.T)
	}
}

func (element *Element) addTermsToStiffnessMatrix(matrix mat.MutableMatrix) {
	var (
		stiffMat                    mat.ReadOnlyMatrix
		trailNodeDofs, leadNodeDofs [3]int
		dofs                        [6]int
		stiffVal                    float64
	)

	for i := 1; i < len(element.nodes); i++ {
		stiffMat = element.globalStiffMat[i-1]
		trailNodeDofs = element.nodes[i-1].DegreesOfFreedomNum()
		leadNodeDofs = element.nodes[i].DegreesOfFreedomNum()
		dofs = [6]int{
			trailNodeDofs[0], trailNodeDofs[1], trailNodeDofs[2],
			leadNodeDofs[0], leadNodeDofs[1], leadNodeDofs[2],
		}

		for row := 0; row < stiffMat.Rows(); row++ {
			for col := 0; col < stiffMat.Cols(); col++ {
				if stiffVal = stiffMat.Value(row, col); !nums.IsCloseToZero(stiffVal) {
					matrix.AddToValue(dofs[row], dofs[col], stiffVal)
				}
			}
		}
	}
}

func (element *Element) addTermsToLoadVector(sysVector vec.MutableVector) {
	var (
		globalTorsor *math.Torsor
		dofs         [3]int
		refFrame     = element.RefFrame()
	)

	for _, node := range element.nodes {
		globalTorsor = node.NetLocalLoadTorsor().ProjectedToGlobal(refFrame)
		dofs = node.DegreesOfFreedomNum()

		sysVector.SetValue(dofs[0], globalTorsor.Fx())
		sysVector.SetValue(dofs[1], globalTorsor.Fy())
		sysVector.SetValue(dofs[2], globalTorsor.Mz())
	}
}
