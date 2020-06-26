package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

// Element after slicing original structural element
type Element struct {
	originalElement *structure.Element
	Nodes           []Node
	globalStiffMat  []mat.ReadOnlyMatrix
}

/* ::::::::::::::: Construction ::::::::::::::: */

/*
MakeElement creates a new element given the original element and the nodes
of the sliced result.
*/
func MakeElement(originalElement *structure.Element, nodes []Node) Element {
	matrices := make([]mat.ReadOnlyMatrix, len(nodes)-1)
	return Element{originalElement, nodes, matrices}
}

/* ::::::::::::::: Properties ::::::::::::::: */

// ID returns the id of the original structural element.
func (e Element) ID() int {
	return e.originalElement.Id
}

// NodesCount returns the number of nodes in the sliced element.
func (e *Element) NodesCount() int {
	return len(e.Nodes)
}

// Geometry returns a pointer to the geometry of the original structural element.
func (e Element) Geometry() *inkgeom.Segment {
	return &e.originalElement.Geometry
}

// StartNodeID returns the id of the start node in the original structural element.
func (e Element) StartNodeID() int {
	return e.originalElement.StartNodeId
}

// EndNodeID returns the id of the end node in the original structural element.
func (e Element) EndNodeID() int {
	return e.originalElement.EndNodeId
}

// StartLink returns the link of the original structural element with the start node.
func (e Element) StartLink() *structure.Constraint {
	return e.originalElement.StartLink
}

// EndLink returns the link of the original structural element with the end node.
func (e Element) EndLink() *structure.Constraint {
	return e.originalElement.EndLink
}

// OriginalElementString returns the string representation of the original
// structural element.
func (e Element) OriginalElementString() string {
	return e.originalElement.String()
}

// Material returns the material defined in the orifinal element.
func (e Element) Material() structure.Material {
	return e.originalElement.Material()
}

// Section returns the material defined in the orifinal element.
func (e Element) Section() structure.Section {
	return e.originalElement.Section()
}

/*
StiffnessGlobalMat generates the local stiffness matrix for the element and applies
the rotation defined by the elements' geometry reference frame.
*/
func (e Element) StiffnessGlobalMat(startT, entT inkgeom.TParam) mat.ReadOnlyMatrix {
	return e.originalElement.StiffnessGlobalMat(startT, entT)
}

/* ::::::::::::::: Methods ::::::::::::::: */

/*
ComputeStiffnessMatrices sets the global stiffness matrices for this element.
Each element has a stiffness matrix between two contiguous nodes, so
in total that makes n - 1 matrices, where n is the number of nodes.
*/
func (e Element) ComputeStiffnessMatrices(c chan<- Element) {
	var trail, lead Node
	for i := 1; i < len(e.Nodes); i++ {
		trail = e.Nodes[i-1]
		lead = e.Nodes[i]
		e.globalStiffMat[i-1] = e.StiffnessGlobalMat(trail.T, lead.T)
	}

	c <- e
}

/*
GlobalStiffMatrixAt returns the global stiffness matrix at position i, that is,
between nodes i and i + 1.
*/
func (e Element) GlobalStiffMatrixAt(i int) mat.ReadOnlyMatrix {
	return e.globalStiffMat[i]
}

/* ::::::::::::::: sort.Interface ::::::::::::::: */

/*
ByGeometryPos implements sort.Interface for []Element based on the position of
the original geometry.
*/
type ByGeometryPos []Element

func (a ByGeometryPos) Len() int {
	return len(a)
}

func (a ByGeometryPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByGeometryPos) Less(i, j int) bool {
	iStart := a[i].Geometry().Start
	jStart := a[j].Geometry().Start
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].Geometry().End
	jEnd := a[j].Geometry().End
	return iEnd.Compare(jEnd) < 0
}
