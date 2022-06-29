package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

const unsetDOF = -1

// A Node represents an intermediate point in a sliced element.
//
// This point has a T Parameter associated, external loads applied and degrees of freedom
// numbering for the global system.
//
// The "leftLocalLoad" is the equivalent load, in local coordinates, from the finite element located
// to the left of the node.
//
// The "rightLocalLoad" is the equivalent load, in local coordinates, from the finite element located
// to the right of the node.
type Node struct {
	T                 nums.TParam
	Position          *g2d.Point
	externalLocalLoad *math.Torsor
	leftLocalLoad     *math.Torsor
	rightLocalLoad    *math.Torsor
	dofs              [3]int
}

// MakeNode creates a new node with given T parameter value, position and local external
// loads {fx, fy, mz}.
func MakeNode(
	t nums.TParam,
	position *g2d.Point,
	fx, fy, mz float64,
) *Node {
	return &Node{
		T:                 t,
		Position:          position,
		externalLocalLoad: math.MakeTorsor(fx, fy, mz),
		leftLocalLoad:     math.MakeNilTorsor(),
		rightLocalLoad:    math.MakeNilTorsor(),
		dofs:              [3]int{unsetDOF, unsetDOF, unsetDOF},
	}
}

// MakeNodeWithDofs creates a new node with given T parameter value, position and local
// external loads {fx, fy, mz} and the given degrees of freedom numbers.
//
// The degrees of freedom numbers in a preprocessed nodes are typically set after the
// node has been created, as they are assigned at the structure level; thus, it's
// necessary to have all nodes before they are assigned degrees of freedom numbers.
//
// This factory function is mostly used for testing purposes.
func MakeNodeWithDofs(
	t nums.TParam,
	position *g2d.Point,
	fx, fy, mz float64,
	dofs [3]int,
) *Node {
	return &Node{
		T:                 t,
		Position:          position,
		externalLocalLoad: math.MakeTorsor(fx, fy, mz),
		leftLocalLoad:     math.MakeNilTorsor(),
		rightLocalLoad:    math.MakeNilTorsor(),
		dofs:              dofs,
	}
}

// MakeUnloadedNode creates a new node with given T parameter value, position, and no loads applied.
func MakeUnloadedNode(t nums.TParam, position *g2d.Point) *Node {
	return &Node{
		T:                 t,
		Position:          position,
		externalLocalLoad: math.MakeNilTorsor(),
		leftLocalLoad:     math.MakeNilTorsor(),
		rightLocalLoad:    math.MakeNilTorsor(),
		dofs:              [3]int{unsetDOF, unsetDOF, unsetDOF},
	}
}

// NetLocalTorsor returns the resulting torsor of adding the external loads and the loads added
// by the left and right finite elements, all projected in local coordinates.
func (n Node) NetLocalTorsor() *math.Torsor {
	return n.externalLocalLoad.Plus(n.leftLocalLoad).Plus(n.rightLocalLoad)
}

// NetLocalFx returns the magnitude of the net force in X, projected in local coordinates.
func (n Node) NetLocalFx() float64 {
	return n.externalLocalLoad.Fx() + n.leftLocalLoad.Fx() + n.rightLocalLoad.Fx()
}

func (n Node) LocalLeftFx() float64 {
	return n.leftLocalLoad.Fx()
}

func (n Node) LocalRightFx() float64 {
	return n.rightLocalLoad.Fx()
}

// NetLocalFy returns the magnitude of the net force in Y, projected in local coordinates.
func (n Node) NetLocalFy() float64 {
	return n.externalLocalLoad.Fy() + n.LocalLeftFy() + n.LocalRightFy()
}

func (n Node) LocalLeftFy() float64 {
	return n.leftLocalLoad.Fy()
}

func (n Node) LocalRightFy() float64 {
	return n.rightLocalLoad.Fy()
}

