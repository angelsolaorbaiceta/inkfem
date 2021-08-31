package io

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestDeserializeNode(t *testing.T) {
	t.Run("deserializes the node", func(t *testing.T) {
		var (
			got  = deserializeNode("1 -> 10.1 20.2 { dx dy rz }")
			want = structure.MakeNode("1", g2d.MakePoint(10.1, 20.2), &structure.FullConstraint)
		)

		if !got.Equals(want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("deserializes the node with scientific notation coordinates", func(t *testing.T) {
		var (
			got  = deserializeNode("1 -> 1e+2 2.0e-2 { dx dy rz }")
			want = structure.MakeNode("1", g2d.MakePoint(100.0, 0.02), &structure.FullConstraint)
		)

		if !got.Equals(want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestDeserializeNodes(t *testing.T) {
	var (
		lines []string = []string{
			"1 -> 10.1 20.2 { dx dy rz }",
			"2 -> 40.1 50.2 { dx dy }",
			"3 -> 70.1 80.2 { }",
		}
		nodes = deserializeNodesByID(lines)

		nodeOne   = structure.MakeNode("1", g2d.MakePoint(10.1, 20.2), &structure.FullConstraint)
		nodeTwo   = structure.MakeNode("2", g2d.MakePoint(40.1, 50.2), &structure.DispConstraint)
		nodeThree = structure.MakeNode("3", g2d.MakePoint(70.1, 80.2), &structure.NilConstraint)
	)

	if size := len(*nodes); size != 3 {
		t.Errorf("Expected 3 nodes, but got %d", size)
	}
	if got := (*nodes)["1"]; !got.Equals(nodeOne) {
		t.Errorf("Expected node %v, but got %v", nodeOne, got)
	}
	if got := (*nodes)["2"]; !got.Equals(nodeTwo) {
		t.Errorf("Expected node %v, but got %v", nodeTwo, got)
	}
	if got := (*nodes)["3"]; !got.Equals(nodeThree) {
		t.Errorf("Expected node %v, but got %v", nodeThree, got)
	}
}
