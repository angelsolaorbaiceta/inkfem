package preprocess

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestPreprocessStructureWithoutLoads(t *testing.T) {
	build.ReadBuildInfo()

	var (
		nodeOne  = structure.MakeNode("n1", g2d.MakePoint(0, 0), &structure.FullConstraint)
		nodeTwo  = structure.MakeNode("n2", g2d.MakePoint(300, 0), &structure.NilConstraint)
		material = structure.MakeUnitMaterial()
		section  = structure.MakeUnitSection()
		bar      = structure.MakeElementBuilder("b1").
				WithStartNode(nodeOne, &structure.FullConstraint).
				WithEndNode(nodeTwo, &structure.NilConstraint).
				WithMaterial(material).
				WithSection(section).
				Build()
		meta = structure.StrMetadata{MajorVersion: 2, MinorVersion: 3}
		str  = structure.Make(meta, map[contracts.StrID]*structure.Node{
			"n1": nodeOne,
			"n2": nodeTwo,
		}, []*structure.Element{bar})
	)

	t.Run("without own weight", func(t *testing.T) {
		var (
			options = PreprocessOptions{IncludeOwnWeight: false}
			result  = StructureModel(str, &options)
		)

		t.Run("has the right number of DOFs", func(t *testing.T) {
			// There are 3 DOFs per node in the sliced bar, which has elementWithoutLoadsSlices
			// plus 1 nodes, plus 3 DOFs for the free end node.
			wantDofs := 3*(elementWithoutLoadsSlices+1) + 3
			assert.Equal(t, wantDofs, result.DofsCount())
		})

		t.Run("the sliced element nodes are equally spaced", func(t *testing.T) {
			var (
				el   = result.Elements()[0]
				tPos = 0.0
				tInc = 1.0 / elementWithoutLoadsSlices
			)

			for _, node := range el.Nodes() {
				got := node.T.Value()
				assert.True(t, nums.FloatsEqual(got, tPos), "got %f, want %f", got, tPos)

				tPos += tInc
			}
		})

		t.Run("the sliced element nodes have no loads", func(t *testing.T) {
			el := result.Elements()[0]
			assert.Nil(t, el.ConcentratedLoads)
			assert.Nil(t, el.DistributedLoads)

			for _, node := range el.Nodes() {
				assert.Equal(t, math.NilTorsor, node.NetLocalLoadTorsor())
			}
		})
	})

	t.Run("with own weight", func(t *testing.T) {
		var (
			options = PreprocessOptions{IncludeOwnWeight: true}
			result  = StructureModel(str, &options)
		)

		t.Run("has the right number of DOFs", func(t *testing.T) {
			// There are 3 DOFs per node in the sliced bar, which has elementWithLoadsSlices
			// plus 1 nodes plus 3 DOFs for the free end node.
			wantDofs := 3*(elementWithLoadsSlices+1) + 3
			assert.Equal(t, wantDofs, result.DofsCount())
		})

		t.Run("the sliced element nodes are equally spaced", func(t *testing.T) {
			var (
				el   = result.Elements()[0]
				tPos = 0.0
				tInc = 1.0 / elementWithLoadsSlices
			)

			for _, node := range el.Nodes() {
				got := node.T.Value()
				assert.True(t, nums.FloatsEqual(got, tPos), "got %f, want %f", got, tPos)

				tPos += tInc
			}
		})

		t.Run("the sliced element nodes have the weight of the bar", func(t *testing.T) {
			el := result.Elements()[0]
			assert.Nil(t, el.ConcentratedLoads)
			assert.NotNil(t, el.DistributedLoads)

			var (
				elLength      = el.Length() / float64(elementWithLoadsSlices)
				wantFyLoadVal = -material.Density * section.Area * elLength
				wantMzLoadVal = -elLength * elLength / 12.0
			)

			for _, node := range el.Nodes() {
				// Fx is always 0.0
				assert.True(
					t,
					nums.FloatsEqual(node.NetLocalFx(), 0.0),
					"(at t = %f) got %f, want 0.0", node.T.Value(), node.NetLocalFx(),
				)

				// The first and last nodes have half the weight of the rest
				wantFy := wantFyLoadVal
				if node.T.Value() == 0.0 || node.T.Value() == 1.0 {
					wantFy /= 2
				}
				assert.True(
					t,
					nums.FloatsEqual(node.NetLocalFy(), wantFy),
					"(at t = %f) got %f, want %f", node.T.Value(), node.NetLocalFy(), wantFy,
				)

				// Mz is 0.0 in the middle nodes, but has a value in the first and last nodes
				wantMz := 0.0
				if node.T.Value() == 0.0 {
					wantMz = wantMzLoadVal
				}
				if node.T.Value() == 1.0 {
					wantMz = -wantMzLoadVal
				}
				assert.True(
					t,
					nums.FloatsEqual(node.NetLocalMz(), wantMz),
					"(at t = %f) got %f, want %f", node.T.Value(), node.NetLocalMz(), wantMz,
				)
			}
		})
	})
}
