package preprocess

import (
	"github.com/angelsolaorbaiceta/inkgeom"
)

type Node struct {
	T inkgeom.TParam
	Position inkgeom.Projectable
    localActions []float64
}

/* Construction */
func MakeNode(t inkgeom.TParam, position inkgeom.Projectable, fx, fy, mz float64) Node {
    return Node{t, position, []float64{fx, fy, mz}}
}

func MakeUnloadedNode(t inkgeom.TParam, position inkgeom.Projectable) Node {
	return Node{t, position, []float64{}}
}

/* Properties */
func (n Node) IsLoaded() bool {
	return len(n.localActions) > 0
}

func (n Node) LocalFx() float64 {
	if n.IsLoaded() {
		return n.localActions[0]
	}

	return 0.0
}

func (n Node) LocalFy() float64 {
	if n.IsLoaded() {
		return n.localActions[1]
	}

	return 0.0
}

func (n Node) LocalMz() float64 {
	if n.IsLoaded() {
		return n.localActions[2]
	}

	return 0.0
}
