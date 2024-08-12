package plot

// plotConfig is a struct that holds the configuration options to control the
// appearance of the structural plot.
type plotConfig struct {
	GeometryColor string
	GeometryWidth int
}

func defaultPlotConfig() plotConfig {
	return plotConfig{
		GeometryColor: "black",
		GeometryWidth: 1,
	}
}
