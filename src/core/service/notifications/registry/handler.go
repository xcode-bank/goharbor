// Copyright 2018 Project Harbor Authors
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

package registry

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/goharbor/harbor/src/common/dao"
	clairdao "github.com/goharbor/harbor/src/common/dao/clair"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/core/api"
	"github.com/goharbor/harbor/src/core/config"
	coreutils "github.com/goharbor/harbor/src/core/utils"
	"github.com/goharbor/harbor/src/replication"
	"github.com/goharbor/harbor/src/replication/adapter"
	rep_event "github.com/goharbor/harbor/src/replication/event"
	"github.com/goharbor/harbor/src/replication/model"
)

// NotificationHandler handles request on /service/notifications/, which listens to registry's events.
type NotificationHandler struct {
	api.BaseController
}

const manifestPattern = `^application/vnd.docker.distribution.manifest.v\d\+(json|prettyjws)`
const vicPrefix = "vic/"

// Post handles POST request, and records audit log or refreshes cache based on event.
func (n *NotificationHandler) Post() {
	var notification models.Notification
	err := json.Unmarshal(n.Ctx.Input.CopyBody(1<<32), &notification)

	if err != nil {
		log.Errorf("failed to decode notification: %v", err)
		return
	}

	events, err := filterEvents(&notification)
	if err != nil {
		log.Errorf("failed to filter events: %v", err)
		return
	}

	for _, event := range events {
		repository := event.Target.Repository
		project, _ := utils.ParseRepository(repository)
		tag := event.Target.Tag
		action := event.Action

		user := event.Actor.Name
		if len(user) == 0 {
			user = "anonymous"
		}

		pro, err := config.GlobalProjectMgr.Get(project)
		if err != nil {
			log.Errorf("failed to get project by name %s: %v", project, err)
			return
		}
		if pro == nil {
			log.Warningf("project %s not found", project)
			continue
		}

		go func() {
			if err := dao.AddAccessLog(models.AccessLog{
				Username:  user,
				ProjectID: pro.ProjectID,
				RepoName:  repository,
				RepoTag:   tag,
				Operation: action,
				OpTime:    time.Now(),
			}); err != nil {
				log.Errorf("failed to add access log: %v", err)
			}
		}()

		if action == "push" {
			go func() {
				exist := dao.RepositoryExists(repository)
				if exist {
					return
				}
				log.Debugf("Add repository %s into DB.", repository)
				repoRecord := models.RepoRecord{
					Name:      repository,
					ProjectID: pro.ProjectID,
				}
				if err := dao.AddRepository(repoRecord); err != nil {
					log.Errorf("Error happens when adding repository: %v", err)
				}
			}()
			if !coreutils.WaitForManifestReady(repository, tag, 5) {
				log.Errorf("Manifest for image %s:%s is not ready, skip the follow up actions.", repository, tag)
				return
			}

			// TODO: handle image delete event and chart event
			go func() {
				e := &rep_event.Event{
					Type: rep_event.EventTypeImagePush,
					Resource: &model.Resource{
						Type: model.ResourceTypeImage,
						Metadata: &model.ResourceMetadata{
							Repository: &model.Repository{
								Name: repository,
								// TODO filling the metadata
							},
							Vtags: []string{tag},
						},
					},
				}
				if err := replication.EventHandler.Handle(e); err != nil {
					log.Errorf("failed to handle event: %v", err)
				}
			}()

			if autoScanEnabled(pro) {
				last, err := clairdao.GetLastUpdate()
				if err != nil {
					log.Errorf("Failed to get last update from Clair DB, error: %v, the auto scan will be skipped.", err)
				} else if last == 0 {
					log.Infof("The Vulnerability data is not ready in Clair DB, the auto scan will be skipped, error %v", err)
				} else if err := coreutils.TriggerImageScan(repository, tag); err != nil {
					log.Warningf("Failed to scan image, repository: %s, tag: %s, error: %v", repository, tag, err)
				}
			}
		}
		if action == "pull" {
			go func() {
				log.Debugf("Increase the repository %s pull count.", repository)
				if err := dao.IncreasePullCount(repository); err != nil {
					log.Errorf("Error happens when increasing pull count: %v", repository)
				}
			}()
		}
	}
}

func filterEvents(notification *models.Notification) ([]*models.Event, error) {
	events := []*models.Event{}

	for _, event := range notification.Events {
		log.Debugf("receive an event: \n----ID: %s \n----target: %s:%s \n----digest: %s \n----action: %s \n----mediatype: %s \n----user-agent: %s", event.ID, event.Target.Repository,
			event.Target.Tag, event.Target.Digest, event.Action, event.Target.MediaType, event.Request.UserAgent)

		isManifest, err := regexp.MatchString(manifestPattern, event.Target.MediaType)
		if err != nil {
			log.Errorf("failed to match the media type against pattern: %v", err)
			continue
		}

		if !isManifest {
			continue
		}

		if checkEvent(&event) {
			events = append(events, &event)
			log.Debugf("add event to collection: %s", event.ID)
			continue
		}
	}

	return events, nil
}

func checkEvent(event *models.Event) bool {
	// push action
	if event.Action == "push" {
		return true
	}
	// if it is pull action, check the user-agent
	userAgent := strings.ToLower(strings.TrimSpace(event.Request.UserAgent))
	if userAgent == "harbor-registry-client" || userAgent == strings.ToLower(adapter.UserAgentReplication) {
		return false
	}
	return true
}

func autoScanEnabled(project *models.Project) bool {
	if !config.WithClair() {
		log.Debugf("Auto Scan disabled because Harbor is not deployed with Clair")
		return false
	}

	return project.AutoScan()
}

// Render returns nil as it won't render any template.
func (n *NotificationHandler) Render() error {
	return nil
}
