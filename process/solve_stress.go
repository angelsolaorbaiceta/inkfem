package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkgeom"
)

func computeStresses(sol *ElementSolution) {
	var (
		trailNode, leadNode                    preprocess.Node
		youngMod                               = sol.Element.Material().YoungMod
		iStrong                                = sol.Element.Section().IStrong
		nIndex, vIndex, mIndex                 = 0, 0, 0
		incX, trailDy, leadDy, trailRz, leadRz float64
		length                                 float64
	)

	for i := 1; i < len(sol.Element.Nodes); i++ {
		trailNode, leadNode = sol.Element.Nodes[i-1], sol.Element.Nodes[i]
		length = sol.Element.Geometry().LengthBetween(trailNode.T, leadNode.T)
		incX = sol.LocalXDispl[i].Value - sol.LocalXDispl[i-1].Value
		trailDy = sol.LocalYDispl[i-1].Value
		leadDy = sol.LocalYDispl[i].Value
		trailRz = sol.LocalZRot[i-1].Value
		leadRz = sol.LocalZRot[i].Value

		/* Axial */
		n := incX * youngMod / length
		sol.AxialStress[nIndex] = PointSolutionValue{
			trailNode.T,
			n - trailNode.LocalFx(),
		}
		sol.AxialStress[nIndex+1] = PointSolutionValue{
			leadNode.T,
			n + leadNode.LocalFx(),
		}
		nIndex += 2

		/* Shear */
		v := (6.0 * youngMod * iStrong / (length * length * length)) * ((2.0 * (trailDy - leadDy)) + (length * (leadRz - trailRz)))
		sol.ShearStress[vIndex] = PointSolutionValue{trailNode.T, v - trailNode.LocalFy()}
		sol.ShearStress[vIndex+1] = PointSolutionValue{leadNode.T, v + leadNode.LocalFy()}
		vIndex += 2

		/* Bending */
		eil2 := youngMod * iStrong / (length * length)
		sol.BendingMoment[mIndex] =
			PointSolutionValue{
				trailNode.T,
				eil2*(-6.0*trailDy+2.0*length*trailRz-6.0*leadDy+4.0*length*leadRz) + trailNode.LocalMz(),
			}
		sol.BendingMoment[mIndex+1] =
			PointSolutionValue{
				inkgeom.AverageT(trailNode.T, leadNode.T),
				(youngMod * iStrong / length) * (leadRz - trailRz),
			}
		sol.BendingMoment[mIndex+2] =
			PointSolutionValue{
				leadNode.T,
				eil2*(-6.0*trailDy+4.0*length*trailRz+6.0*leadDy-2.0*length*leadRz) + leadNode.LocalMz(),
			}
		mIndex += 3
	}
}
