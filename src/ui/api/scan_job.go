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
	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/ui/utils"

	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ScanJobAPI handles request to /api/scanJobs/:id/log
type ScanJobAPI struct {
	BaseController
	jobID       int64
	projectName string
	jobUUID     string
}

// Prepare validates that whether user has read permission to the project of the repo the scan job scanned.
func (sj *ScanJobAPI) Prepare() {
	sj.BaseController.Prepare()
	if !sj.SecurityCtx.IsAuthenticated() {
		sj.HandleUnauthorized()
		return
	}
	id, err := sj.GetInt64FromPath(":id")
	if err != nil {
		sj.HandleBadRequest("invalid ID")
		return
	}
	sj.jobID = id

	data, err := dao.GetScanJob(id)
	if err != nil {
		log.Errorf("Failed to load job data for job: %d, error: %v", id, err)
		sj.CustomAbort(http.StatusInternalServerError, "Failed to get Job data")
	}
	projectName := strings.SplitN(data.Repository, "/", 2)[0]
	if !sj.SecurityCtx.HasReadPerm(projectName) {
		log.Errorf("User does not have read permission for project: %s", projectName)
		sj.HandleForbidden(sj.SecurityCtx.GetUsername())
	}
	sj.projectName = projectName
	sj.jobUUID = data.UUID
}

//GetLog ...
func (sj *ScanJobAPI) GetLog() {
	logBytes, err := utils.GetJobServiceClient().GetJobLog(sj.jobUUID)
	if err != nil {
		sj.HandleInternalServerError(fmt.Sprintf("Failed to get job logs, uuid: %s, error: %v", sj.jobUUID, err))
		return
	}
	sj.Ctx.ResponseWriter.Header().Set(http.CanonicalHeaderKey("Content-Length"), strconv.Itoa(len(logBytes)))
	sj.Ctx.ResponseWriter.Header().Set(http.CanonicalHeaderKey("Content-Type"), "text/plain")
	_, err = sj.Ctx.ResponseWriter.Write(logBytes)
	if err != nil {
		sj.HandleInternalServerError(fmt.Sprintf("Failed to write job logs, uuid: %s, error: %v", sj.jobUUID, err))
	}

}
