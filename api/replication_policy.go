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
	policyID int64
	policy   *models.RepPolicy
}

// Prepare validates whether the user has system admin role
// and parsed the policy ID if it exists
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
	idStr := pa.Ctx.Input.Param(":id")
	if len(idStr) > 0 {
		pa.policyID, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Errorf("Error parsing policy id: %s, error: %v", idStr, err)
			pa.CustomAbort(http.StatusBadRequest, "invalid policy id")
		}
		p, err := dao.GetRepPolicy(pa.policyID)
		if err != nil {
			log.Errorf("Error occurred in GetRepPolicy, error: %v", err)
			pa.CustomAbort(http.StatusInternalServerError, "Internal error.")
		}
		if p == nil {
			pa.CustomAbort(http.StatusNotFound, fmt.Sprintf("policy does not exist, id: %v", pa.policyID))
		}
		pa.policy = p
	}
}

// Get gets all the policies according to the project
func (pa *RepPolicyAPI) Get() {
	projectID, err := pa.GetInt64("project_id")
	if err != nil {
		log.Errorf("Failed to get project id, error: %v", err)
		pa.RenderError(http.StatusBadRequest, "Invalid project id")
		return
	}
	policies, err := dao.GetRepPolicyByProject(projectID)
	if err != nil {
		log.Errorf("Failed to query policies from db, error: %v", err)
		pa.RenderError(http.StatusInternalServerError, "Failed to query policies")
		return
	}
	pa.Data["json"] = policies
	pa.ServeJSON()
}

// Post creates a policy, and if it is enbled, the replication will be triggered right now.
func (pa *RepPolicyAPI) Post() {
	policy := models.RepPolicy{}
	pa.DecodeJSONReq(&policy)
	pid, err := dao.AddRepPolicy(policy)
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

type enablementReq struct {
	Enabled int `json:"enabled"`
}

// UpdateEnablement changes the enablement of the policy
func (pa *RepPolicyAPI) UpdateEnablement() {
	e := enablementReq{}
	pa.DecodeJSONReq(&e)
	if e.Enabled != 0 && e.Enabled != 1 {
		pa.RenderError(http.StatusBadRequest, "invalid enabled value")
		return
	}

	if pa.policy.Enabled == e.Enabled {
		return
	}

	if err := dao.UpdateRepPolicyEnablement(pa.policyID, e.Enabled); err != nil {
		log.Errorf("Failed to update policy enablement in DB, error: %v", err)
		pa.RenderError(http.StatusInternalServerError, "Internal Error")
		return
	}

	if e.Enabled == 1 {
		go func() {
			if err := TriggerReplication(pa.policyID, "", nil, models.RepOpTransfer); err != nil {
				log.Errorf("failed to trigger replication of %d: %v", pa.policyID, err)
			} else {
				log.Infof("replication of %d triggered", pa.policyID)
			}
		}()
	} else {
		go func() {
			if err := postReplicationAction(pa.policyID, "stop"); err != nil {
				log.Errorf("failed to stop replication of %d: %v", pa.policyID, err)
			} else {
				log.Infof("try to stop replication of %d", pa.policyID)
			}
		}()
	}
}
