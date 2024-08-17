package plot

// plotConfig is a struct that holds the configuration options to control the
// appearance of the structural plot.
type plotConfig struct {
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

func defaultPlotConfig() *plotConfig {
	return &plotConfig{
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
