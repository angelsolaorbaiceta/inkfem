package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
)

type ElementsSeq struct {
	elements        []*Element
	materialsByName map[string]*Material
}

// ElementsCount is the number of elements in the structure.
func (el *ElementsSeq) ElementsCount() int {
	return len(el.elements)
}

// Elements returns a slice containing all elements.
func (el *ElementsSeq) Elements() []*Element {
	return el.elements
}

// GetElementById returns the element with the given id or panics.
// This operation has an O(n) time complexity as it needs to iterate over all elements.
func (el *ElementsSeq) GetElementById(id contracts.StrID) *Element {
	for _, element := range el.elements {
		if element.GetID() == id {
			return element
		}
	}

	panic(fmt.Sprintf("Can't find element with id %s", id))
}

// MaterialsCount is the number of different materials in the structure.
// Two materials are considered different if their names are.
func (el *ElementsSeq) MaterialsCount() int {
	return len(el.GetMaterialsByName())
}

// GetMaterialsByName returns a map of materials by material name.
func (el *ElementsSeq) GetMaterialsByName() map[string]*Material {
	if el.materialsByName == nil {
		el.materialsByName = make(map[string]*Material)

		for _, element := range el.elements {
			el.materialsByName[element.material.Name] = element.material
		}
	}

	return el.materialsByName
}
