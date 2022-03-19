package io

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// A LinesReader reads lines from a buffered scanner one by one, ignoring blank and commented
// lines.
type LinesReader struct {
	scanner        *bufio.Scanner
	nextLine       *string
	nextLineNumber int
}

// MakeLinesReader creates a lines using the passed in reader.
func MakeLinesReader(reader io.Reader) *LinesReader {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	return &LinesReader{scanner: scanner, nextLineNumber: 1}
}

// ReadNext checks if there are more lines available and reads the next one.
// If HasMoreLines returns "true", GetNextLine and GetNextLineNumber can be called to get the
// next line and its original line number.
func (lr *LinesReader) ReadNext() bool {
	lr.nextLine = nil
	var line string

	for lr.scanner.Scan() {
		line = strings.TrimSpace(lr.scanner.Text())
		lr.nextLineNumber += 1

		if ShouldIgnoreLine(line) {
			continue
		}

		lr.nextLine = &line

		return true
	}

	return false
}

// GetNextLine returns the next line read by the ReadNext method.
func (lr *LinesReader) GetNextLine() string {
	lr.ensureHasNextLine()
	return *lr.nextLine
}

// GetNextLineNumber returns the next line number read by the ReadNext method.
func (lr *LinesReader) GetNextLineNumber() int {
	lr.ensureHasNextLine()
	return lr.nextLineNumber
}

// GetNextLines gets the next "count" lines from the reader, ignoring comments and blank lines.
// Panics if there're not enough lines left in the reader.
func (lr *LinesReader) GetNextLines(count int) []string {
	var (
		line  string
		lines = make([]string, count)
	)

	for i := 0; i < count; i++ {
		if !lr.ReadNext() {
			panic(fmt.Sprintf("Couldn't read all expected %d lines", count))
		}

		line = lr.GetNextLine()
		lines[i] = line
	}

	return lines
}

func (lr *LinesReader) ensureHasNextLine() {
	if lr.nextLine == nil {
		panic("Can't read more lines from LinesReader")
	}
}
