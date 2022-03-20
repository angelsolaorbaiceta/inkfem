package def

import (
	"bytes"
	"fmt"
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
)

func TestWriteDefinition(t *testing.T) {
	var (
		str    = inkio.MakeTestOriginalStructure()
		writer bytes.Buffer
	)

	Write(str, &writer)
	fmt.Println(writer.String())
	// var gotLines []string
	// for _, line := range strings.Split(writer.String(), "\n") {
	// 	if line != "" {
	// 		gotLines = append(gotLines, line)
	// 	}
	// }
}
