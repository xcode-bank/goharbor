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

package label

import "github.com/goharbor/harbor/src/pkg/retention/res"

const (
	With    = "with labels"
	Without = "without labels"
)

// selector is for label selector
type selector struct {
	// Pre defined pattern decorations
	// "with" or "without"
	decoration string
	// Label list
	labels []string
}

// Select candidates by regular expressions
func (s *selector) Select(artifacts []*res.Candidate) ([]*res.Candidate, error) {
	return nil, nil
}

// New is factory method for list selector
func New(decoration string, pattern interface{}) res.Selector {
	return &selector{
		decoration: decoration,
		labels:     pattern.([]string),
	}
}
