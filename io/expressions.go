/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package io

import "fmt"

const (
	floatExpr     = `-?\d+\.?\d*`
	validNameExpr = `[\w\-_ ]+`
	validIDExpr   = `\d+`

	nameGrpExpr             = `'(?P<name>` + validNameExpr + `)'`
	idGrpExpr               = `(?P<id>` + validIDExpr + `)`
	arrowExpr               = `\s*->\s*`
	loadTermExpr            = `(?P<term>[fm]{1}[xyz]{1})\s+`
	loadElementID           = `(?P<element>` + validIDExpr + `)\s+`
	distributedLoadRefExpr  = `(?P<ref>[lg]{1})(?:d{1})\s+`
	concentratedLoadRefExpr = `(?P<ref>[lg]{1})(?:c{1})\s+`
)

func floatGroupAndSpaceExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)\s+`, groupName, floatExpr)
}

func floatGroupAndOptinalSpaceExpr(groupName string) string {
	return fmt.Sprintf(`(?P<%s>%s)\s*`, groupName, floatExpr)
}
