// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quota

import (
	"encoding/json"
)

const (
	// ResourceCount count, in number
	ResourceCount ResourceName = "count"
	// ResourceStorage storage size, in bytes
	ResourceStorage ResourceName = "storage"
)

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// ResourceList is a set of (resource name, value) pairs.
type ResourceList map[ResourceName]int64

func (resources ResourceList) String() string {
	bytes, _ := json.Marshal(resources)
	return string(bytes)
}

// NewResourceList returns resource list from string
func NewResourceList(s string) (ResourceList, error) {
	var resources ResourceList
	if err := json.Unmarshal([]byte(s), &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

// Equals returns true if the two lists are equivalent
func Equals(a ResourceList, b ResourceList) bool {
	if len(a) != len(b) {
		return false
	}

	for key, value1 := range a {
		value2, found := b[key]
		if !found {
			return false
		}
		if value1 != value2 {
			return false
		}
	}

	return true
}

// Add returns the result of a + b for each named resource
func Add(a ResourceList, b ResourceList) ResourceList {
	result := ResourceList{}
	for key, value := range a {
		if other, found := b[key]; found {
			value = value + other
		}
		result[key] = value
	}

	for key, value := range b {
		if _, found := result[key]; !found {
			result[key] = value
		}
	}
	return result
}

// Subtract returns the result of a - b for each named resource
func Subtract(a ResourceList, b ResourceList) ResourceList {
	result := ResourceList{}
	for key, value := range a {
		if other, found := b[key]; found {
			value = value - other
		}
		result[key] = value
	}

	for key, value := range b {
		if _, found := result[key]; !found {
			result[key] = -value
		}
	}

	return result
}
