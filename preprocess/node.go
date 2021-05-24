package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const (
	unsetDOF = -1
	fxIndex  = 0
	fyIndex  = 1
	mzIndex  = 2
)

/*
A Node represents an intermediate point in a sliced element.

This point has a T Parameter associated, loads applied and degrees of freedom numbering for
the global system.

The `leftLocalLoad` is the equivalent load, in local coordinates, from the finite element located
to the left of the node.

The `rightLocalLoad` is the equivalent load, in local coordinates, from the finite element located
to the right of the node.
*/
type Node struct {
	T                 inkgeom.TParam
	Position          g2d.Projectable
	externalLocalLoad [3]float64
	leftLocalLoad     [3]float64
	rightLocalLoad    [3]float64
	dofs              [3]int
}

/*
MakeNode creates a new node with given T parameter value, position and local external
loads {fx, fy, mz}.
*/
func MakeNode(
	t inkgeom.TParam,
	position g2d.Projectable,
	fx, fy, mz float64,
) *Node {
	return &Node{
		T:                 t,
		Position:          position,
		externalLocalLoad: [3]float64{fx, fy, mz},
		leftLocalLoad:     [3]float64{0, 0, 0},
		rightLocalLoad:    [3]float64{0, 0, 0},
		dofs:              [3]int{unsetDOF, unsetDOF, unsetDOF},
	}
}

// MakeUnloadedNode creates a new node with given T parameter value, position, and no loads applied.
func MakeUnloadedNode(t inkgeom.TParam, position g2d.Projectable) *Node {
	return &Node{t, position, [3]float64{}, [3]float64{}, [3]float64{}, [3]int{unsetDOF, unsetDOF, unsetDOF}}
}

// NetLocalFx returns the magnitude of the local force in X.
func (n Node) NetLocalFx() float64 {
	return n.externalLocalLoad[fxIndex] + n.leftLocalLoad[fxIndex] + n.rightLocalLoad[fxIndex]
}

// NetLocalFy returns the magnitude of the net local force in Y.
func (n Node) NetLocalFy() float64 {
	return n.externalLocalLoad[fyIndex] + n.LocalLeftFy() + n.LocalRightFy()
}

/*
LocalLeftFy returns the magnitude of the local force in Y coming from the finite element
to the left of the node.
*/
func (n Node) LocalLeftFy() float64 {
	return n.leftLocalLoad[fyIndex]
}

/*
LocalRightFy returns the magnitude of the local force in Y coming from the finite element
to the right of the node.
*/
func (n Node) LocalRightFy() float64 {
	return n.rightLocalLoad[fyIndex]
}

// NetLocalMz returns the magnitude of the local moment about Z.
func (n Node) NetLocalMz() float64 {
	return n.externalLocalLoad[mzIndex] + n.LocalLeftMz() + n.LocalRightMz()
}

/*
LocalLeftMz returns the magnitude of the local moment around Z coming from the finite element
to the left of the node.
*/
func (n Node) LocalLeftMz() float64 {
	return n.leftLocalLoad[mzIndex]
}

/*
LocalRightMz returns the magnitude of the local moment around Z coming from the finite element
to the right of the node.
*/
func (n Node) LocalRightMz() float64 {
	return n.rightLocalLoad[mzIndex]
}

// NetLocalLoadVector returns the array of net local load values {Fx, Fy, Mz}.
func (n Node) NetLocalLoadVector() [3]float64 {
	return [3]float64{
		n.NetLocalFx(),
		n.NetLocalFy(),
		n.NetLocalMz(),
	}
}

/*
SetDegreesOfFreedomNum adds degrees of freedom numbers to the node.

These degrees of freedom numbers are also the position in the system of equations for the
corresponding stiffness terms.
*/
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.dofs[0] = dx
	n.dofs[1] = dy
	n.dofs[2] = rz
}

// DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
func (n Node) DegreesOfFreedomNum() [3]int {
	return n.dofs
}

/*
HasDegreesOfFreedomNum returns true if the node has already been assigned degress of freedom.

If any of the DOFs is -1, it's assumed that this node hasn't been assigned DOFs.
*/
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.dofs[0] != unsetDOF ||
		n.dofs[1] != unsetDOF ||
		n.dofs[2] != unsetDOF
}

// DistanceTo computes the distance between this an other node.
func (n *Node) DistanceTo(other *Node) float64 {
	return n.Position.DistanceTo(other.Position)
}

// AddLocalExternalLoad adds the given load values to the load applied from the left finite element.
func (n *Node) AddLocalExternalLoad(fx, fy, mz float64) {
	n.externalLocalLoad[fxIndex] += fx
	n.externalLocalLoad[fyIndex] += fy
	n.externalLocalLoad[mzIndex] += mz
}

/*
AddLocalLeftLoad adds the given load values to the load applied from the finite element where this
node is to the left of it (where this node is the element's trailing node).
*/
func (n *Node) AddLocalLeftLoad(fx, fy, mz float64) {
	n.leftLocalLoad[fxIndex] += fx
	n.leftLocalLoad[fyIndex] += fy
	n.leftLocalLoad[mzIndex] += mz
}

/*
AddLocalRightLoad adds the given load values to the load applied from the finite element where this
node is to the right of it (where this node is the element's leading node).
*/
func (n *Node) AddLocalRightLoad(fx, fy, mz float64) {
	n.rightLocalLoad[fxIndex] += fx
	n.rightLocalLoad[fyIndex] += fy
	n.rightLocalLoad[mzIndex] += mz
}

/* <-- Stringer --> */

func (n Node) String() string {
	loads := fmt.Sprintf("{%f %f %f}", n.NetLocalFx(), n.NetLocalFy(), n.NetLocalMz())
	return fmt.Sprintf(
		"%f: %f %f | %s | DOF: %v",
		n.T.Value(),
		n.Position.X,
		n.Position.Y,
		loads,
		n.DegreesOfFreedomNum(),
	)
}
