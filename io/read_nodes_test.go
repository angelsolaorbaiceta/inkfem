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
			"4 -> 20.5 50.2 { } | [3 4 50]",
		}
		nodes = deserializeNodesByID(lines)

		nodeOne   = structure.MakeNodeAtPosition("1", 10.1, 20.2, &structure.FullConstraint)
		nodeTwo   = structure.MakeNodeAtPosition("2", 40.1, 50.2, &structure.DispConstraint)
		nodeThree = structure.MakeNodeAtPosition("3", 70.1, 80.2, &structure.NilConstraint)
		nodeFour  = structure.MakeNodeAtPosition("4", 20.5, 50.2, &structure.NilConstraint).SetDegreesOfFreedomNum(3, 4, 50)
	)

	if got := nodes["1"]; !got.Equals(nodeOne) {
		t.Errorf("Want node %v, got %v", nodeOne, got)
	}
	if got := nodes["2"]; !got.Equals(nodeTwo) {
		t.Errorf("Want node %v, got %v", nodeTwo, got)
	}
	if got := nodes["3"]; !got.Equals(nodeThree) {
		t.Errorf("Want node %v, got %v", nodeThree, got)
	}
	if got := nodes["4"]; !got.Equals(nodeFour) {
		t.Errorf("Want node %v, got %v", nodeFour, got)
	}
}
