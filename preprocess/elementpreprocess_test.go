package preprocess

import (
    // "fmt"
    "testing"
    "github.com/angelsolaorbaiceta/inkgeom"
    "github.com/angelsolaorbaiceta/inkfem/structure"
    "github.com/angelsolaorbaiceta/inkfem/structure/load"
)

/* Axial Member */
func TestSliceAxialMemberNodePositions(t *testing.T) {
    element := makeElementWithLoads(make([]load.Load, 0))
    slicedEl := nonSlicedElement(element)

    if len(slicedEl.Nodes) != 2 {
        t.Error("Expected element to have two nodes")
    }

    if !slicedEl.Nodes[0].Position.Equals(element.StartPoint()) {
        t.Error("First node's position was not as expected")
    }
    if !slicedEl.Nodes[1].Position.Equals(element.EndPoint()) {
        t.Error("Last node's position was not as expected")
    }
}

/* Non Axial Member */
func TestSliceNonAxialMemberNodePositions(t *testing.T) {
    element := makeElementWithLoads(make([]load.Load, 0))
    slicedEl := slicedElement(element, 2)

    if len(slicedEl.Nodes) != 3 {
        t.Error("Expected element to have three nodes")
    }

    if !slicedEl.Nodes[0].Position.Equals(element.StartPoint()) {
        t.Error("First node's position was not as expected")
    }
    if !slicedEl.Nodes[1].Position.Equals(element.PointAt(inkgeom.MakeTParam(0.5))) {
        t.Error("Middle node's position was not as expected")
    }
    if !slicedEl.Nodes[2].Position.Equals(element.EndPoint()) {
        t.Error("Last node's position was not as expected")
    }
}

/* Utils */
func makeElementWithLoads(loads []load.Load) structure.Element {
    return structure.MakeElement(
        1,
        structure.MakeFreeNodeFromProjs(1, 1.0, 2.0),
        structure.MakeFreeNodeFromProjs(2, 3.0, 4.0),
        structure.MakeDispConstraint(),
        structure.MakeDispConstraint(),
        structure.MakeUnitMaterial(),
        structure.MakeUnitSection(),
        loads)
}
