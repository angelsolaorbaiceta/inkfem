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

	ShearColor     string
	ShearFillColor string

	AxialColor     string
	AxialFillColor string

	BendingColor     string
	BendingFillColor string
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
		DistLoadArrowSize: 20,

		ShearColor:     "#FB8C00",
		ShearFillColor: "#FFA72633",

		AxialColor:     "#8E24AA",
		AxialFillColor: "#AB47BC33",

		BendingColor:     "#1E88E5",
		BendingFillColor: "#42A5F533",
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
		DistLoadArrowSize: 20,

		ShearColor:     "#FFA726",
		ShearFillColor: "#FFA72633",

		AxialColor:     "#CE93D8",
		AxialFillColor: "#CE93D833",

		BendingColor:     "#42A5F5",
		BendingFillColor: "#42A5F533",
	}
}
