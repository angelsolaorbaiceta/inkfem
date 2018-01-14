package preprocess

import (
	"sort"
	"sync"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

/*
DoStructure preprocesses the structure by concurrently slicing each of the structural members.
*/
func DoStructure(s structure.Structure, wg *sync.WaitGroup) Structure {
	channel := make(chan Element, len(s.Elements))

	for _, element := range s.Elements {
		wg.Add(1)
		go DoElement(element, channel, wg)
	}
	wg.Wait()
	close(channel)

	var slicedElements []Element
	for slicedEl := range channel {
		slicedElements = append(slicedElements, slicedEl)
	}

	str := Structure{s.Metadata, s.Nodes, slicedElements}
	assignDof(&str)

	return str
}

func assignDof(s *Structure) {
	sort.Sort(ByGeometryPos(s.Elements))
}
