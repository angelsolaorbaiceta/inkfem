package sol

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
)

func TestWriteSolution(t *testing.T) {
	var (
		sol    = inkio.MakeTestSolution()
		writer bytes.Buffer
	)

	Write(sol, &writer)
	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

	fmt.Println(gotLines)

	t.Run("first line is always the header with the version", func(t *testing.T) {
		var (
			want = fmt.Sprintf("inkfem v%d.%d", sol.Metadata.MajorVersion, sol.Metadata.MinorVersion)
			got  = gotLines[0]
		)

		if got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})
}
