/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const unsetDOF = -1

/*
A Node represents an intermediate point in a sliced element.

This point has a T Parameter associated, loads applied and degrees of freedom
numbering for the global system.

The `leftLocalLoad` is the equivalent load, in local coordinates, from the finite
element located to the left of the node.

The `rightLocalLoad` is the equivalent load, in local coordinates, from the finite
element located to the right of the node.
*/
type Node struct {
	T                 inkgeom.TParam
	Position          g2d.Projectable
	externalLocalLoad [3]float64
	leftLocalLoad     [3]float64
	rightLocalLoad    [3]float64
	dofs              [3]int
}

/* <-- Construction --> */

/*
MakeNode creates a new node with given T parameter value, position and local
external loads {fx, fy, mz}.
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

/*
MakeUnloadedNode creates a new node with given T parameter value, position,
and no loads applied.
*/
func MakeUnloadedNode(t inkgeom.TParam, position g2d.Projectable) *Node {
	return &Node{t, position, [3]float64{}, [3]float64{}, [3]float64{}, [3]int{unsetDOF, unsetDOF, unsetDOF}}
}

/* <-- Properties --> */

/*
NetLocalFx returns the magnitude of the local force in X. 0.0 if it has no loads applied.
*/
func (n Node) NetLocalFx() float64 {
	return n.externalLocalLoad[0] + n.leftLocalLoad[0] + n.rightLocalLoad[0]
}

/*
NetLocalFy returns the magnitude of the local force in Y. 0.0 if it has no loads applied.
*/
func (n Node) NetLocalFy() float64 {
	return n.externalLocalLoad[1] + n.LocalLeftFy() + n.LocalRightFy()
}

func (n Node) LocalLeftFy() float64 {
	return n.leftLocalLoad[1]
}

func (n Node) LocalRightFy() float64 {
	return n.rightLocalLoad[1]
}

/*
NetLocalMz returns the magnitude of the local moment about Z. 0.0 if it has no loads applied.
*/
func (n Node) NetLocalMz() float64 {
	return n.externalLocalLoad[2] + n.leftLocalLoad[2] + n.rightLocalLoad[2]
}

/*
NetLocalLoadVector returns the array of net local load values {Fx, Fy, Mz}.
*/
func (n Node) NetLocalLoadVector() [3]float64 {
	return [3]float64{
		n.NetLocalFx(),
		n.NetLocalFy(),
		n.NetLocalMz(),
	}
}

/*
SetDegreesOfFreedomNum adds degrees of freedom numbers to the node.

These degrees of freedom numbers are also the position in the system of equations
for the corresponding stiffness terms.
*/
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.dofs[0] = dx
	n.dofs[1] = dy
	n.dofs[2] = rz
}

/*
DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
*/
func (n Node) DegreesOfFreedomNum() [3]int {
	return n.dofs
}

/*
HasDegreesOfFreedomNum returns true if the node has already been assigned degress of
freedom.

If any of the DOFs is -1, it's assumed that this node hasn't been assigned DOFs.
*/
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.dofs[0] != unsetDOF ||
		n.dofs[1] != unsetDOF ||
		n.dofs[2] != unsetDOF
}

/* <-- Methods --> */

/*
DistanceTo computes the distance between this an other node.
*/
func (n *Node) DistanceTo(other *Node) float64 {
	return n.Position.DistanceTo(other.Position)
}

/*
AddLocalExternalLoad adds the given load values to the load applied from the
left finite element.
*/
func (n *Node) AddLocalExternalLoad(fx, fy, mz float64) {
	n.externalLocalLoad[0] += fx
	n.externalLocalLoad[1] += fy
	n.externalLocalLoad[2] += mz
}

/*
AddLocalLeftLoad adds the given load values to the load applied from the
left finite element.
*/
func (n *Node) AddLocalLeftLoad(fx, fy, mz float64) {
	n.leftLocalLoad[0] += fx
	n.leftLocalLoad[1] += fy
	n.leftLocalLoad[2] += mz
}

/*
AddLocalRightLoad adds the given load values to the load applied from the
right finite element.
*/
func (n *Node) AddLocalRightLoad(fx, fy, mz float64) {
	n.rightLocalLoad[0] += fx
	n.rightLocalLoad[1] += fy
	n.rightLocalLoad[2] += mz
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
