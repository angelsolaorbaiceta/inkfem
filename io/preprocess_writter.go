/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import (
	"bufio"
	"os"
	"text/template"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

const solTemplatePath = "io/templates/preprocess.template.txt"

// PreprocessedStructureToFile Writes the given preprocessed structure to a file.
func PreprocessedStructureToFile(structure *preprocess.Structure, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic("Could not create file for preprocessed structure")
	}
	defer file.Close()

	var (
		tmpl   = template.Must(template.ParseFiles(solTemplatePath))
		writer = bufio.NewWriter(file)
	)

	tmpl.Execute(writer, structure)
	writer.Flush()
}
