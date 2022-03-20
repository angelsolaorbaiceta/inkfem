package preprocess

// ByGeometryPos implements sort.Interface for []Element based on the position of the
// original geometry.
type ByGeometryPos []*Element

func (a ByGeometryPos) Len() int {
	return len(a)
}

func (a ByGeometryPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByGeometryPos) Less(i, j int) bool {
	iStart := a[i].StartPoint()
	jStart := a[j].StartPoint()
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].EndPoint()
	jEnd := a[j].EndPoint()
	return iEnd.Compare(jEnd) < 0
}
