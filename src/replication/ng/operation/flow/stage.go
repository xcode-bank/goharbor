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

package flow

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goharbor/harbor/src/common/utils/log"
	adp "github.com/goharbor/harbor/src/replication/ng/adapter"
	"github.com/goharbor/harbor/src/replication/ng/dao/models"
	"github.com/goharbor/harbor/src/replication/ng/model"
	"github.com/goharbor/harbor/src/replication/ng/operation/execution"
	"github.com/goharbor/harbor/src/replication/ng/operation/scheduler"
	"github.com/goharbor/harbor/src/replication/ng/util"
)

// get/create the source registry, destination registry, source adapter and destination adapter
func initialize(policy *model.Policy) (adp.Adapter, adp.Adapter, error) {
	var srcAdapter, dstAdapter adp.Adapter
	var err error

	// create the source registry adapter
	srcFactory, err := adp.GetFactory(policy.SrcRegistry.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get adapter factory for registry type %s: %v", policy.SrcRegistry.Type, err)
	}
	srcAdapter, err = srcFactory(policy.SrcRegistry)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create adapter for source registry %s: %v", policy.SrcRegistry.URL, err)
	}

	// create the destination registry adapter
	dstFactory, err := adp.GetFactory(policy.DestRegistry.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get adapter factory for registry type %s: %v", policy.DestRegistry.Type, err)
	}
	dstAdapter, err = dstFactory(policy.DestRegistry)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create adapter for destination registry %s: %v", policy.DestRegistry.URL, err)
	}
	log.Debug("replication flow initialization completed")
	return srcAdapter, dstAdapter, nil
}

// fetch resources from the source registry
func fetchResources(adapter adp.Adapter, policy *model.Policy) ([]*model.Resource, error) {
	resTypes := []model.ResourceType{}
	filters := []*model.Filter{}
	for _, filter := range policy.Filters {
		if filter.Type != model.FilterTypeResource {
			filters = append(filters, filter)
			continue
		}
		resTypes = append(resTypes, filter.Value.(model.ResourceType))
	}
	if len(resTypes) == 0 {
		info, err := adapter.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get the adapter info: %v", err)
		}
		resTypes = append(resTypes, info.SupportedResourceTypes...)
	}

	resources := []*model.Resource{}
	// convert the adapter to different interfaces according to its required resource types
	for _, typ := range resTypes {
		var res []*model.Resource
		var err error
		if typ == model.ResourceTypeRepository {
			// images
			reg, ok := adapter.(adp.ImageRegistry)
			if !ok {
				return nil, fmt.Errorf("the adapter doesn't implement the ImageRegistry interface")
			}
			res, err = reg.FetchImages(policy.SrcNamespaces, filters)
		} else if typ == model.ResourceTypeChart {
			// charts
			reg, ok := adapter.(adp.ChartRegistry)
			if !ok {
				return nil, fmt.Errorf("the adapter doesn't implement the ChartRegistry interface")
			}
			res, err = reg.FetchCharts(policy.SrcNamespaces, filters)
		} else {
			return nil, fmt.Errorf("unsupported resource type %s", typ)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s: %v", typ, err)
		}
		resources = append(resources, res...)
		log.Debugf("fetch %s completed", typ)
	}

	log.Debug("fetch resources from the source registry completed")
	return resources, nil
}

// apply the filters to the resources and returns the filtered resources
func filterResources(resources []*model.Resource, filters []*model.Filter) ([]*model.Resource, error) {
	res := []*model.Resource{}
	for _, resource := range resources {
		match := true
		for _, filter := range filters {
			switch filter.Type {
			case model.FilterTypeResource:
				resourceType, ok := filter.Value.(string)
				if !ok {
					return nil, fmt.Errorf("%v is not a valid string", filter.Value)
				}
				if model.ResourceType(resourceType) != resource.Type {
					match = false
					break
				}
			case model.FilterTypeName:
				pattern, ok := filter.Value.(string)
				if !ok {
					return nil, fmt.Errorf("%v is not a valid string", filter.Value)
				}
				if resource.Metadata == nil {
					match = false
					break
				}
				// TODO filter only the repository part?
				m, err := util.Match(pattern, resource.Metadata.GetResourceName())
				if err != nil {
					return nil, err
				}
				if !m {
					match = false
					break
				}
			case model.FilterTypeTag:
				pattern, ok := filter.Value.(string)
				if !ok {
					return nil, fmt.Errorf("%v is not a valid string", filter.Value)
				}
				if resource.Metadata == nil {
					match = false
					break
				}
				versions := []string{}
				for _, version := range resource.Metadata.Vtags {
					m, err := util.Match(pattern, version)
					if err != nil {
						return nil, err
					}
					if m {
						versions = append(versions, version)
					}
				}
				if len(versions) == 0 {
					match = false
					break
				}
				// NOTE: the property "Vtags" of the origin resource struct is overrided here
				resource.Metadata.Vtags = versions
			case model.FilterTypeLabel:
			// TODO add support to label
			default:
				return nil, fmt.Errorf("unsupportted filter type: %v", filter.Type)
			}
		}
		if match {
			res = append(res, resource)
		}
	}
	log.Debug("filter resources completed")
	return res, nil
}

