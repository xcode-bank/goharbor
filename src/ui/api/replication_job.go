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

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/replication/core"
	api_models "github.com/vmware/harbor/src/ui/api/models"
	"github.com/vmware/harbor/src/ui/config"
	"github.com/vmware/harbor/src/ui/utils"
)

// RepJobAPI handles request to /api/replicationJobs /api/replicationJobs/:id/log
type RepJobAPI struct {
	BaseController
	jobID int64
}

// Prepare validates that whether user has system admin role
func (ra *RepJobAPI) Prepare() {
	ra.BaseController.Prepare()
	if !ra.SecurityCtx.IsAuthenticated() {
		ra.HandleUnauthorized()
		return
	}

	if !(ra.Ctx.Request.Method == http.MethodGet || ra.SecurityCtx.IsSysAdmin()) {
		ra.HandleForbidden(ra.SecurityCtx.GetUsername())
		return
	}

	if len(ra.GetStringFromPath(":id")) != 0 {
		id, err := ra.GetInt64FromPath(":id")
		if err != nil {
			ra.CustomAbort(http.StatusBadRequest, "ID is invalid")
		}
		ra.jobID = id
	}

}

// List filters jobs according to the parameters
func (ra *RepJobAPI) List() {

	policyID, err := ra.GetInt64("policy_id")
	if err != nil || policyID <= 0 {
		ra.CustomAbort(http.StatusBadRequest, "invalid policy_id")
	}

	policy, err := core.GlobalController.GetPolicy(policyID)
	if err != nil {
		log.Errorf("failed to get policy %d: %v", policyID, err)
		ra.CustomAbort(http.StatusInternalServerError, "")
	}

	if policy.ID == 0 {
		ra.CustomAbort(http.StatusNotFound, fmt.Sprintf("policy %d not found", policyID))
	}

	if !ra.SecurityCtx.HasAllPerm(policy.ProjectIDs[0]) {
		ra.HandleForbidden(ra.SecurityCtx.GetUsername())
		return
	}

	repository := ra.GetString("repository")
	statuses := ra.GetStrings("status")

	var startTime *time.Time
	startTimeStr := ra.GetString("start_time")
	if len(startTimeStr) != 0 {
		i, err := strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			ra.CustomAbort(http.StatusBadRequest, "invalid start_time")
		}
		t := time.Unix(i, 0)
		startTime = &t
	}

	var endTime *time.Time
	endTimeStr := ra.GetString("end_time")
	if len(endTimeStr) != 0 {
		i, err := strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			ra.CustomAbort(http.StatusBadRequest, "invalid end_time")
		}
		t := time.Unix(i, 0)
		endTime = &t
	}

	page, pageSize := ra.GetPaginationParams()

	jobs, total, err := dao.FilterRepJobs(policyID, repository, statuses,
		startTime, endTime, pageSize, pageSize*(page-1))
	if err != nil {
		log.Errorf("failed to filter jobs according policy ID %d, repository %s, status %v, start time %v, end time %v: %v",
			policyID, repository, statuses, startTime, endTime, err)
		ra.CustomAbort(http.StatusInternalServerError, "")
	}

	ra.SetPaginationHeader(total, page, pageSize)

	ra.Data["json"] = jobs
	ra.ServeJSON()
}

// Delete ...
func (ra *RepJobAPI) Delete() {
	if ra.jobID == 0 {
		ra.CustomAbort(http.StatusBadRequest, "id is nil")
	}

	job, err := dao.GetRepJob(ra.jobID)
	if err != nil {
		log.Errorf("failed to get job %d: %v", ra.jobID, err)
		ra.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	if job == nil {
		ra.CustomAbort(http.StatusNotFound, fmt.Sprintf("job %d not found", ra.jobID))
	}

	if job.Status == models.JobPending || job.Status == models.JobRunning {
		ra.CustomAbort(http.StatusBadRequest, fmt.Sprintf("job is %s, can not be deleted", job.Status))
	}

	if err = dao.DeleteRepJob(ra.jobID); err != nil {
		log.Errorf("failed to deleted job %d: %v", ra.jobID, err)
		ra.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

// GetLog ...
func (ra *RepJobAPI) GetLog() {
	if ra.jobID == 0 {
		ra.CustomAbort(http.StatusBadRequest, "id is nil")
	}

	job, err := dao.GetRepJob(ra.jobID)
	if err != nil {
		ra.HandleInternalServerError(fmt.Sprintf("failed to get replication job %d: %v", ra.jobID, err))
		return
	}

	if job == nil {
		ra.HandleNotFound(fmt.Sprintf("replication job %d not found", ra.jobID))
		return
	}

	policy, err := core.GlobalController.GetPolicy(job.PolicyID)
	if err != nil {
		ra.HandleInternalServerError(fmt.Sprintf("failed to get policy %d: %v", job.PolicyID, err))
		return
	}

	if !ra.SecurityCtx.HasAllPerm(policy.ProjectIDs[0]) {
		ra.HandleForbidden(ra.SecurityCtx.GetUsername())
		return
	}

	url := buildJobLogURL(strconv.FormatInt(ra.jobID, 10), ReplicationJobType)
	err = utils.RequestAsUI(http.MethodGet, url, nil, utils.NewJobLogRespHandler(&ra.BaseAPI))
	if err != nil {
		ra.RenderError(http.StatusInternalServerError, err.Error())
		return
	}
}

// StopJobs stop replication jobs for the policy
func (ra *RepJobAPI) StopJobs() {
	req := &api_models.StopJobsReq{}
	ra.DecodeJSONReqAndValidate(req)

	policy, err := core.GlobalController.GetPolicy(req.PolicyID)
	if err != nil {
		ra.HandleInternalServerError(fmt.Sprintf("failed to get policy %d: %v", req.PolicyID, err))
		return
	}

	if policy.ID == 0 {
		ra.CustomAbort(http.StatusNotFound, fmt.Sprintf("policy %d not found", req.PolicyID))
	}

	if err = config.GlobalJobserviceClient.StopReplicationJobs(req.PolicyID); err != nil {
		ra.HandleInternalServerError(fmt.Sprintf("failed to stop replication jobs of policy %d: %v", req.PolicyID, err))
		return
	}
}

//TODO:add Post handler to call job service API to submit jobs by policy
