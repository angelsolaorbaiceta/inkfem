package structure

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

type ElementBuilder struct {
	id                 contracts.StrID
	startNode, endNode *Node
	startLink, endLink *Constraint
	material           *Material
	section            *Section
	concentratedLoads  []*load.ConcentratedLoad
	distributedLoads   []*load.DistributedLoad
}

func MakeElementBuilder(id contracts.StrID) *ElementBuilder {
	builder := ElementBuilder{}
	builder.id = id

	return &builder
}

func (builder *ElementBuilder) WithStartNode(
	startNode *Node,
	startLink *Constraint,
) *ElementBuilder {
	builder.startNode = startNode
	builder.startLink = startLink

	return builder
}

func (builder *ElementBuilder) WithEndNode(
	endNode *Node,
	endLink *Constraint,
) *ElementBuilder {
	builder.endNode = endNode
	builder.endLink = endLink

	return builder
}

func (builder *ElementBuilder) WithMaterial(material *Material) *ElementBuilder {
	builder.material = material
	return builder
}

func (builder *ElementBuilder) WithSection(section *Section) *ElementBuilder {
	builder.section = section
	return builder
}

func (builder *ElementBuilder) IncludeOwnWeightLoad() *ElementBuilder {
	builder.ensureSectionAndMaterial()

	var (
		loadValue = -builder.section.Area * builder.material.Density
		load      = load.MakeDistributed(load.FY, false, nums.MinT, loadValue, nums.MaxT, loadValue)
	)
	builder.distributedLoads = append(builder.distributedLoads, load)

	return builder
}

func (builder *ElementBuilder) AddConcentratedLoads(loads []*load.ConcentratedLoad) *ElementBuilder {
	builder.concentratedLoads = append(builder.concentratedLoads, loads...)
	return builder
}

func (builder *ElementBuilder) AddDistributedLoads(loads []*load.DistributedLoad) *ElementBuilder {
	builder.distributedLoads = append(builder.distributedLoads, loads...)
	return builder
}

func (builder *ElementBuilder) AddDistributedLoad(load *load.DistributedLoad) *ElementBuilder {
	builder.distributedLoads = append(builder.distributedLoads, load)
	return builder
}

func (builder ElementBuilder) Build() *Element {
	builder.ensureNodesInfo()
	builder.ensureSectionAndMaterial()

	return &Element{
		id:                builder.id,
		startNodeID:       builder.startNode.GetID(),
		endNodeID:         builder.endNode.GetID(),
		geometry:          g2d.MakeSegment(builder.startNode.Position, builder.endNode.Position),
		startLink:         builder.startLink,
		endLink:           builder.endLink,
		material:          builder.material,
		section:           builder.section,
		ConcentratedLoads: builder.concentratedLoads,
		DistributedLoads:  builder.distributedLoads,
	}
}

func (builder ElementBuilder) ensureNodesInfo() {
	if builder.startNode == nil || builder.startLink == nil {
		panic("The start node information isn't defined")
	}
	if builder.endNode == nil || builder.endLink == nil {
		panic("The end node information isn't defined")
	}
}

func (builder ElementBuilder) ensureSectionAndMaterial() {
	if builder.section == nil {
		panic("The section isn't defined")
	}
	if builder.material == nil {
		panic("The material isn't defined")
	}
}
