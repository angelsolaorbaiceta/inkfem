package io

import "fmt"

const (
	floatExpr      = `[+-]?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?`
	validNameExpr  = `[\w\-_ ]+`
	validIDExpr    = `[\w\-_]+`
	constraintExpr = `{[drxyz ]*}`
)

const (
	OptionalSpaceExpr       = `\s*`
	SpaceExpr               = `\s+`
	NameGrpName             = "name"
	NameGrpExpr             = `'(?P<` + NameGrpName + `>` + validNameExpr + `)'`
	IdGrpName               = "id"
	IdGrpExpr               = `(?P<` + IdGrpName + `>` + validIDExpr + `)`
	ArrowExpr               = `\s*->\s*`
	LoadTermExpr            = `(?P<term>[fm]{1}[xyz]{1})\s+`
	LoadElementID           = `(?P<element>` + validIDExpr + `)\s+`
	DistributedLoadRefExpr  = `(?P<ref>[lg]{1})d\s+`
	ConcentratedLoadRefExpr = `(?P<ref>[lg]{1})c\s+`
	DofGrpName              = "dof"
	DofGrpExpr              = `(?P<` + DofGrpName + `>\[\d+ \d+ \d+\])`
)

func FloatGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)`, groupName, floatExpr)
}

func IdGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)`, groupName, validIDExpr)
}

func NameGroupExpr(groupName string) string {
	return fmt.Sprintf(`'(?P<%s>%s)'`, groupName, validNameExpr)
}

func ConstraintGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)`, groupName, constraintExpr)
}

func TorsorGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>{%s %s %s})`, groupName, floatExpr, floatExpr, floatExpr)
}
