package utils

type Identifiable interface {
	Id() int
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
	return a[i].Id() < a[j].Id()
}
