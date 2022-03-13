package io

import "os"

// CreateFile returns a new open file to be writen to. Panics if the file can't be created.
func CreateFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could't create file: " + err.Error())
	}

	return file
}
