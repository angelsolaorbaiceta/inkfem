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
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	elementWithLoadsSlices    = 10
	elementWithoutLoadsSlices = 6
)

/*
DoElement preprocesses the given structural element subdividing it as corresponds.
The result is sent through a channel.
*/
func DoElement(element *structure.Element, c chan<- *Element) {
	if element.IsAxialMember() {
		c <- sliceAxialElement(element)
	} else if element.HasLoadsApplied() {
		c <- sliceLoadedElement(element, elementWithLoadsSlices)
	} else {
		c <- sliceElementWithoutLoads(element, elementWithoutLoadsSlices)
	}
}
