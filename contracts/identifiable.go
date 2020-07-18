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

package contracts

import "sort"

// StrID is the type used for structural data ids
type StrID = int

/*
Identifiable is anything that can be referenced using an integer number.
*/
type Identifiable interface {
	GetID() StrID
}

// ByID implements the sort.Interface for []Identifiable based in their id.
type ByID []Identifiable

func (a ByID) Len() int {
	return len(a)
}

func (a ByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByID) Less(i, j int) bool {
	return a[i].GetID() < a[j].GetID()
}

// SortByID sorts (in place) a slice of identifiable elements.
func SortByID(elements []Identifiable) {
	sort.Sort(ByID(elements))
}
