package io

import "fmt"

const (
	floatExpr         = `[+-]?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?`
	validNameExpr     = `[\w\-_ ]+`
	validIDExpr       = `[\w\-_]+`
	constraintExpr    = `{[drxyz ]*}`
	optionalSpaceExpr = `\s*`
	spaceExpr         = `\s+`

	nameGrpExpr             = `'(?P<name>` + validNameExpr + `)'`
	idGrpExpr               = `(?P<id>` + validIDExpr + `)`
	arrowExpr               = `\s*->\s*`
	loadTermExpr            = `(?P<term>[fm]{1}[xyz]{1})\s+`
	loadElementID           = `(?P<element>` + validIDExpr + `)\s+`
	distributedLoadRefExpr  = `(?P<ref>[lg]{1})d\s+`
	concentratedLoadRefExpr = `(?P<ref>[lg]{1})c\s+`
)

func floatGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)`, groupName, floatExpr)
}

func idGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)`, groupName, validIDExpr)
}

func nameGroupExpr(groupName string) string {
	return fmt.Sprintf(`'(?P<%s>%s)'`, groupName, validNameExpr)
}

func constraintGroupExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)`, groupName, constraintExpr)
}
