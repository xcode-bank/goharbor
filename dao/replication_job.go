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

package dao

import (
	"fmt"

	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/vmware/harbor/models"
)

// AddRepTarget ...
func AddRepTarget(target models.RepTarget) (int64, error) {
	o := GetOrmer()
	return o.Insert(&target)
}

// GetRepTarget ...
func GetRepTarget(id int64) (*models.RepTarget, error) {
	o := GetOrmer()
	t := models.RepTarget{ID: id}
	err := o.Read(&t)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

// GetRepTargetByName ...
func GetRepTargetByName(name string) (*models.RepTarget, error) {
	o := GetOrmer()
	t := models.RepTarget{Name: name}
	err := o.Read(&t, "Name")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

// DeleteRepTarget ...
func DeleteRepTarget(id int64) error {
	o := GetOrmer()
	_, err := o.Delete(&models.RepTarget{ID: id})
	return err
}

// UpdateRepTarget ...
func UpdateRepTarget(target models.RepTarget) error {
	o := GetOrmer()
	_, err := o.Update(&target, "URL", "Name", "Username", "Password")
	return err
}

// FilterRepTargets filters targets by name
func FilterRepTargets(name string) ([]*models.RepTarget, error) {
	o := GetOrmer()

	var args []interface{}

	sql := `select * from replication_target `
	if len(name) != 0 {
		sql += `where name like ? `
		args = append(args, "%"+name+"%")
	}
	sql += `order by creation_time`

	var targets []*models.RepTarget

	if _, err := o.Raw(sql, args).QueryRows(&targets); err != nil {
		return nil, err
	}

	return targets, nil
}

// AddRepPolicy ...
func AddRepPolicy(policy models.RepPolicy) (int64, error) {
	o := GetOrmer()
	sqlTpl := `insert into replication_policy (name, project_id, target_id, enabled, description, cron_str, start_time, creation_time, update_time ) values (?, ?, ?, ?, ?, ?, %s, NOW(), NOW())`
	var sql string
	if policy.Enabled == 1 {
		sql = fmt.Sprintf(sqlTpl, "NOW()")
	} else {
		sql = fmt.Sprintf(sqlTpl, "NULL")
	}
	p, err := o.Raw(sql).Prepare()
	if err != nil {
		return 0, err
	}
	r, err := p.Exec(policy.Name, policy.ProjectID, policy.TargetID, policy.Enabled, policy.Description, policy.CronStr)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	return id, err
}

// GetRepPolicy ...
func GetRepPolicy(id int64) (*models.RepPolicy, error) {
	o := GetOrmer()
	sql := `select * from replication_policy where id = ?`

	var policy models.RepPolicy

	if err := o.Raw(sql, id).QueryRow(&policy); err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &policy, nil
}

// FilterRepPolicies filters policies by name and project ID
func FilterRepPolicies(name string, projectID int64) ([]*models.RepPolicy, error) {
	o := GetOrmer()

	var args []interface{}

	sql := `select rp.id, rp.project_id, p.name as project_name, rp.target_id, 
				rt.name as target_name, rp.name, rp.enabled, rp.description,
				rp.cron_str, rp.start_time, rp.creation_time, rp.update_time  
			from replication_policy rp 
			join project p on rp.project_id=p.project_id 
			join replication_target rt on rp.target_id=rt.id `

	if len(name) != 0 && projectID != 0 {
		sql += `where rp.name like ? and rp.project_id = ? `
		args = append(args, "%"+name+"%")
		args = append(args, projectID)
	} else if len(name) != 0 {
		sql += `where rp.name like ? `
		args = append(args, "%"+name+"%")
	} else if projectID != 0 {
		sql += `where rp.project_id = ? `
		args = append(args, projectID)
	}

	sql += `order by rp.creation_time`

	var policies []*models.RepPolicy
	if _, err := o.Raw(sql, args).QueryRows(&policies); err != nil {
		return nil, err
	}
	return policies, nil
}

// GetRepPolicyByName ...
func GetRepPolicyByName(name string) (*models.RepPolicy, error) {
	o := GetOrmer()
	sql := `select * from replication_policy where name = ?`

	var policy models.RepPolicy

	if err := o.Raw(sql, name).QueryRow(&policy); err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &policy, nil
}

// GetRepPolicyByProject ...
func GetRepPolicyByProject(projectID int64) ([]*models.RepPolicy, error) {
	o := GetOrmer()
	sql := `select * from replication_policy where project_id = ?`

	var policies []*models.RepPolicy

	if _, err := o.Raw(sql, projectID).QueryRows(&policies); err != nil {
		return nil, err
	}

	return policies, nil
}

// GetRepPolicyByTarget ...
func GetRepPolicyByTarget(targetID int64) ([]*models.RepPolicy, error) {
	o := GetOrmer()
	sql := `select * from replication_policy where target_id = ?`

	var policies []*models.RepPolicy

	if _, err := o.Raw(sql, targetID).QueryRows(&policies); err != nil {
		return nil, err
	}

	return policies, nil
}

// UpdateRepPolicy ...
func UpdateRepPolicy(policy *models.RepPolicy) error {
	o := GetOrmer()
	_, err := o.Update(policy, "TargetID", "Name", "Enabled", "Description", "CronStr")
	return err
}

// DeleteRepPolicy ...
func DeleteRepPolicy(id int64) error {
	o := GetOrmer()
	_, err := o.Delete(&models.RepPolicy{ID: id})
	return err
}

// UpdateRepPolicyEnablement ...
func UpdateRepPolicyEnablement(id int64, enabled int) error {
	o := GetOrmer()
	p := models.RepPolicy{
		ID:      id,
		Enabled: enabled}
	_, err := o.Update(&p, "Enabled")

	return err
}

// EnableRepPolicy ...
func EnableRepPolicy(id int64) error {
	return UpdateRepPolicyEnablement(id, 1)
}

// DisableRepPolicy ...
func DisableRepPolicy(id int64) error {
	return UpdateRepPolicyEnablement(id, 0)
}

// AddRepJob ...
func AddRepJob(job models.RepJob) (int64, error) {
	o := GetOrmer()
	if len(job.Status) == 0 {
		job.Status = models.JobPending
	}
	if len(job.TagList) > 0 {
		job.Tags = strings.Join(job.TagList, ",")
	}
	return o.Insert(&job)
}

// GetRepJob ...
func GetRepJob(id int64) (*models.RepJob, error) {
	o := GetOrmer()
	j := models.RepJob{ID: id}
	err := o.Read(&j)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	genTagListForJob(&j)
	return &j, nil
}

// GetRepJobByPolicy ...
func GetRepJobByPolicy(policyID int64) ([]*models.RepJob, error) {
	var res []*models.RepJob
	_, err := repJobPolicyIDQs(policyID).All(&res)
	genTagListForJob(res...)
	return res, err
}

// FilterRepJobs filters jobs by repo and policy ID
func FilterRepJobs(repo string, policyID int64) ([]*models.RepJob, error) {
	o := GetOrmer()

	var args []interface{}

	sql := `select * from replication_job `

	if len(repo) != 0 && policyID != 0 {
		sql += `where repository like ? and policy_id = ? `
		args = append(args, "%"+repo+"%")
		args = append(args, policyID)
	} else if len(repo) != 0 {
		sql += `where repository like ? `
		args = append(args, "%"+repo+"%")
	} else if policyID != 0 {
		sql += `where policy_id = ? `
		args = append(args, policyID)
	}

	sql += `order by creation_time`

	var jobs []*models.RepJob
	if _, err := o.Raw(sql, args).QueryRows(&jobs); err != nil {
		return nil, err
	}

	genTagListForJob(jobs...)

	return jobs, nil
}

// GetRepJobToStop get jobs that are possibly being handled by workers of a certain policy.
func GetRepJobToStop(policyID int64) ([]*models.RepJob, error) {
	var res []*models.RepJob
	_, err := repJobPolicyIDQs(policyID).Filter("status__in", models.JobPending, models.JobRunning).All(&res)
	genTagListForJob(res...)
	return res, err
}

func repJobPolicyIDQs(policyID int64) orm.QuerySeter {
	o := GetOrmer()
	return o.QueryTable("replication_job").Filter("policy_id", policyID)
}

// DeleteRepJob ...
func DeleteRepJob(id int64) error {
	o := GetOrmer()
	_, err := o.Delete(&models.RepJob{ID: id})
	return err
}

// UpdateRepJobStatus ...
func UpdateRepJobStatus(id int64, status string) error {
	o := GetOrmer()
	j := models.RepJob{
		ID:     id,
		Status: status,
	}
	num, err := o.Update(&j, "Status")
	if num == 0 {
		err = fmt.Errorf("Failed to update replication job with id: %d %s", id, err.Error())
	}
	return err
}

func genTagListForJob(jobs ...*models.RepJob) {
	for _, j := range jobs {
		if len(j.Tags) > 0 {
			j.TagList = strings.Split(j.Tags, ",")
		}
	}
}
