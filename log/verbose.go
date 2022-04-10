package log

import (
	"fmt"
	"log"
	"time"

	"github.com/angelsolaorbaiceta/inkmath/lineq"
)

var (
	isVerbose = false

	readFileStartTime   time.Time
	readFileElapsedTime time.Duration

	preprocessStartTime   time.Time
	preprocessElapsedTime time.Duration

	assembleSystemStartTime   time.Time
	assembleSystemElapsedTime time.Duration

	solveSystemStartTime   time.Time
	solveSystemElapsedTime time.Duration

	computeStressesStartTime time.Time
	computeStressesEndTime   time.Duration
)

// SetVerbosity sets the verbosity flag value.
// The log functions will only output to the standard output if this flag's value is true.
func SetVerbosity(verbose bool) {
	isVerbose = verbose
}

// StartProcess should be called when the solving process starts.
func StartProcess() {
	if isVerbose {
		// TODO: read version from file
		log.Printf("---------- [ inkfem v1.0 ] ----------\n")
	}
}

// StartReadFile should be called when the process of reading and parsing the input structure
// file is about to start.
func StartReadFile() {
	if isVerbose {
		readFileStartTime = time.Now()
	}
}

// EndReadFile should be called when the process of reading and parsing the input structure
// has been completed successfully.
func EndReadFile(fileType string, nodesCount, elementsCount int) {
	if isVerbose {
		readFileElapsedTime = time.Since(readFileStartTime)
		message := fmt.Sprintf(
			"read '%s' file (%d nodes and %d bars)", fileType, nodesCount, elementsCount,
		)
		writeDone(message, readFileElapsedTime)
	}
}

// StartPreprocess should be called when the preprocessing of the structure is about to start.
func StartPreprocess() {
	if isVerbose {
		preprocessStartTime = time.Now()
	}
}

// EndPreprocess should be called when the structure has been successfully preprocessed.
func EndPreprocess() {
	if isVerbose {
		preprocessElapsedTime = time.Since(preprocessStartTime)
		writeDone("stucture preprocessed", preprocessElapsedTime)
	}
}

// StartAssembleSysEqs should be called when the structure's system of equations is about
// to be assembled.
func StartAssembleSysEqs() {
	if isVerbose {
		assembleSystemStartTime = time.Now()
	}
}

// EndAssembleSysEqs should be called when the structure's system of equations has been
// completely assembled.
func EndAssembleSysEqs(sysSize int) {
	if isVerbose {
		assembleSystemElapsedTime = time.Since(assembleSystemStartTime)
		message := fmt.Sprintf("assembled system of %d equations", sysSize)
		writeDone(message, assembleSystemElapsedTime)
	}
}

// StartSolveSysEqs should be called when the structure's system of equations is about
// to be solved.
func StartSolveSysEqs() {
	if isVerbose {
		solveSystemStartTime = time.Now()
	}
}

var lastProgressPercentage = -1

// SolveSysProgress should be called at each iteration of the solving process.
func SolveSysProgress(progress lineq.IterativeSolverProgress) {
	if isVerbose && progress.ProgressPercentage%5 == 0 && progress.ProgressPercentage > lastProgressPercentage {
		lastProgressPercentage = progress.ProgressPercentage
		log.Printf(
			"[solver] %3d%%, %d iterations, error ~ %f\n",
			progress.ProgressPercentage, progress.IterCount, progress.Error,
		)
	}
}

// EndSolveSysEqs should be called when the structure's system of equations has been solved.
func EndSolveSysEqs(iterations int, minError float64) {
	if isVerbose {
		solveSystemElapsedTime = time.Since(solveSystemStartTime)
		message := fmt.Sprintf(
			"solved system of equations in %d iterations, error = %f", iterations, minError,
		)
		writeDone(message, solveSystemElapsedTime)
	}
}

// StartComputeStresses should be called when the sliced elements stresses are about to
// start being computed.
func StartComputeStresses() {
	if isVerbose {
		computeStressesStartTime = time.Now()
	}
}

// EndComputeStresses should be called when the stresses on all elements have been computed.
func EndComputeStresses() {
	if isVerbose {
		computeStressesEndTime = time.Since(computeStressesStartTime)
		writeDone("computed stresses for all elements", computeStressesEndTime)
	}
}

// Result should be called at the end of the execution to display the overall
// execution time results.
func Result() {
	if isVerbose {
		totalTime := readFileElapsedTime +
			preprocessElapsedTime +
			assembleSystemElapsedTime +
			solveSystemElapsedTime +
			computeStressesEndTime

		log.Printf("Total time: %s\n", totalTime)
	}
}

func writeDone(message string, elapsedTime time.Duration) {
	log.Printf("%s (took %s)\n", message, elapsedTime)
}
