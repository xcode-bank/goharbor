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

package service

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/vmware/harbor/api"
	"github.com/vmware/harbor/dao"
	"github.com/vmware/harbor/models"
	"github.com/vmware/harbor/service/cache"
	"github.com/vmware/harbor/utils/log"
	"github.com/vmware/harbor/utils/registry"

	"github.com/astaxie/beego"
)

// NotificationHandler handles request on /service/notifications/, which listens to registry's events.
type NotificationHandler struct {
	beego.Controller
}

const manifestPattern = `^application/vnd.docker.distribution.manifest.v\d\+json`

// Post handles POST request, and records audit log or refreshes cache based on event.
func (n *NotificationHandler) Post() {
	var notification models.Notification
	//log.Info("Notification Handler triggered!\n")
	//	log.Infof("request body in string: %s", string(n.Ctx.Input.CopyBody()))
	err := json.Unmarshal(n.Ctx.Input.CopyBody(1<<32), &notification)

	if err != nil {
		log.Errorf("error while decoding json: %v", err)
		return
	}
	var username, action, repo, project, repoTag string
	var matched bool
	for _, e := range notification.Events {
		matched, err = regexp.MatchString(manifestPattern, e.Target.MediaType)
		if err != nil {
			log.Errorf("Failed to match the media type against pattern, error: %v", err)
			matched = false
		}
		if matched && (strings.HasPrefix(e.Request.UserAgent, "docker") ||
			strings.ToLower(strings.TrimSpace(e.Request.UserAgent)) == strings.ToLower(registry.UserAgent)) {
			username = e.Actor.Name
			action = e.Action
			repo = e.Target.Repository
			repoTag = e.Target.Tag
			log.Debugf("repo tag is : %v ", repoTag)

			if strings.Contains(repo, "/") {
				project = repo[0:strings.LastIndex(repo, "/")]
			}
			if username == "" {
				username = "anonymous"
			}

			if username == "job-service-user" {
				return
			}

			go dao.AccessLog(username, project, repo, repoTag, action)
			if action == "push" {
				go func() {
					err2 := cache.RefreshCatalogCache()
					if err2 != nil {
						log.Errorf("Error happens when refreshing cache: %v", err2)
					}
				}()

				go api.TriggerReplicationByRepository(repo, []string{repoTag}, models.RepOpTransfer)
			}
		}
	}

}

// Render returns nil as it won't render any template.
func (n *NotificationHandler) Render() error {
	return nil
}
