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

	DistLoadColor string
	DistLoadWidth int
}

func defaultPlotConfig() *plotConfig {
	return &plotConfig{
		GeometryColor: "black",
		GeometryWidth: 2,

		ExternalConstColor: "black",
		ExternalConstWidth: 2,

		NodeRadius:       10,
		ConstraintLength: 80,

		DistLoadColor: "red",
		DistLoadWidth: 1,
	}
}
