package preprocess

import "github.com/angelsolaorbaiceta/inkfem/structure"

type ElementsSeq struct {
	elements        []*Element
	materialsByName map[string]*structure.Material
}

// ElementsCount returns the number of elements in the original structure.
func (el *ElementsSeq) ElementsCount() int {
	return len(el.elements)
}

// Elements returns a slice containing all elements.
func (el *ElementsSeq) Elements() []*Element {
	return el.elements
}

// MaterialsCount is the number of different materials in the structure.
// Two materials are considered different if their names are.
func (el *ElementsSeq) MaterialsCount() int {
	return len(el.GetMaterialsByName())
}

// GetMaterialsByName returns a map of materials by material name.
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