// NetLocalMz returns the magnitude of the net moment about Z, projected in local coordinates.
func (n Node) NetLocalMz() float64 {
	return n.externalLocalLoad.Mz() + n.LocalLeftMz() + n.LocalRightMz()
}

func (n Node) LocalLeftMz() float64 {
	return n.leftLocalLoad.Mz()
}

func (n Node) LocalRightMz() float64 {
	return n.rightLocalLoad.Mz()
}

// NetLocalLoadTorsor returns the torsor of net load values {Fx, Fy, Mz} projected in
// local coordinates.
func (n Node) NetLocalLoadTorsor() *math.Torsor {
	return math.MakeTorsor(
		n.NetLocalFx(),
		n.NetLocalFy(),
		n.NetLocalMz(),
	)
}

// SetDegreesOfFreedomNum adds degrees of freedom numbers to the node.
// These degrees of freedom numbers are also the position in the system of equations for the
// corresponding stiffness terms.
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.dofs[0] = dx
	n.dofs[1] = dy
	n.dofs[2] = rz
}

// DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
// Panics if the degrees of freedom haven't been set.
func (n Node) DegreesOfFreedomNum() [3]int {
	if n.dofs[0] == unsetDOF || n.dofs[1] == unsetDOF || n.dofs[2] == unsetDOF {
		panic("Degrees of freedom not set for preprocessed node")
	}

	return n.dofs
}

// HasDegreesOfFreedomNum returns true if the node has already been assigned degress of freedom.
// If any of the DOFs is -1, it's assumed that this node hasn't been assigned DOFs.
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.dofs[0] != unsetDOF ||
		n.dofs[1] != unsetDOF ||
		n.dofs[2] != unsetDOF
}

// DistanceTo computes the distance between this an other node.
func (n *Node) DistanceTo(other *Node) float64 {
	return n.Position.DistanceTo(other.Position)
}

// AddLocalExternalLoad adds the given load values to the externally applied load.
func (n *Node) AddLocalExternalLoad(loadTorsor *math.Torsor) {
	n.externalLocalLoad = n.externalLocalLoad.Plus(loadTorsor)
}

// AddLocalLeftLoad adds the given load values to the load applied from the finite element where this
// node is to the left of it (where this node is the element's trailing node).
func (n *Node) AddLocalLeftLoad(fx, fy, mz float64) {
	n.leftLocalLoad = n.leftLocalLoad.PlusComponents(fx, fy, mz)
}

// AddLocalRightLoad adds the given load values to the load applied from the finite element where this
// node is to the right of it (where this node is the element's leading node).
func (n *Node) AddLocalRightLoad(fx, fy, mz float64) {
	n.rightLocalLoad = n.rightLocalLoad.PlusComponents(fx, fy, mz)
}

// Equals returns true if this and the other node equal.
func (n *Node) Equals(other *Node) bool {
	return n.T.Equals(other.T) &&
		n.Position.Equals(other.Position) &&
		n.externalLocalLoad.Equals(other.externalLocalLoad) &&
		n.leftLocalLoad.Equals(other.leftLocalLoad) &&
		n.rightLocalLoad.Equals(other.rightLocalLoad) &&
		n.dofs[0] == other.dofs[0] &&
		n.dofs[1] == other.dofs[1] &&
		n.dofs[2] == other.dofs[2]
}

// String representation of the node.
// This method is used for serialization, thus if the format is changed, the preprocessed
// file format might be affected.
func (n Node) String() string {
	return fmt.Sprintf(
		"%f : %f %f\n\text   : %s\n\tleft  : %s\n\tright : %s\n\tnet   : %s\n\tdof   : %v",
		n.T.Value(),
		n.Position.X(),
		n.Position.Y(),
		n.externalLocalLoad,
		n.leftLocalLoad,
		n.rightLocalLoad,
		n.NetLocalLoadTorsor(),
		n.DegreesOfFreedomNum(),
	)
}
