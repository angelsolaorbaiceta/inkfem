package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

/*
Node represents an intermediate point in a sliced element.
This point has a T Parameter associated, loads applied and degrees of freedom
numbering for the global system.
*/
type Node struct {
	T            inkgeom.TParam
	Position     g2d.Projectable
	localActions [3]float64
	globalDof    [3]int
}

/* <-- Construction --> */

/*
MakeNode creates a new node with given T parameter value, position and local
loads {fx, fy, mz}.
*/
func MakeNode(
	t inkgeom.TParam,
	position g2d.Projectable,
	fx, fy, mz float64,
) *Node {
	return &Node{t, position, [3]float64{fx, fy, mz}, [3]int{0, 0, 0}}
}

/*
MakeUnloadedNode creates a new node with given T parameter value and position.
It has no loads applied.
*/
func MakeUnloadedNode(t inkgeom.TParam, position g2d.Projectable) *Node {
	return &Node{t, position, [3]float64{}, [3]int{0, 0, 0}}
}

/* <-- Properties --> */

/*
LocalFx returns the magnitude of the local force in X. 0.0 if it has no loads
applied.
*/
func (n Node) LocalFx() float64 {
	return n.localActions[0]
}

/*
LocalFy returns the magnitude of the local force in Y. 0.0 if it has no loads
applied.
*/
func (n Node) LocalFy() float64 {
	return n.localActions[1]
}

/*
LocalMz returns the magnitude of the local moment about Z. 0.0 if it has no
loads applied.
*/
func (n Node) LocalMz() float64 {
	return n.localActions[2]
}

/*
LocalActions returns the array of local load values {Fx, Fy, Mz}.
*/
func (n Node) LocalActions() [3]float64 {
	return n.localActions
}

/*
SetDegreesOfFreedomNum adds degrees of freedom numbers to the node. These
degrees of freedom numbers are also the position in the system of equations
for the corresponding stiffness terms.
*/
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.globalDof[0] = dx
	n.globalDof[1] = dy
	n.globalDof[2] = rz
}

/*
DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
*/
func (n Node) DegreesOfFreedomNum() [3]int {
	return n.globalDof
}

/*
HasDegreesOfFreedomNum returns true if the node has already been assigned
degress of freedom.

If all DOFs are 0, it is assumed that this node hasn't been assigned DOFs.
*/
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.globalDof[0] != 0 || n.globalDof[1] != 0 || n.globalDof[2] != 0
}

/* <-- Methods --> */

/*
AddLoad adds the given load to the node applied load.
*/
func (n *Node) AddLoad(localComponents [3]float64) {
	n.localActions[0] += localComponents[0]
	n.localActions[1] += localComponents[1]
	n.localActions[2] += localComponents[2]
}

/* <-- Stringer --> */

func (n Node) String() string {
	loads := fmt.Sprintf("{%f %f %f}", n.LocalFx(), n.LocalFy(), n.LocalMz())
	return fmt.Sprintf(
		"%f: %f %f | %s | DOF: %v",
		n.T.Value(),
		n.Position.X,
		n.Position.Y,
		loads,
		n.DegreesOfFreedomNum(),
	)
}
