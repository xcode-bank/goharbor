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

package retention

import (
	"fmt"
	"time"

	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/pkg/project"
	"github.com/goharbor/harbor/src/pkg/repository"
	"github.com/goharbor/harbor/src/pkg/retention/policy"
	"github.com/goharbor/harbor/src/pkg/retention/q"
	"github.com/goharbor/harbor/src/pkg/scheduler"
)

// APIController to handle the requests related with retention
type APIController interface {
	// Handle the related hooks from the job service and launch the corresponding actions if needed
	//
	//  Arguments:
	//    PolicyID string         : uuid of the retention policy
	//    event *job.StatusChange : event object sent by job service
	//
	//  Returns:
	//    common error object if any errors occurred
	HandleHook(policyID string, event *job.StatusChange) error

	GetRetention(id int64) (*policy.Metadata, error)

	CreateRetention(p *policy.Metadata) error

	UpdateRetention(p *policy.Metadata) error

	DeleteRetention(id int64) error

	TriggerRetentionExec(policyID int64, trigger string, dryRun bool) error

	OperateRetentionExec(eid int64, action string) error

	ListRetentionExecs(policyID int64, query *q.Query) ([]*Execution, error)

	ListRetentionExecTasks(executionID int64, query *q.Query) ([]*Task, error)

	GetRetentionExecTaskLog(taskID int64) ([]byte, error)
}

// DefaultAPIController ...
type DefaultAPIController struct {
	manager        Manager
	launcher       Launcher
	projectManager project.Manager
	repositoryMgr  repository.Manager
	scheduler      scheduler.Scheduler
}

const (
	// SchedulerCallback ...
	SchedulerCallback = "SchedulerCallback"
)

// TriggerParam ...
type TriggerParam struct {
	PolicyID int64
	Trigger  string
}

// GetRetention Get Retention
func (r *DefaultAPIController) GetRetention(id int64) (*policy.Metadata, error) {
	return r.manager.GetPolicy(id)
}

// CreateRetention Create Retention
func (r *DefaultAPIController) CreateRetention(p *policy.Metadata) error {
	if p.Trigger.Kind == policy.TriggerKindSchedule {
		if p.Trigger.Settings != nil {
			cron, ok := p.Trigger.Settings[policy.TriggerSettingsCron]
			if ok {
				jobid, err := r.scheduler.Schedule(cron.(string), SchedulerCallback, TriggerParam{
					PolicyID: p.ID,
					Trigger:  ExecutionTriggerSchedule,
				})
				if err != nil {
					return err
				}
				p.Trigger.References[policy.TriggerReferencesJobid] = jobid
			}
		}
	}
	if _, err := r.manager.CreatePolicy(p); err != nil {
		return err
	}
	return nil
}

// UpdateRetention Update Retention
func (r *DefaultAPIController) UpdateRetention(p *policy.Metadata) error {
	p0, err := r.manager.GetPolicy(p.ID)
	if err != nil {
		return err
	}
	needUn := false
	needSch := false

	if p0.Trigger.Kind != p.Trigger.Kind {
		if p0.Trigger.Kind == policy.TriggerKindSchedule {
			needUn = true
		}

		if p.Trigger.Kind == policy.TriggerKindSchedule {
			needSch = true
		}
	} else {
		switch p.Trigger.Kind {
		case policy.TriggerKindSchedule:
			if p0.Trigger.Settings["cron"] != p.Trigger.Settings["cron"] {
				// unschedule old
				if len(p0.Trigger.References[policy.TriggerReferencesJobid].(string)) > 0 {
					needUn = true
				}
				// schedule new
				if len(p.Trigger.Settings[policy.TriggerSettingsCron].(string)) > 0 {
					// valid cron
					needSch = true
				}
			}
		case "":

		default:
			return fmt.Errorf("not support Trigger %s", p.Trigger.Kind)
		}
	}
	if needUn {
		err = r.scheduler.UnSchedule(p0.Trigger.References[policy.TriggerReferencesJobid].(int64))
		if err != nil {
			return err
		}
	}
	if needSch {
		jobid, err := r.scheduler.Schedule(p.Trigger.Settings[policy.TriggerSettingsCron].(string), SchedulerCallback, TriggerParam{
			PolicyID: p.ID,
			Trigger:  ExecutionTriggerSchedule,
		})
		if err != nil {
			return err
		}
		p.Trigger.References[policy.TriggerReferencesJobid] = jobid
	}

	return r.manager.UpdatePolicy(p)
}

