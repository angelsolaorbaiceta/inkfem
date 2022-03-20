package structure

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
)

type MaterialsByName = map[string]*Material
type SectionsByName = map[string]*Section
type NodesByIdMap = map[contracts.StrID]*Node
type ConcLoadsById = map[contracts.StrID][]*load.ConcentratedLoad
type DistLoadsById = map[contracts.StrID][]*load.DistributedLoad

type StructureData struct {
	Nodes             NodesByIdMap
	Materials         MaterialsByName
	Sections          SectionsByName
	ConcentratedLoads ConcLoadsById
	DistributedLoads  DistLoadsById
}
