package process

// SolveOptions includes configuration parameters for structural solving process.
type SolveOptions struct {
	SaveSysMatrixImage    bool
	OutputPath            string
	SafeChecks            bool
	MaxDisplacementsError float64
}