// assemble the destination resources by filling the metadata, registry and override properties
func assembleDestinationResources(adapter adp.Adapter, resources []*model.Resource,
	policy *model.Policy) ([]*model.Resource, error) {
	result := []*model.Resource{}
	var namespace *model.Namespace
	if len(policy.DestNamespace) > 0 {
		namespace = &model.Namespace{
			Name: policy.DestNamespace,
		}
	}
	for _, resource := range resources {
		metadata, err := adapter.ConvertResourceMetadata(resource.Metadata, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to convert the resource metadata of %s: %v", resource.Metadata.GetResourceName(), err)
		}
		res := &model.Resource{
			Type:         resource.Type,
			Metadata:     metadata,
			Registry:     policy.DestRegistry,
			ExtendedInfo: resource.ExtendedInfo,
			Deleted:      resource.Deleted,
			Override:     policy.Override,
		}
		result = append(result, res)
	}
	log.Debug("assemble the destination resources completed")
	return result, nil
}

// do the prepare work for pushing/uploading the resources: create the namespace or repository
func prepareForPush(adapter adp.Adapter, resources []*model.Resource) error {
	// TODO need to consider how to handle that both contains public/private namespace
	for _, resource := range resources {
		name := resource.Metadata.GetResourceName()
		if err := adapter.PrepareForPush(resource); err != nil {
			return fmt.Errorf("failed to do the prepare work for pushing/uploading %s: %v", name, err)
		}
		log.Debugf("the prepare work for pushing/uploading %s completed", name)
	}
	return nil
}

// preprocess
func preprocess(scheduler scheduler.Scheduler, srcResources, dstResources []*model.Resource) ([]*scheduler.ScheduleItem, error) {
	items, err := scheduler.Preprocess(srcResources, dstResources)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess the resources: %v", err)
	}
	log.Debug("preprocess the resources completed")
	return items, nil
}

// create task records in database
func createTasks(mgr execution.Manager, executionID int64, items []*scheduler.ScheduleItem) error {
	for _, item := range items {
		operation := "copy"
		if item.DstResource.Deleted {
			operation = "deletion"
		}

		task := &models.Task{
			ExecutionID:  executionID,
			Status:       models.TaskStatusInitialized,
			ResourceType: string(item.SrcResource.Type),
			SrcResource:  getResourceName(item.SrcResource),
			DstResource:  getResourceName(item.DstResource),
			Operation:    operation,
		}

		if item.DstResource.Invalid {
			task.Status = models.TaskStatusFailed
			task.EndTime = time.Now()
		}

		id, err := mgr.CreateTask(task)
		if err != nil {
			// if failed to create the task for one of the items,
			// the whole execution is marked as failure and all
			// the items will not be submitted
			return fmt.Errorf("failed to create task records for the execution %d: %v", executionID, err)
		}

		item.TaskID = id
		log.Debugf("task record %d for the execution %d created", id, executionID)
	}
	return nil
}

// schedule the replication tasks and update the task's status
// returns the count of tasks which have been scheduled and the error
func schedule(scheduler scheduler.Scheduler, executionMgr execution.Manager, items []*scheduler.ScheduleItem) (int, error) {
	results, err := scheduler.Schedule(items)
	if err != nil {
		return 0, fmt.Errorf("failed to schedule the tasks: %v", err)
	}

	allFailed := true
	n := len(results)
	for _, result := range results {
		// if the task is failed to be submitted, update the status of the
		// task as failure
		if result.Error != nil {
			log.Errorf("failed to schedule the task %d: %v", result.TaskID, result.Error)
			if err = executionMgr.UpdateTaskStatus(result.TaskID, models.TaskStatusFailed); err != nil {
				log.Errorf("failed to update the task status %d: %v", result.TaskID, err)
			}
			continue
		}
		allFailed = false
		// if the task is submitted successfully, update the status, job ID and start time
		if err = executionMgr.UpdateTaskStatus(result.TaskID, models.TaskStatusPending, models.TaskStatusInitialized); err != nil {
			log.Errorf("failed to update the task status %d: %v", result.TaskID, err)
		}
		if err = executionMgr.UpdateTask(&models.Task{
			ID:        result.TaskID,
			JobID:     result.JobID,
			StartTime: time.Now(),
		}, "JobID", "StartTime"); err != nil {
			log.Errorf("failed to update the task %d: %v", result.TaskID, err)
		}
		log.Debugf("the task %d scheduled", result.TaskID)
	}
	// if all the tasks are failed, return err
	if allFailed {
		return n, errors.New("all tasks are failed")
	}
	return n, nil
}

// return the name with format "res_name" or "res_name:[vtag1,vtag2,vtag3]"
// if the resource has vtags
func getResourceName(res *model.Resource) string {
	if res == nil {
		return ""
	}
	meta := res.Metadata
	if meta == nil {
		return ""
	}
	if len(meta.Vtags) == 0 {
		return meta.GetResourceName()
	}

	if len(meta.Vtags) <= 5 {
		return meta.GetResourceName() + ":[" + strings.Join(meta.Vtags, ",") + "]"
	}

	return fmt.Sprintf("%s:[%s ... %d in total]", meta.GetResourceName(), strings.Join(meta.Vtags[:5], ","), len(meta.Vtags))
}
