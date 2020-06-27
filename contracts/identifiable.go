package contracts

import "sort"

// Identifiable is anything that can be referenced using an integer number.
type Identifiable interface {
	GetId() int
}

// ByID implements the sort.Interface for []Identifiable based in their id.
type ByID []Identifiable

func (a ByID) Len() int {
	return len(a)
}

func (a ByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByID) Less(i, j int) bool {
	return a[i].GetId() < a[j].GetId()
}

// SortById sorts (in place) a slice of identifiable elements.
func SortById(elements []Identifiable) {
	sort.Sort(ByID(elements))
}
