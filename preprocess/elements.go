package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

type ElementsSeq struct {
	elements        []*Element
	materialsByName map[string]*structure.Material
	sectionsByName  map[string]*structure.Section
}

// ElementsCount returns the number of elements in the original structure.
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

// MaterialsCount is the number of different materials used in the elements.
// Two materials are considered different if their names are.
func (el *ElementsSeq) MaterialsCount() int {
	return len(el.GetMaterialsByName())
}

// GetMaterialsByName returns a map of all used materials by name.
func (el *ElementsSeq) GetMaterialsByName() map[string]*structure.Material {
	if el.materialsByName == nil {
		el.materialsByName = make(map[string]*structure.Material)
		var material *structure.Material

		for _, element := range el.elements {
			material = element.Material()
			el.materialsByName[material.Name] = material
		}
	}

	return el.materialsByName
}

// SectionsCount is the number of different sections used in the elements.
// Two sections are considered different if their names are.
func (el *ElementsSeq) SectionsCount() int {
	return len(el.GetSectionsByName())
}

// GetSectionsByName returns a map of all used sections by name.
func (el *ElementsSeq) GetSectionsByName() map[string]*structure.Section {
	if el.sectionsByName == nil {
		el.sectionsByName = make(map[string]*structure.Section)
		var section *structure.Section

		for _, element := range el.elements {
			section = element.Section()
			el.sectionsByName[section.Name] = section
		}
	}

	return el.sectionsByName
}
