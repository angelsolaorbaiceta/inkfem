package utils

import "sort"

type Identifiable interface {
	GetId() int
}

// ByID implements the sort.Interface for []Identifiable based in their id
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

func SortById(elements []Identifiable) {
	sort.Sort(ByID(elements))
}
