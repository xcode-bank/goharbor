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

package robot

import (
	"context"
	"sync"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/rbac"
	"github.com/goharbor/harbor/src/controller/project"
	"github.com/goharbor/harbor/src/pkg/permission/evaluator"
	"github.com/goharbor/harbor/src/pkg/permission/types"
	"github.com/goharbor/harbor/src/pkg/robot/model"
)

// SecurityContext implements security.Context interface based on database
type SecurityContext struct {
	robot     *model.Robot
	ctl       project.Controller
	policy    []*types.Policy
	evaluator evaluator.Evaluator
	once      sync.Once
}

// NewSecurityContext ...
func NewSecurityContext(robot *model.Robot, policy []*types.Policy) *SecurityContext {
	return &SecurityContext{
		ctl:    project.Ctl,
		robot:  robot,
		policy: policy,
	}
}

// Name returns the name of the security context
func (s *SecurityContext) Name() string {
	return "robot"
}

// IsAuthenticated returns true if the user has been authenticated
func (s *SecurityContext) IsAuthenticated() bool {
	return s.robot != nil
}

// GetUsername returns the username of the authenticated user
// It returns null if the user has not been authenticated
func (s *SecurityContext) GetUsername() string {
	if !s.IsAuthenticated() {
		return ""
	}
	return s.robot.Name
}

// IsSysAdmin robot cannot be a system admin
func (s *SecurityContext) IsSysAdmin() bool {
	return false
}

// IsSolutionUser robot cannot be a system admin
func (s *SecurityContext) IsSolutionUser() bool {
	return false
}

// Can returns whether the robot can do action on resource
func (s *SecurityContext) Can(ctx context.Context, action types.Action, resource types.Resource) bool {
	s.once.Do(func() {
		s.evaluator = rbac.NewProjectEvaluator(s.ctl, rbac.NewBuilderForPolicies(s.GetUsername(), s.policy, filterRobotPolicies))
	})

	return s.evaluator != nil && s.evaluator.HasPermission(ctx, resource, action)
}

func filterRobotPolicies(p *models.Project, policies []*types.Policy) []*types.Policy {
	namespace := rbac.NewProjectNamespace(p.ProjectID)

	var results []*types.Policy
	for _, policy := range policies {
		if types.ResourceAllowedInNamespace(policy.Resource, namespace) {
			results = append(results, policy)
			// give the PUSH action a pull access
			if policy.Action == rbac.ActionPush {
				results = append(results, &types.Policy{Resource: policy.Resource, Action: rbac.ActionPull})
			}
		}
	}
	return results
}
