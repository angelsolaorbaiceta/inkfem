package main

import (
	_ "embed"

	"github.com/angelsolaorbaiceta/inkfem/cmd"
	inkio "github.com/angelsolaorbaiceta/inkfem/io"
)

//go:embed VERSION
var versionString string

func main() {
	inkio.SetBinaryVersion(versionString)
	cmd.Execute()
}
