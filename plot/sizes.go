package plot

import (
	"fmt"
	"sort"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// barLengthsStats holds the statistics of the lengths of the bars in a structure.
type barLengthsStats struct {
	// MinLength is the length of the shorter bar in the structure.
	MinLength float64
	// MaxLength is the length of the longer bar in the structure.
	MaxLength float64
	// AvgLength is the average length of the bars in the structure.
	AvgLength float64
	// P50Length is the length of the bar that is in the middle of the sorted
	// list of lengths. The value is exact if the number of bars is odd, and the
	// average of the two middle values if the number of bars is even.
	P50Length float64
}

const (
	// xSmallMedianThreshold is the threshold that determines if the median of the
	// bar lengths is considered small. If the median is smaller than this value,
	// the scale of the plot needs to be adjusted.
	//
	// A median below this value suggest the units in which the structure is defined
	// are meters.
	xSmallMedianThreshold = 5.0

	// smallMedianThreshold is the threshold that determines if the median of the
	// bar lengths is considered small. If the median is smaller than this value,
	// the scale of the plot needs to be adjusted.
	//
	// A median below this value suggest the units in which the structure is defined
	// are feet.
	smallMedianThreshold = 17.0

	// mediumMedianThreshold is the threshold that determines if the median of the
	// bar lengths is considered medium. If the median is smaller than this value,
	// the scale of the plot needs to be adjusted.
	//
	// A median below this value suggest the units in which the structure is defined
	// are inches.
	mediumMedianThreshold = 200.0
)

// calculateBarLengthStats calculates the statistics of the lengths of the bars
// in the structure. It includes the minimum, maximum, average, and median lengths.
//
// These values can be used to determine the scale of the plot, for instance.
func calculateBarLengthStats(st *structure.Structure) (*barLengthsStats, error) {
	if st.ElementsCount() == 0 {
		return nil, fmt.Errorf("structure has no elements")
	}

	var (
		totalLength float64 = 0
		count               = st.ElementsCount()
		lengths             = make([]float64, count)
	)

	for i, bar := range st.Elements() {
		length := bar.Length()

		totalLength += length
		lengths[i] = length
	}

	sort.Float64s(lengths)

	return &barLengthsStats{
		MinLength: lengths[0],
		MaxLength: lengths[count-1],
		AvgLength: totalLength / float64(count),
		P50Length: lengthsMedian(lengths),
	}, nil
}

func lengthsMedian(lengths []float64) float64 {
	count := len(lengths)

	if count%2 == 0 {
		return (lengths[count/2-1] + lengths[count/2]) / 2
	}

	return lengths[count/2]
}
