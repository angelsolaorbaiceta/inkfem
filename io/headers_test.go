package io

import "testing"

func TestSectionHeaders(t *testing.T) {
	t.Run("is a section header line", func(t *testing.T) {
		line := "|nodes|"

		if !IsSectionHeaderLine(line) {
			t.Errorf("expected %q to be a section header line", line)
		}
	})

	// To guarantee compatibiilty with version v1.0, we need to support the old section format
	// that includes the count of elements below it.
	t.Run("is a deprecated header line with count", func(t *testing.T) {
		lines := []string{"|nodes| 5", "|nodes|5"}

		for _, line := range lines {
			if !IsSectionHeaderLine(line) {
				t.Errorf("expected %q to be a section header line", line)
			}
		}
	})

	t.Run("parses a section header", func(t *testing.T) {
		line := "|nodes|"

		if section := ParseSectionHeader(line); section != "nodes" {
			t.Errorf("expected %q to parse to %q", line, section)
		}
	})
}
