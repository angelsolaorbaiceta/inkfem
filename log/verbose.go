package log

import (
	"fmt"
	"time"
)

var (
	isVerbose               = false
	assembleSystemStartTime time.Time
	solveSystemStartTime    time.Time
)

/*
SetVerbosity sets the verbosity flag value.

The log functions will only output to the standard output if this flag's value is true.
*/
func SetVerbosity(verbose bool) {
	fmt.Println("Setting verbosity")
	isVerbose = verbose
}

/*
StartSolve should be called when the solving process starts.
*/
func StartSolve() {
	if isVerbose {
		fmt.Printf("----- [ inkfem ] -----\n")
	}
}

/*
StartAssembleSysEqs should be called when the structure's system of equations is about
to be assembled.
*/
func StartAssembleSysEqs() {
	if isVerbose {
		fmt.Println("> assembling system of equations...")
		assembleSystemStartTime = time.Now()
	}
}

/*
EndAssembleSysEqs should be called when the structure's system of equations has been
completely assembled.
*/
func EndAssembleSysEqs(sysSize int) {
	if isVerbose {
		elapsedTime := time.Since(assembleSystemStartTime)
		fmt.Printf("[DONE][%s] assembled system with %d equations\n", elapsedTime, sysSize)
	}
}

/*
StartSolveSysEqs should be called when the structure's system of equations is about
to be solved.
*/
func StartSolveSysEqs() {
	if isVerbose {
		fmt.Println("> solving sytem of equations for global displacements")
		solveSystemStartTime = time.Now()
	}
}

/*
EndSolveSysEqs should be called when the structure's system of equations has been solved.
*/
func EndSolveSysEqs(iterations int, minError float64) {
	if isVerbose {
		elapsedTime := time.Since(solveSystemStartTime)
		fmt.Printf(
			"[DONE][%s] solved system in %d iterations, error = %f\n",
			elapsedTime,
			iterations,
			minError,
		)
	}
}
