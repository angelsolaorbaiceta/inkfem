package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const unsetDOF = -1

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
	externalLocalLoad *math.Torsor
	leftLocalLoad     *math.Torsor
	rightLocalLoad    *math.Torsor
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
		externalLocalLoad: math.MakeTorsor(fx, fy, mz),
		leftLocalLoad:     math.MakeNilTorsor(),
		rightLocalLoad:    math.MakeNilTorsor(),
		dofs:              [3]int{unsetDOF, unsetDOF, unsetDOF},
	}
}

// MakeUnloadedNode creates a new node with given T parameter value, position, and no loads applied.
func MakeUnloadedNode(t inkgeom.TParam, position g2d.Projectable) *Node {
	return &Node{
		T:                 t,
		Position:          position,
		externalLocalLoad: math.MakeNilTorsor(),
		leftLocalLoad:     math.MakeNilTorsor(),
		rightLocalLoad:    math.MakeNilTorsor(),
		dofs:              [3]int{unsetDOF, unsetDOF, unsetDOF},
	}
}

// NetLocalFx returns the magnitude of the local force in X.
func (n Node) NetLocalFx() float64 {
	return n.externalLocalLoad.Fx() + n.leftLocalLoad.Fx() + n.rightLocalLoad.Fx()
}

func (n Node) LocalLeftFx() float64 {
	return n.leftLocalLoad.Fx()
}

func (n Node) LocalRightFx() float64 {
	return n.rightLocalLoad.Fx()
}

// NetLocalFy returns the magnitude of the net local force in Y.
func (n Node) NetLocalFy() float64 {
	return n.externalLocalLoad.Fy() + n.LocalLeftFy() + n.LocalRightFy()
}

func (n Node) LocalLeftFy() float64 {
	return n.leftLocalLoad.Fy()
}

func (n Node) LocalRightFy() float64 {
	return n.rightLocalLoad.Fy()
}

// NetLocalMz returns the magnitude of the local moment about Z.
func (n Node) NetLocalMz() float64 {
	return n.externalLocalLoad.Mz() + n.LocalLeftMz() + n.LocalRightMz()
}

func (n Node) LocalLeftMz() float64 {
	return n.leftLocalLoad.Mz()
}

func (n Node) LocalRightMz() float64 {
	return n.rightLocalLoad.Mz()
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
	n.externalLocalLoad = n.externalLocalLoad.PlusComponents(fx, fy, mz)
}

/*
AddLocalLeftLoad adds the given load values to the load applied from the finite element where this
node is to the left of it (where this node is the element's trailing node).
*/
func (n *Node) AddLocalLeftLoad(fx, fy, mz float64) {
	n.leftLocalLoad = n.leftLocalLoad.PlusComponents(fx, fy, mz)
}

/*
AddLocalRightLoad adds the given load values to the load applied from the finite element where this
node is to the right of it (where this node is the element's leading node).
*/
func (n *Node) AddLocalRightLoad(fx, fy, mz float64) {
	n.rightLocalLoad = n.rightLocalLoad.PlusComponents(fx, fy, mz)
}

func (n Node) String() string {
	var (
		leftLoads  = fmt.Sprintf("{%f %f %f}", n.LocalLeftFx(), n.LocalLeftFy(), n.LocalLeftMz())
		rightLoads = fmt.Sprintf("{%f %f %f}", n.LocalRightFx(), n.LocalRightFy(), n.LocalRightMz())
		loads      = fmt.Sprintf("{%f %f %f}", n.NetLocalFx(), n.NetLocalFy(), n.NetLocalMz())
	)

	return fmt.Sprintf(
		"%f : %f %f \n\t left  : %s \n\t right : %s \n\t net   : %s \n\t dof   : %v",
		n.T.Value(),
		n.Position.X,
		n.Position.Y,
		leftLoads,
		rightLoads,
		loads,
		n.DegreesOfFreedomNum(),
	)
}
