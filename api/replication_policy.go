/*
    Copyright (c) 2016 VMware, Inc. All Rights Reserved.
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

package api

import (
	"fmt"

	"net/http"
	"strconv"

	"github.com/vmware/harbor/dao"
	"github.com/vmware/harbor/models"
	"github.com/vmware/harbor/utils/log"
)

// RepPolicyAPI handles /api/replicationPolicies /api/replicationPolicies/:id/enablement
type RepPolicyAPI struct {
	BaseAPI
}

// Prepare validates whether the user has system admin role
func (pa *RepPolicyAPI) Prepare() {
	uid := pa.ValidateUser()
	var err error
	isAdmin, err := dao.IsAdminRole(uid)
	if err != nil {
		log.Errorf("Failed to Check if the user is admin, error: %v, uid: %d", err, uid)
	}
	if !isAdmin {
		pa.CustomAbort(http.StatusForbidden, "")
	}
}

// Get ...
func (pa *RepPolicyAPI) Get() {
	id := pa.GetIDFromURL()
	policy, err := dao.GetRepPolicy(id)
	if err != nil {
		log.Errorf("failed to get policy %d: %v", id, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if policy == nil {
		pa.CustomAbort(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	pa.Data["json"] = policy
	pa.ServeJSON()
}

// List filters policies by name and project_id, if name and project_id
// are nil, List returns all policies
func (pa *RepPolicyAPI) List() {
	name := pa.GetString("name")
	projectIDStr := pa.GetString("project_id")

	var projectID int64
	var err error

	if len(projectIDStr) != 0 {
		projectID, err = strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil || projectID <= 0 {
			pa.CustomAbort(http.StatusBadRequest, "invalid project ID")
		}
	}

	policies, err := dao.FilterRepPolicies(name, projectID)
	if err != nil {
		log.Errorf("failed to filter policies %s project ID %d: %v", name, projectID, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	pa.Data["json"] = policies
	pa.ServeJSON()
}

// Post creates a policy, and if it is enbled, the replication will be triggered right now.
func (pa *RepPolicyAPI) Post() {
	policy := &models.RepPolicy{}
	pa.DecodeJSONReqAndValidate(policy)

	po, err := dao.GetRepPolicyByName(policy.Name)
	if err != nil {
		log.Errorf("failed to get policy %s: %v", policy.Name, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if po != nil {
		pa.CustomAbort(http.StatusConflict, "name is already used")
	}

	project, err := dao.GetProjectByID(policy.ProjectID)
	if err != nil {
		log.Errorf("failed to get project %d: %v", policy.ProjectID, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if project == nil {
		pa.CustomAbort(http.StatusBadRequest, fmt.Sprintf("project %d does not exist", policy.ProjectID))
	}

	target, err := dao.GetRepTarget(policy.TargetID)
	if err != nil {
		log.Errorf("failed to get target %d: %v", policy.TargetID, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if target == nil {
		pa.CustomAbort(http.StatusBadRequest, fmt.Sprintf("target %d does not exist", policy.TargetID))
	}

	pid, err := dao.AddRepPolicy(*policy)
	if err != nil {
		log.Errorf("Failed to add policy to DB, error: %v", err)
		pa.RenderError(http.StatusInternalServerError, "Internal Error")
		return
	}

	if policy.Enabled == 1 {
		go func() {
			if err := TriggerReplication(pid, "", nil, models.RepOpTransfer); err != nil {
				log.Errorf("failed to trigger replication of %d: %v", pid, err)
			} else {
				log.Infof("replication of %d triggered", pid)
			}
		}()
	}

	pa.Redirect(http.StatusCreated, strconv.FormatInt(pid, 10))
}

// Put modifies name and description of policy
func (pa *RepPolicyAPI) Put() {
	id := pa.GetIDFromURL()
	originalPolicy, err := dao.GetRepPolicy(id)
	if err != nil {
		log.Errorf("failed to get policy %d: %v", id, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if originalPolicy == nil {
		pa.CustomAbort(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	policy := &models.RepPolicy{}
	pa.DecodeJSONReq(policy)
	policy.ProjectID = originalPolicy.ProjectID
	policy.TargetID = originalPolicy.TargetID
	pa.Validate(policy)

	if policy.Name != originalPolicy.Name {
		po, err := dao.GetRepPolicyByName(policy.Name)
		if err != nil {
			log.Errorf("failed to get policy %s: %v", policy.Name, err)
			pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		if po != nil {
			pa.CustomAbort(http.StatusConflict, "name is already used")
		}
	}

	policy.ID = id

	if err = dao.UpdateRepPolicy(policy); err != nil {
		log.Errorf("failed to update policy %d: %v", id, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if policy.Enabled == originalPolicy.Enabled {
		return
	}

	//enablement has been modified
	if policy.Enabled == 1 {
		go func() {
			if err := TriggerReplication(id, "", nil, models.RepOpTransfer); err != nil {
				log.Errorf("failed to trigger replication of %d: %v", id, err)
			} else {
				log.Infof("replication of %d triggered", id)
			}
		}()
	} else {
		go func() {
			if err := postReplicationAction(id, "stop"); err != nil {
				log.Errorf("failed to stop replication of %d: %v", id, err)
			} else {
				log.Infof("try to stop replication of %d", id)
			}
		}()
	}
}

type enablementReq struct {
	Enabled int `json:"enabled"`
}

// UpdateEnablement changes the enablement of the policy
func (pa *RepPolicyAPI) UpdateEnablement() {
	id := pa.GetIDFromURL()
	policy, err := dao.GetRepPolicy(id)
	if err != nil {
		log.Errorf("failed to get policy %d: %v", id, err)
		pa.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if policy == nil {
		pa.CustomAbort(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	e := enablementReq{}
	pa.DecodeJSONReq(&e)
	if e.Enabled != 0 && e.Enabled != 1 {
		pa.RenderError(http.StatusBadRequest, "invalid enabled value")
		return
	}

	if policy.Enabled == e.Enabled {
		return
	}

	if err := dao.UpdateRepPolicyEnablement(id, e.Enabled); err != nil {
		log.Errorf("Failed to update policy enablement in DB, error: %v", err)
		pa.RenderError(http.StatusInternalServerError, "Internal Error")
		return
	}

	if e.Enabled == 1 {
		go func() {
			if err := TriggerReplication(id, "", nil, models.RepOpTransfer); err != nil {
				log.Errorf("failed to trigger replication of %d: %v", id, err)
			} else {
				log.Infof("replication of %d triggered", id)
			}
		}()
	} else {
		go func() {
			if err := postReplicationAction(id, "stop"); err != nil {
				log.Errorf("failed to stop replication of %d: %v", id, err)
			} else {
				log.Infof("try to stop replication of %d", id)
			}
		}()
	}
}
