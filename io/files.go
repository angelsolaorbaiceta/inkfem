package io

import "os"

// CreateFile returns a new open file to be writen to. Panics if the file can't be created.
// Don't forget to close the file.
func CreateFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could't create file: " + err.Error())
	}

	return file
}

// OpenFile returns an existing file to be writen to. Panics if the file can't be opened.
// Don't forget to close the file.
func OpenFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Couldn't open file: " + err.Error())
	}

	return file
}
