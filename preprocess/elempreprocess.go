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
        c <- slicedElement(e)
    }
}

/* <---------- Non Sliced ----------> */
func nonSlicedElement(e structure.Element) Element {
    return MakeElement(
        e,
        []Node{
            MakeUnloadedNode(inkgeom.MakeTParam(inkgeom.MIN_T), e.StartPoint()),
            MakeUnloadedNode(inkgeom.MakeTParam(inkgeom.MAX_T), e.EndPoint())})
}

/* <---------- Sliced ----------> */
func slicedElement(e structure.Element) Element {
    return MakeElement(e, make([]Node, 10))
}
