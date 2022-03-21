package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
)

type ElementsSeq struct {
	elements        []*Element
	materialsByName map[string]*Material
	sectionsByName  map[string]*Section
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

// MaterialsCount is the number of different materials used in the elements.
// Two materials are considered different if their names are.
func (el *ElementsSeq) MaterialsCount() int {
	return len(el.GetMaterialsByName())
}

// GetMaterialsByName returns a map of all used materials by name.
func (el *ElementsSeq) GetMaterialsByName() map[string]*Material {
	if el.materialsByName == nil {
		el.materialsByName = make(map[string]*Material)
		var material *Material

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
func (el *ElementsSeq) GetSectionsByName() map[string]*Section {
	if el.sectionsByName == nil {
		el.sectionsByName = make(map[string]*Section)
		var section *Section

		for _, element := range el.elements {
			section = element.Section()
			el.sectionsByName[section.Name] = section
		}
	}

	return el.sectionsByName
}

// LoadsCount is the total number of concentrated and distributed loads applied to all elements.
func (el *ElementsSeq) LoadsCount() int {
	count := 0

	for _, element := range el.elements {
		count += element.LoadsCount()
	}

	return count
}

// GetAllConcentratedLoads returns a slice containing all the concentrated loads applied
// to the elements.
func (el *ElementsSeq) GetAllConcentratedLoads() []*load.ConcentratedLoad {
	var loads []*load.ConcentratedLoad

	for _, element := range el.elements {
		loads = append(loads, element.ConcentratedLoads...)
	}

	return loads
}
