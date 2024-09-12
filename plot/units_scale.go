package plot

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const (
	// xSmallScaleFactor is the scale factor applied to the structure geometry
	// when the units appear to be in the order of meters.
	xSmallScaleFactor unitsScale = 150.0
	// smallScaleFactor is the scale factor applied to the structure geometry
	// when the units appear to be in the order of feet.
	smallScaleFactor unitsScale = 50.0
	// mediumScaleFactor is the scale factor applied to the structure geometry
	// when the units appear to be in the order of inches.
	mediumScaleFactor unitsScale = 4.0
)

// unitsScale is a scale factor that applied to the structure geometry to make
// it look nice. This is necessary due to the fact that the library svggo uses
// integer coordinates, and the bars might be too small to be visible if their
// length is in meters, for instance.
type unitsScale float64

func (s unitsScale) value() float64 {
	return float64(s)
}

func (s unitsScale) applyToPoint(p *g2d.Point) *g2d.Point {
	return g2d.MakePoint(
		p.X()*float64(s),
		p.Y()*float64(s),
	)
}

func (s unitsScale) applyToLength(l float64) float64 {
	return l * float64(s)
}

// determineUnitsScale determines the scale factor that should be applied to the
// plot to make the bars visible and proportional to the size of the structure.
//
// Depending on the length units used (cm, m, ft, in), the scale factor will be
// different. The median length of the bars is used to determine the scale factor.
//
// The optimal drawing is adjusted for cm units. If cm are detected, the scale
// factor is 1.0.
func determineUnitsScale(st *structure.Structure) unitsScale {
	stats, err := calculateBarLengthStats(st)
	if err != nil {
		panic(err)
	}

	median := stats.P50Length

	if median < xSmallMedianThreshold {
		return xSmallScaleFactor
	}

	if median < smallMedianThreshold {
		return smallScaleFactor
	}

	if median < mediumMedianThreshold {
		return mediumScaleFactor
	}

	return 1.0
}
