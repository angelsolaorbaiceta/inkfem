package plot

// plotConfig is a struct that holds the configuration options to control the
// appearance of the structural plot.
type plotConfig struct {
	GeometryColor string
	GeometryWidth int

	ExternalConstColor string
	ExternalConstWidth int
}

func defaultPlotConfig() *plotConfig {
	return &plotConfig{
		GeometryColor:      "black",
		GeometryWidth:      2,
		ExternalConstColor: "black",
		ExternalConstWidth: 2,
	}
}
