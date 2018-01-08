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
	localActions []float64
}

/* Construction */

// MakeNode creates a new node with given T parameter value, position and local loads {fx, fy, mz}.
func MakeNode(t inkgeom.TParam, position inkgeom.Projectable, fx, fy, mz float64) Node {
	return Node{t, position, []float64{fx, fy, mz}}
}

// MakeUnloadedNode creates a new node with given T parameter value and position. It has no loads applied.
func MakeUnloadedNode(t inkgeom.TParam, position inkgeom.Projectable) Node {
	return Node{t, position, []float64{}}
}

/* Properties */

// IsLoaded returns true if this node has loads applied to it. False otherwise.
func (n Node) IsLoaded() bool {
	return len(n.localActions) > 0
}

// LocalFx returns the magnitude of the local force in X. 0.0 if it has no loads applied.
func (n Node) LocalFx() float64 {
	if n.IsLoaded() {
		return n.localActions[0]
	}

	return 0.0
}

// LocalFy returns the magnitude of the local force in Y. 0.0 if it has no loads applied.
func (n Node) LocalFy() float64 {
	if n.IsLoaded() {
		return n.localActions[1]
	}

	return 0.0
}

// LocalMz returns the magnitude of the local moment about Z. 0.0 if it has no loads applied.
func (n Node) LocalMz() float64 {
	if n.IsLoaded() {
		return n.localActions[2]
	}

	return 0.0
}

/* Stringer */
func (n Node) String() string {
	var loads string
	if n.IsLoaded() {
		loads = fmt.Sprintf("{%f %f %f}", n.LocalFx(), n.LocalFy(), n.LocalMz())
	} else {
		loads = "{}"
	}

	return fmt.Sprintf("%f: {%f, %f} | %s", n.T.Value(), n.Position.X, n.Position.Y, loads)
}
