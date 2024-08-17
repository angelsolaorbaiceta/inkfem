package plot

// PlotConfig is a struct that holds the configuration options to control the
// appearance of the structural plot.
type PlotConfig struct {
	GeometryColor string
	GeometryWidth int

	ExternalConstColor string
	ExternalConstWidth int

	NodeRadius       int
	ConstraintLength int

	DistLoadColor     string
	DistLoadFillColor string
	DistLoadWidth     int
	DistLoadArrowSize int
}

func DefaultPlotConfig() *PlotConfig {
	return &PlotConfig{
		GeometryColor: "black",
		GeometryWidth: 2,

		ExternalConstColor: "black",
		ExternalConstWidth: 2,

		NodeRadius:       10,
		ConstraintLength: 80,

		DistLoadColor:     "#558B2F",
		DistLoadFillColor: "#9CCC6533",
		DistLoadWidth:     1,
		DistLoadArrowSize: 30,
	}
}

func DarkPlotConfig() *PlotConfig {
	return &PlotConfig{
		GeometryColor: "#FAFAFA",
		GeometryWidth: 2,

		ExternalConstColor: "#FAFAFA",
		ExternalConstWidth: 2,

		NodeRadius:       10,
		ConstraintLength: 80,

		DistLoadColor:     "#9CCC65",
		DistLoadFillColor: "#9CCC6533",
		DistLoadWidth:     1,
		DistLoadArrowSize: 30,
	}
}
