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

package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/harbor/src/common/utils/test"
	"github.com/vmware/harbor/src/replication"
	"github.com/vmware/harbor/src/replication/models"
	"github.com/vmware/harbor/src/replication/source"
)

func TestMain(m *testing.M) {
	GlobalController = NewDefaultController(ControllerConfig{})
	// set the policy manager used by GlobalController with a fake policy manager
	controller := GlobalController.(*DefaultController)
	controller.policyManager = &test.FakePolicyManager{}
	os.Exit(m.Run())
}

func TestNewDefaultController(t *testing.T) {
	controller := NewDefaultController(ControllerConfig{})
	assert.NotNil(t, controller)
}

func TestInit(t *testing.T) {
	assert.Nil(t, GlobalController.Init())
}

func TestCreatePolicy(t *testing.T) {
	_, err := GlobalController.CreatePolicy(models.ReplicationPolicy{
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindManual,
		},
	})
	assert.Nil(t, err)
}

func TestUpdatePolicy(t *testing.T) {
	assert.Nil(t, GlobalController.UpdatePolicy(models.ReplicationPolicy{
		ID: 2,
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindManual,
		},
	}))
}

func TestRemovePolicy(t *testing.T) {
	assert.Nil(t, GlobalController.RemovePolicy(1))
}

func TestGetPolicy(t *testing.T) {
	_, err := GlobalController.GetPolicy(1)
	assert.Nil(t, err)
}

func TestGetPolicies(t *testing.T) {
	_, err := GlobalController.GetPolicies(models.QueryParameter{})
	assert.Nil(t, err)
}

func TestReplicate(t *testing.T) {
	// TODO
}

func TestGetCandidates(t *testing.T) {
	policy := &models.ReplicationPolicy{
		ID: 1,
		Filters: []models.Filter{
			models.Filter{
				Kind:    replication.FilterItemKindTag,
				Pattern: "*",
			},
		},
		Trigger: &models.Trigger{
			Kind: replication.TriggerKindImmediate,
		},
	}

	sourcer := source.NewSourcer()

	candidates := []models.FilterItem{
		models.FilterItem{
			Kind:  replication.FilterItemKindTag,
			Value: "library/hello-world:release-1.0",
		},
		models.FilterItem{
			Kind:  replication.FilterItemKindTag,
			Value: "library/hello-world:latest",
		},
	}
	metadata := map[string]interface{}{
		"candidates": candidates,
	}
	result := getCandidates(policy, sourcer, metadata)
	assert.Equal(t, 2, len(result))

	policy.Filters = []models.Filter{
		models.Filter{
			Kind:    replication.FilterItemKindTag,
			Pattern: "release-*",
		},
	}
	result = getCandidates(policy, sourcer, metadata)
	assert.Equal(t, 1, len(result))
}

func TestBuildFilterChain(t *testing.T) {
	policy := &models.ReplicationPolicy{
		ID: 1,
		Filters: []models.Filter{
			models.Filter{
				Kind:    replication.FilterItemKindRepository,
				Pattern: "*",
			},
			models.Filter{
				Kind:    replication.FilterItemKindTag,
				Pattern: "*",
			},
		},
	}

	sourcer := source.NewSourcer()

	chain := buildFilterChain(policy, sourcer)
	assert.Equal(t, 2, len(chain.Filters()))
}
