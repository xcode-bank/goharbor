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

package operation

import (
	"testing"

	"github.com/goharbor/harbor/src/replication/ng/config"
	"github.com/goharbor/harbor/src/replication/ng/dao/models"
	"github.com/goharbor/harbor/src/replication/ng/model"
	"github.com/goharbor/harbor/src/replication/ng/scheduler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakedExecutionManager struct{}

func (f *fakedExecutionManager) Create(*models.Execution) (int64, error) {
	return 1, nil
}
func (f *fakedExecutionManager) List(...*models.ExecutionQuery) (int64, []*models.Execution, error) {
	return 1, []*models.Execution{
		{
			ID: 1,
		},
	}, nil
}
func (f *fakedExecutionManager) Get(int64) (*models.Execution, error) {
	return &models.Execution{
		ID: 1,
	}, nil
}
func (f *fakedExecutionManager) Update(*models.Execution, ...string) error {
	return nil
}
func (f *fakedExecutionManager) Remove(int64) error {
	return nil
}
func (f *fakedExecutionManager) RemoveAll(int64) error {
	return nil
}
func (f *fakedExecutionManager) CreateTask(*models.Task) (int64, error) {
	return 1, nil
}
func (f *fakedExecutionManager) ListTasks(...*models.TaskQuery) (int64, []*models.Task, error) {
	return 1, []*models.Task{
		{
			ID: 1,
		},
	}, nil
}
func (f *fakedExecutionManager) GetTask(int64) (*models.Task, error) {
	return &models.Task{
		ID: 1,
	}, nil
}
func (f *fakedExecutionManager) UpdateTask(*models.Task, ...string) error {
	return nil
}
func (f *fakedExecutionManager) UpdateTaskStatus(int64, string, ...string) error {
	return nil
}
func (f *fakedExecutionManager) RemoveTask(int64) error {
	return nil
}
func (f *fakedExecutionManager) RemoveAllTasks(int64) error {
	return nil
}
func (f *fakedExecutionManager) GetTaskLog(int64) ([]byte, error) {
	return []byte("message"), nil
}

type fakedRegistryManager struct{}

func (f *fakedRegistryManager) Add(*model.Registry) (int64, error) {
	return 0, nil
}
func (f *fakedRegistryManager) List(...*model.RegistryQuery) (int64, []*model.Registry, error) {
	return 0, nil, nil
}
func (f *fakedRegistryManager) Get(id int64) (*model.Registry, error) {
	var registry *model.Registry
	switch id {
	case 1:
		registry = &model.Registry{
			ID:   1,
			Type: model.RegistryTypeHarbor,
		}
	}
	return registry, nil
}
func (f *fakedRegistryManager) GetByName(name string) (*model.Registry, error) {
	return nil, nil
}
func (f *fakedRegistryManager) Update(*model.Registry, ...string) error {
	return nil
}
func (f *fakedRegistryManager) Remove(int64) error {
	return nil
}
func (f *fakedRegistryManager) HealthCheck() error {
	return nil
}

type fakedScheduler struct{}

func (f *fakedScheduler) Preprocess(src []*model.Resource, dst []*model.Resource) ([]*scheduler.ScheduleItem, error) {
	items := []*scheduler.ScheduleItem{}
	for i, res := range src {
		items = append(items, &scheduler.ScheduleItem{
			SrcResource: res,
			DstResource: dst[i],
		})
	}
	return items, nil
}
func (f *fakedScheduler) Schedule(items []*scheduler.ScheduleItem) ([]*scheduler.ScheduleResult, error) {
	results := []*scheduler.ScheduleResult{}
	for _, item := range items {
		results = append(results, &scheduler.ScheduleResult{
			TaskID: item.TaskID,
			Error:  nil,
		})
	}
	return results, nil
}
func (f *fakedScheduler) Stop(id string) error {
	return nil
}

var ctl = NewController(&fakedExecutionManager{}, &fakedRegistryManager{}, &fakedScheduler{})

func TestStartReplication(t *testing.T) {
	config.Config = &config.Configuration{}
	// the resource contains Vtags whose length isn't 1
	policy := &model.Policy{}
	resource := &model.Resource{
		Type: model.ResourceTypeRepository,
		Metadata: &model.ResourceMetadata{
			Name:  "library/hello-world",
			Vtags: []string{"1.0", "2.0"},
		},
	}
	_, err := ctl.StartReplication(policy, resource)
	require.NotNil(t, err)

	// replicate resource deletion
	resource.Metadata.Vtags = []string{"1.0"}
	resource.Deleted = true
	id, err := ctl.StartReplication(policy, resource)
	require.Nil(t, err)
	assert.Equal(t, int64(1), id)

	// replicate resource copy
	resource.Deleted = false
	id, err = ctl.StartReplication(policy, resource)
	require.Nil(t, err)
	assert.Equal(t, int64(1), id)

	// nil resource
	id, err = ctl.StartReplication(policy, nil)
	require.Nil(t, err)
	assert.Equal(t, int64(1), id)
}

func TestStopReplication(t *testing.T) {
	err := ctl.StopReplication(1)
	require.Nil(t, err)
}

func TestListExecutions(t *testing.T) {
	n, executions, err := ctl.ListExecutions()
	require.Nil(t, err)
	assert.Equal(t, int64(1), n)
	assert.Equal(t, int64(1), executions[0].ID)
}

func TestGetExecution(t *testing.T) {
	execution, err := ctl.GetExecution(1)
	require.Nil(t, err)
	assert.Equal(t, int64(1), execution.ID)
}

func TestListTasks(t *testing.T) {
	n, tasks, err := ctl.ListTasks()
	require.Nil(t, err)
	assert.Equal(t, int64(1), n)
	assert.Equal(t, int64(1), tasks[0].ID)
}

func TestGetTask(t *testing.T) {
	task, err := ctl.GetTask(1)
	require.Nil(t, err)
	assert.Equal(t, int64(1), task.ID)
}

func TestUpdateTaskStatus(t *testing.T) {
	err := ctl.UpdateTaskStatus(1, "running")
	require.Nil(t, err)
}

func TestGetTaskLog(t *testing.T) {
	log, err := ctl.GetTaskLog(1)
	require.Nil(t, err)
	assert.Equal(t, "message", string(log))
}
