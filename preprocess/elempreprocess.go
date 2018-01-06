package preprocess

import (
    "sync"
    "github.com/angelsolaorbaiceta/inkgeom"
    "github.com/angelsolaorbaiceta/inkfem/structure"
)

// Preprocess the structural element subdividing it as corresponds.
func DoElement(e structure.Element, c chan Element, wg *sync.WaitGroup) {
    defer wg.Done()

    if e.IsAxialMember() {
        c <- nonSlicedElement(e)
    } else {
        c <- slicedElement(e, 12)
    }
}

/* <---------- Non Sliced ----------> */
func nonSlicedElement(e structure.Element) Element {
    return MakeElement(
        e,
        []Node{ // TODO: add loads
            MakeUnloadedNode(inkgeom.MIN_T, e.StartPoint()),
            MakeUnloadedNode(inkgeom.MAX_T, e.EndPoint()),
        })
}

/* <---------- Sliced ----------> */
func slicedElement(e structure.Element, times int) Element {
    tPos := inkgeom.SubTParamCompleteRangeTimes(times)

    nodes := make([]Node, len(tPos))
    for i := 0; i < len(tPos); i++ { // TODO: add loads
        nodes[i] = MakeUnloadedNode(tPos[i], e.PointAt(tPos[i]))
    }

    return MakeElement(e, nodes)
}
