package preprocess

import (
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

	return Structure{s.Metadata, s.Nodes, slicedElements}
}
