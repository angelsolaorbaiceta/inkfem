package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

/*
ComputeElementStiffnessMatrices sets the global stiffness matrices of the given
element. Each element has a stiffness matrix between two contiguous nodes, so
in total that makes n - 1 matrices, where n is the number of nodes.
*/
func ComputeElementStiffnessMatrices(e preprocess.Element, c chan<- preprocess.Element) {
	// var trail, lead preprocess.Node
	for i := 0; i < len(e.Nodes); i++ {
		// trail = e.Nodes[i-1]
		// lead = e.Nodes[i]
	}

	c <- e
}