// DeleteRetention Delete Retention
func (r *DefaultAPIController) DeleteRetention(id int64) error {
	p, err := r.manager.GetPolicy(id)
	if err != nil {
		return err
	}
	if p.Trigger.Kind == policy.TriggerKindSchedule && len(p.Trigger.Settings[policy.TriggerSettingsCron].(string)) > 0 {
		err = r.scheduler.UnSchedule(p.Trigger.References[policy.TriggerReferencesJobid].(int64))
		if err != nil {
			return err
		}
	}

	return r.manager.DeletePolicyAndExec(id)
}

// TriggerRetentionExec Trigger Retention Execution
func (r *DefaultAPIController) TriggerRetentionExec(policyID int64, trigger string, dryRun bool) error {
	p, err := r.manager.GetPolicy(policyID)
	if err != nil {
		return err
	}

	exec := &Execution{
		PolicyID:  policyID,
		StartTime: time.Now(),
		Status:    ExecutionStatusInProgress,
		Trigger:   trigger,
		DryRun:    dryRun,
	}
	id, err := r.manager.CreateExecution(exec)
	num, err := r.launcher.Launch(p, id, dryRun)
	if err != nil {
		return err
	}
	if num == 0 {
		exec := &Execution{
			ID:      id,
			EndTime: time.Now(),
			Status:  ExecutionStatusSucceed,
		}
		err = r.manager.UpdateExecution(exec)
		if err != nil {
			return err
		}
	}
	return err

}

// OperateRetentionExec Operate Retention Execution
func (r *DefaultAPIController) OperateRetentionExec(eid int64, action string) error {
	e, err := r.manager.GetExecution(eid)
	if err != nil {
		return err
	}
	if e == nil {
		return fmt.Errorf("execution %d not found", eid)
	}
	switch action {
	case "stop":
		return r.launcher.Stop(eid)
	default:
		return fmt.Errorf("not support action %s", action)
	}
}

// ListRetentionExecs List Retention Executions
func (r *DefaultAPIController) ListRetentionExecs(policyID int64, query *q.Query) ([]*Execution, error) {
	return r.manager.ListExecutions(policyID, query)
}

// ListRetentionExecTasks List Retention Execution Histories
func (r *DefaultAPIController) ListRetentionExecTasks(executionID int64, query *q.Query) ([]*Task, error) {
	q1 := &q.TaskQuery{
		ExecutionID: executionID,
		PageNumber:  query.PageNumber,
		PageSize:    query.PageSize,
	}
	return r.manager.ListTasks(q1)
}

// GetRetentionExecTaskLog Get Retention Execution Task Log
func (r *DefaultAPIController) GetRetentionExecTaskLog(taskID int64) ([]byte, error) {
	return r.manager.GetTaskLog(taskID)
}

// HandleHook HandleHook
func (r *DefaultAPIController) HandleHook(policyID string, event *job.StatusChange) error {
	panic("implement me")
}

// NewAPIController ...
func NewAPIController(projectManager project.Manager, repositoryMgr repository.Manager, scheduler scheduler.Scheduler, retentionLauncher Launcher) APIController {
	return &DefaultAPIController{
		manager:        NewManager(),
		launcher:       retentionLauncher,
		projectManager: projectManager,
		repositoryMgr:  repositoryMgr,
		scheduler:      scheduler,
	}
}
