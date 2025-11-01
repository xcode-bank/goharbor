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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vmware/harbor/dao"
	"github.com/vmware/harbor/models"
	"github.com/vmware/harbor/utils/log"
)

func checkProjectPermission(userID int, projectID int64) bool {
	exist, err := dao.IsAdminRole(userID)
	if err != nil {
		log.Errorf("Error occurred in IsAdminRole, error: %v", err)
		return false
	}
	if exist {
		return true
	}
	roleList, err := dao.GetUserProjectRoles(userID, projectID)
	if err != nil {
		log.Errorf("Error occurred in GetUserProjectRoles, error: %v", err)
		return false
	}
	return len(roleList) > 0
}

func checkUserExists(name string) int {
	u, err := dao.GetUser(models.User{Username: name})
	if err != nil {
		log.Errorf("Error occurred in GetUser, error: %v", err)
		return 0
	}
	if u != nil {
		return u.UserID
	}
	return 0
}

// TriggerReplication triggers the replication according to the policy
func TriggerReplication(policyID int64, repository, operation string) error {
	data := struct {
		PolicyID  int64  `json:"policy_id"`
		Repo      string `json:"repository"`
		Operation string `json:"operation"`
	}{
		PolicyID:  policyID,
		Repo:      repository,
		Operation: operation,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := buildReplicationURL()

	resp, err := http.DefaultClient.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return fmt.Errorf("%d %s", resp.StatusCode, string(b))
}

// GetPoliciesByRepository returns policies according the repository
func GetPoliciesByRepository(repository string) ([]*models.RepPolicy, error) {
	repository = strings.TrimSpace(repository)
	repository = strings.TrimRight(repository, "/")
	projectName := repository[:strings.LastIndex(repository, "/")]

	project, err := dao.GetProjectByName(projectName)
	if err != nil {
		return nil, err
	}

	policies, err := dao.GetRepPolicyByProject(project.ProjectID)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func TriggerReplicationByRepository(repository, operation string) {
	policies, err := GetPoliciesByRepository(repository)
	if err != nil {
		log.Errorf("failed to get policies for repository %s: %v", repository, err)
		return
	}

	for _, policy := range policies {
		if err := TriggerReplication(policy.ProjectID, repository, operation); err != nil {
			log.Errorf("failed to trigger replication of %d for %s: %v", policy.ID, repository, err)
		} else {
			log.Infof("replication of %d for %s triggered", policy.ID, repository)
		}
	}
}

func buildReplicationURL() string {
	return "http://job_service/api/replicationJobs"
}

func buildJobLogURL(jobID string) string {
	return fmt.Sprintf("http://job_service/api/replicationJobs/%s/log", jobID)
}
