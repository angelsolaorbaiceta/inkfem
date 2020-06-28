package log

import (
	"fmt"
	"time"
)

var (
	isVerbose                = false
	readFileStartTime        time.Time
	preprocessStartTime      time.Time
	assembleSystemStartTime  time.Time
	solveSystemStartTime     time.Time
	computeStressesStartTime time.Time
)

/*
SetVerbosity sets the verbosity flag value.

The log functions will only output to the standard output if this flag's value is true.
*/
func SetVerbosity(verbose bool) {
	isVerbose = verbose
}

/*
StartProcess should be called when the solving process starts.
*/
func StartProcess() {
	if isVerbose {
		// TODO: read version from file
		fmt.Printf("----- [ inkfem v1.0 ] -----\n")
	}
}

/*
StartReadFile should be called when the process of reading and parsing the input structure
file is about to start.
*/
func StartReadFile() {
	if isVerbose {
		writeStart("reading structure from file")
		readFileStartTime = time.Now()
	}
}

/*
EndReadFile should be called when the process of reading and parsing the input structure
has been completed successfully.
*/
func EndReadFile(nodesCount, elementsCount int) {
	if isVerbose {
		elapsedTime := time.Since(readFileStartTime)
		message := fmt.Sprintf(
			"reading input file (%d nodes and %d elements)", nodesCount, elementsCount,
		)
		writeDone(message, elapsedTime)
	}
}

/*
StartPreprocess should be called when the preprocessing of the structure is about
to start.
*/
func StartPreprocess() {
	if isVerbose {
		writeStart("preprocessing structure")
		preprocessStartTime = time.Now()
	}
}

/*
EndPreprocess should be called when the structure has been successfully preprocessed.
*/
func EndPreprocess() {
	if isVerbose {
		elapsedTime := time.Since(preprocessStartTime)
		writeDone("stucture preprocessed", elapsedTime)
	}
}

/*
StartAssembleSysEqs should be called when the structure's system of equations is about
to be assembled.
*/
func StartAssembleSysEqs() {
	if isVerbose {
		writeStart("assembling system of equations")
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
		message := fmt.Sprintf("assembled system of %d equations", sysSize)
		writeDone(message, elapsedTime)
	}
}

/*
StartSolveSysEqs should be called when the structure's system of equations is about
to be solved.
*/
func StartSolveSysEqs() {
	if isVerbose {
		writeStart("solving sytem of equations for global displacements")
		solveSystemStartTime = time.Now()
	}
}

/*
EndSolveSysEqs should be called when the structure's system of equations has been solved.
*/
func EndSolveSysEqs(iterations int, minError float64) {
	if isVerbose {
		elapsedTime := time.Since(solveSystemStartTime)
		message := fmt.Sprintf(
			"solved system in %d iterations, error = %f", iterations, minError,
		)
		writeDone(message, elapsedTime)
	}
}

/*
StartComputeStresses should be called when the sliced elements stresses are about to
start being computed.
*/
func StartComputeStresses() {
	if isVerbose {
		writeStart("solving element stresses")
		computeStressesStartTime = time.Now()
	}
}

/*
EndComputeStresses should be called when the stresses on all elements have been computed.
*/
func EndComputeStresses() {
	if isVerbose {
		elapsedTime := time.Since(computeStressesStartTime)
		writeDone("computed stresses for all elements", elapsedTime)
	}
}

func writeStart(message string) {
	fmt.Println()
	fmt.Println("> " + message + "...")
}

func writeDone(message string, elapsedTime time.Duration) {
	fmt.Printf("[DONE in %s] "+message+"\n", elapsedTime)
}
