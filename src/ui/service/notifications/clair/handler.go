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

package clair

import (
	"encoding/json"
	"time"

	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils"
	"github.com/vmware/harbor/src/common/utils/clair"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/ui/api"
	"github.com/vmware/harbor/src/ui/config"
)

const (
	rescanInterval = 15 * time.Minute
)

var (
	clairClient = clair.NewClient(config.ClairEndpoint(), nil)
)

// Handler handles reqeust on /service/notifications/clair/, which listens to clair's notifications.
// When there's unexpected error it will silently fail without removing the notification such that it will be triggered again.
type Handler struct {
	api.BaseController
}

// Handle ...
func (h *Handler) Handle() {
	var ne models.ClairNotificationEnvelope
	if err := json.Unmarshal(h.Ctx.Input.CopyBody(1<<32), &ne); err != nil {
		log.Errorf("Failed to decode the request: %v", err)
		return
	}
	log.Debugf("Received notification from Clair, name: %s", ne.Notification.Name)
	notification, err := clairClient.GetNotification(ne.Notification.Name)
	if err != nil {
		log.Errorf("Failed to get notification details from Clair, name: %s, err: %v", ne.Notification.Name, err)
		return
	}
	ns := make(map[string]bool)
	if old := notification.Old; old != nil {
		if vuln := old.Vulnerability; vuln != nil {
			log.Debugf("old vulnerability namespace: %s", vuln.NamespaceName)
			ns[vuln.NamespaceName] = true
		}
	}
	if new := notification.New; new != nil {
		if vuln := new.Vulnerability; vuln != nil {
			log.Debugf("new vulnerability namespace: %s", vuln.NamespaceName)
			ns[vuln.NamespaceName] = true
		}
	}
	for k, v := range ns {
		if v {
			if err := dao.SetClairVulnTimestamp(k, time.Now()); err == nil {
				log.Debugf("Updated the timestamp for namespaces: %s", k)
			} else {
				log.Warningf("Failed to update the timestamp for namespaces: %s, error: %v", k, err)
			}
		}
	}
	if utils.ScanOverviewMarker().Check() {
		go func() {
			<-time.After(rescanInterval)
			l, err := dao.ListImgScanOverviews()
			if err != nil {
				log.Errorf("Failed to list scan overview records, error: %v", err)
				return
			}
			for _, e := range l {
				if err := clair.UpdateScanOverview(e.Digest, e.DetailsKey); err != nil {
					log.Errorf("Failed to refresh scan overview for image: %s", e.Digest)
				} else {
					log.Debugf("Refreshed scan overview for record with digest: %s", e.Digest)
				}
			}
		}()
		utils.ScanOverviewMarker().Mark()
	} else {
		log.Debugf("There is a rescan scheduled at %v already, skip.", utils.ScanOverviewMarker().Next())
	}
	if err := clairClient.DeleteNotification(ne.Notification.Name); err != nil {
		log.Warningf("Failed to remove notification from Clair, name: %s", ne.Notification.Name)
	} else {
		log.Debugf("Removed notification from Clair, name: %s", ne.Notification.Name)
	}
}
