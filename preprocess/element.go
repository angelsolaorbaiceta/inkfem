package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// Element after slicing original structural element.
type Element struct {
	*structure.Element
	Nodes          []*Node
	globalStiffMat []mat.ReadOnlyMatrix
}

// MakeElement creates a new element given the original element and the nodes of the sliced result.
func MakeElement(originalElement *structure.Element, nodes []*Node) *Element {
	matrices := make([]mat.ReadOnlyMatrix, len(nodes)-1)
	return &Element{originalElement, nodes, matrices}
}

// NodesCount returns the number of nodes in the sliced element.
func (element Element) NodesCount() int {
	return len(element.Nodes)
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

	for i := 1; i < len(element.Nodes); i++ {
		trail = element.Nodes[i-1]
		lead = element.Nodes[i]
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

	for i := 1; i < len(element.Nodes); i++ {
		stiffMat = element.globalStiffMatrixAt(i - 1)
		trailNodeDofs = element.Nodes[i-1].DegreesOfFreedomNum()
		leadNodeDofs = element.Nodes[i].DegreesOfFreedomNum()
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

// globalStiffMatrixAt returns the global stiffness matrix at position i, that is,
// between nodes i and i + 1.
func (e Element) globalStiffMatrixAt(i int) mat.ReadOnlyMatrix {
	return e.globalStiffMat[i]
}

func (element *Element) addTermsToLoadVector(sysVector vec.MutableVector) {
	var (
		globalTorsor *math.Torsor
		dofs         [3]int
		refFrame     = element.RefFrame()
	)

	for _, node := range element.Nodes {
		globalTorsor = node.NetLocalLoadTorsor().ProjectedToGlobal(refFrame)
		dofs = node.DegreesOfFreedomNum()

		sysVector.SetValue(dofs[0], globalTorsor.Fx())
		sysVector.SetValue(dofs[1], globalTorsor.Fy())
		sysVector.SetValue(dofs[2], globalTorsor.Mz())
	}
}

// ByGeometryPos implements sort.Interface for []Element based on the position of the
// original geometry.
type ByGeometryPos []*Element

func (a ByGeometryPos) Len() int {
	return len(a)
}

func (a ByGeometryPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByGeometryPos) Less(i, j int) bool {
	iStart := a[i].StartPoint()
	jStart := a[j].StartPoint()
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].EndPoint()
	jEnd := a[j].EndPoint()
	return iEnd.Compare(jEnd) < 0
}
