# Preprocess Package

This package defines the _preprocessed_ or _sliced_ structure model which is used for the Finite Element Analysis.
It also provides the means for slicing or preprocessing the structure as it is defined in the _structure_ package.

## Slicing the structure model

A structure model (`structure.Structure`) can be sliced using the `preprocess.StructureModel` function.
This function is the only public function this package exports.
The result is a `preprocess.Structure`.
