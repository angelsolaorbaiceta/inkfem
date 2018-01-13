package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
)

/*
Node represents an intermediate point in a sliced element.
This point has a T Parameter associated, loads applied and degrees of freedom
numbering for the global system.
*/
type Node struct {
	T            inkgeom.TParam
	Position     inkgeom.Projectable
	localActions [3]float64
}

/* Construction */

// MakeNode creates a new node with given T parameter value, position and local loads {fx, fy, mz}.
func MakeNode(t inkgeom.TParam, position inkgeom.Projectable, fx, fy, mz float64) Node {
	return Node{t, position, [3]float64{fx, fy, mz}}
}

// MakeUnloadedNode creates a new node with given T parameter value and position. It has no loads applied.
func MakeUnloadedNode(t inkgeom.TParam, position inkgeom.Projectable) Node {
	return Node{t, position, [3]float64{}}
}

/* Properties */

// LocalFx returns the magnitude of the local force in X. 0.0 if it has no loads applied.
func (n Node) LocalFx() float64 {
	return n.localActions[0]
}

// LocalFy returns the magnitude of the local force in Y. 0.0 if it has no loads applied.
func (n Node) LocalFy() float64 {
	return n.localActions[1]
}

// LocalMz returns the magnitude of the local moment about Z. 0.0 if it has no loads applied.
func (n Node) LocalMz() float64 {
	return n.localActions[2]
}

// AddLoad adds the given load to the node applied load.
func (n *Node) AddLoad(localComponents [3]float64) {
	n.localActions[0] += localComponents[0]
	n.localActions[1] += localComponents[1]
	n.localActions[2] += localComponents[2]
}

/* Stringer */
func (n Node) String() string {
	loads := fmt.Sprintf("{%f %f %f}", n.LocalFx(), n.LocalFy(), n.LocalMz())
	return fmt.Sprintf("%f: {%f, %f} | %s", n.T.Value(), n.Position.X, n.Position.Y, loads)
}
