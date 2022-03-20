package io

import "strings"

// DefinitionFileExt is the extension of the structure definition files.
const DefinitionFileExt = ".inkfem"

// PreFileExt is the extension of the preprocessed structure files.
const PreFileExt = ".inkfempre"

// SolFileExt is the extension of the solved structure files.
const SolFileExt = ".inkfemsol"

// IsDefinitionFile returns true if the file extension in the path is .inkfem.
func IsDefinitionFile(path string) bool {
	return strings.HasSuffix(path, DefinitionFileExt)
}

// IsPreprocessedFile returns true if the file extension in the path is .inkfempre.
func IsPreprocessedFile(path string) bool {
	return strings.HasSuffix(path, PreFileExt)
}
