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

package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

/*
Structure result of preprocessing original structure, ready to be solved.
The elements of a preprocessed structure are already sliced.
*/
type Structure struct {
	Metadata  structure.StrMetadata
	Nodes     map[contracts.StrID]*structure.Node
	Elements  []*Element
	DofsCount int
}

/*
NodesCount returns the number of nodes in the original structure.
*/
func (s *Structure) NodesCount() int {
	return len(s.Nodes)
}

/*
ElementsCount returns the number of elements in the original structure.
*/
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}
