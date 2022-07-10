# Preprocess Package

This package defines the _preprocessed_ or _sliced_ structure model which is used for the Finite Element Analysis.
It also provides the means for slicing or preprocessing the structure as it is defined in the _structure_ package.

A _sliced structure_ is one whose bars have been chopped into smaller finite elements where the loads originally applied in the bar are distributed in their nodes.
The loads are distributed in a way such that the work done by the equivalent nodal forces equals the work done by the load.

## Slicing the structure model

A structure model (`structure.Structure`) can be sliced using the `preprocess.StructureModel` function.
This function is the only public function this package exports.
The result is a `preprocess.Structure`.

Each bar is sliced differently depending on whether it's subject to only axial stress and the loads applied to it.
When a bar is an _axial member_, it isn't sliced at all.
A bar is an _axial member_ if the following three conditions are met: 

1. it's pinned in both ends, 
2. has no distributed loads applied to it, and
3. has forces (not moments) applied only on it's nodes.

Non-axial bars are sliced according to whether they have loads applied to them or not.
If they haven't got any load applied, they are sliced into a small number of finite elements.
If they have loads applied to them, they are sliced into a slightly larger number of elements, plus the positions where a concentrated load is applied and the start and end positions of distributed loads.
These locations in a bar's directrix are important to consider, as it's where stress discontinuities take place.