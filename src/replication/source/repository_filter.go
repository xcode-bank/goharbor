// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
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

package source

import (
	"strings"

	"github.com/vmware/harbor/src/common/utils"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/replication"
	"github.com/vmware/harbor/src/replication/models"
	"github.com/vmware/harbor/src/replication/registry"
)

// RepositoryFilter implement Filter interface to filter repository
type RepositoryFilter struct {
	pattern   string
	convertor Convertor
}

// NewRepositoryFilter returns an instance of RepositoryFilter
func NewRepositoryFilter(pattern string, registry registry.Adaptor) *RepositoryFilter {
	return &RepositoryFilter{
		pattern:   pattern,
		convertor: NewRepositoryConvertor(registry),
	}
}

// Init ...
func (r *RepositoryFilter) Init() error {
	return nil
}

// GetConvertor ...
func (r *RepositoryFilter) GetConvertor() Convertor {
	return r.convertor
}

// DoFilter filters repository and image(according to the repository part) and drops any other resource types
func (r *RepositoryFilter) DoFilter(items []models.FilterItem) []models.FilterItem {
	candidates := []string{}
	for _, item := range items {
		candidates = append(candidates, item.Value)
	}
	log.Debugf("repository filter candidates: %v", candidates)

	result := []models.FilterItem{}
	for _, item := range items {
		if item.Kind != replication.FilterItemKindRepository && item.Kind != replication.FilterItemKindTag {
			log.Warningf("unsupported type %s for repository filter, drop", item.Kind)
			continue
		}

		repository := item.Value
		if item.Kind == replication.FilterItemKindTag {
			repository = strings.SplitN(repository, ":", 2)[0]
		}

		if len(r.pattern) == 0 {
			log.Debugf("pattern is null, add %s to the repository filter result list", item.Value)
			result = append(result, item)
		} else {
			// trim the project
			_, repository = utils.ParseRepository(repository)
			matched, err := match(r.pattern, repository)
			if err != nil {
				log.Errorf("failed to match pattern %s to value %s: %v", r.pattern, repository, err)
				break
			}
			if matched {
				log.Debugf("pattern %s matched, add %s to the repository filter result list", r.pattern, item.Value)
				result = append(result, item)
			}
		}
	}
	return result
}
