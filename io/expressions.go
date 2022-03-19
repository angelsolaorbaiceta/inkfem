package io

import "fmt"

const (
	floatExpr         = `[+-]?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?`
	validNameExpr     = `[\w\-_ ]+`
	validIDExpr       = `[\w\-_]+`
	constraintExpr    = `{[drxyz ]*}`
	optionalSpaceExpr = `\s*`
	spaceExpr         = `\s+`
)

const (
	NameGrpExpr             = `'(?P<name>` + validNameExpr + `)'`
	IdGrpExpr               = `(?P<id>` + validIDExpr + `)`
	ArrowExpr               = `\s*->\s*`
	LoadTermExpr            = `(?P<term>[fm]{1}[xyz]{1})\s+`
	LoadElementID           = `(?P<element>` + validIDExpr + `)\s+`
	DistributedLoadRefExpr  = `(?P<ref>[lg]{1})d\s+`
	ConcentratedLoadRefExpr = `(?P<ref>[lg]{1})c\s+`
	DofGroupName            = "dof"
	DofGroup                = `(?:\| \[(?P<` + DofGroupName + `>\d+ \d+ \d+)\])?`
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
