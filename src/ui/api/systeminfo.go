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
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/vmware/harbor/src/common"
	"github.com/vmware/harbor/src/common/dao"
	clairdao "github.com/vmware/harbor/src/common/dao/clair"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils"
	"github.com/vmware/harbor/src/common/utils/clair"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/ui/config"
)

//SystemInfoAPI handle requests for getting system info /api/systeminfo
type SystemInfoAPI struct {
	BaseController
}

const defaultRootCert = "/etc/ui/ca/ca.crt"
const harborVersionFile = "/harbor/VERSION"

//SystemInfo models for system info.
type SystemInfo struct {
	HarborStorage Storage `json:"storage"`
}

//Storage models for storage.
type Storage struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

//namespaces stores all name spaces on Clair, it should be initialised only once.
type clairNamespaces struct {
	sync.RWMutex
	l     []string
	clair *clair.Client
}

func (n *clairNamespaces) get() ([]string, error) {
	n.Lock()
	defer n.Unlock()
	if len(n.l) == 0 {
		m := make(map[string]struct{})
		if n.clair == nil {
			n.clair = clair.NewClient(config.ClairEndpoint(), nil)
		}
		list, err := n.clair.ListNamespaces()
		if err != nil {
			return n.l, err
		}
		for _, n := range list {
			ns := strings.Split(n, ":")[0]
			m[ns] = struct{}{}
		}
		for k := range m {
			n.l = append(n.l, k)
		}
	}
	return n.l, nil
}

var (
	namespaces = &clairNamespaces{}
)

//GeneralInfo wraps common systeminfo for anonymous request
type GeneralInfo struct {
	WithNotary                  bool                             `json:"with_notary"`
	WithClair                   bool                             `json:"with_clair"`
	WithAdmiral                 bool                             `json:"with_admiral"`
	AdmiralEndpoint             string                           `json:"admiral_endpoint"`
	AuthMode                    string                           `json:"auth_mode"`
	RegistryURL                 string                           `json:"registry_url"`
	ProjectCreationRestrict     string                           `json:"project_creation_restriction"`
	SelfRegistration            bool                             `json:"self_registration"`
	HasCARoot                   bool                             `json:"has_ca_root"`
	HarborVersion               string                           `json:"harbor_version"`
	NextScanAll                 int64                            `json:"next_scan_all"`
	ClairVulnStatus             *models.ClairVulnerabilityStatus `json:"clair_vulnerability_status,omitempty"`
	RegistryStorageProviderName string                           `json:"registry_storage_provider_name"`
}

// validate for validating user if an admin.
func (sia *SystemInfoAPI) validate() {
	if !sia.SecurityCtx.IsAuthenticated() {
		sia.HandleUnauthorized()
		sia.StopRun()
	}

	if !sia.SecurityCtx.IsSysAdmin() {
		sia.HandleForbidden(sia.SecurityCtx.GetUsername())
		sia.StopRun()
	}
}

// GetVolumeInfo gets specific volume storage info.
func (sia *SystemInfoAPI) GetVolumeInfo() {
	sia.validate()

	capacity, err := config.AdminserverClient.Capacity()
	if err != nil {
		log.Errorf("failed to get capacity: %v", err)
		sia.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	systemInfo := SystemInfo{
		HarborStorage: Storage{
			Total: capacity.Total,
			Free:  capacity.Free,
		},
	}

	sia.Data["json"] = systemInfo
	sia.ServeJSON()
}

//GetCert gets default self-signed certificate.
func (sia *SystemInfoAPI) GetCert() {
	sia.validate()
	if _, err := os.Stat(defaultRootCert); err == nil {
		sia.Ctx.Output.Header("Content-Type", "application/octet-stream")
		sia.Ctx.Output.Header("Content-Disposition", "attachment; filename=ca.crt")
		http.ServeFile(sia.Ctx.ResponseWriter, sia.Ctx.Request, defaultRootCert)
	} else if os.IsNotExist(err) {
		log.Error("No certificate found.")
		sia.CustomAbort(http.StatusNotFound, "No certificate found.")
	} else {
		log.Errorf("Unexpected error: %v", err)
		sia.CustomAbort(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

// GetGeneralInfo returns the general system info, which is to be called by anonymous user
func (sia *SystemInfoAPI) GetGeneralInfo() {
	cfg, err := config.GetSystemCfg()
	if err != nil {
		log.Errorf("Error occured getting config: %v", err)
		sia.CustomAbort(http.StatusInternalServerError, "Unexpected error")
	}
	var registryURL string
	if l := strings.Split(cfg[common.ExtEndpoint].(string), "://"); len(l) > 1 {
		registryURL = l[1]
	} else {
		registryURL = l[0]
	}
	_, caStatErr := os.Stat(defaultRootCert)
	harborVersion := sia.getVersion()
	info := GeneralInfo{
		AdmiralEndpoint:             cfg[common.AdmiralEndpoint].(string),
		WithAdmiral:                 config.WithAdmiral(),
		WithNotary:                  config.WithNotary(),
		WithClair:                   config.WithClair(),
		AuthMode:                    cfg[common.AUTHMode].(string),
		ProjectCreationRestrict:     cfg[common.ProjectCreationRestriction].(string),
		SelfRegistration:            cfg[common.SelfRegistration].(bool),
		RegistryURL:                 registryURL,
		HasCARoot:                   caStatErr == nil,
		HarborVersion:               harborVersion,
		RegistryStorageProviderName: cfg[common.RegistryStorageProviderName].(string),
	}
	if info.WithClair {
		info.ClairVulnStatus = getClairVulnStatus()
		t := utils.ScanAllMarker().Next().UTC().Unix()
		if t < 0 {
			t = 0
		}
		info.NextScanAll = t
	}
	sia.Data["json"] = info
	sia.ServeJSON()
}

// getVersion gets harbor version.
func (sia *SystemInfoAPI) getVersion() string {
	version, err := ioutil.ReadFile(harborVersionFile)
	if err != nil {
		log.Errorf("Error occured getting harbor version: %v", err)
		return ""
	}
	return string(version[:])
}

func getClairVulnStatus() *models.ClairVulnerabilityStatus {
	res := &models.ClairVulnerabilityStatus{}
	last, err := clairdao.GetLastUpdate()
	if err != nil {
		log.Errorf("Failed to get last update from Clair DB, error: %v", err)
		res.OverallUTC = 0
	} else {
		res.OverallUTC = last
		log.Debugf("Clair vuln DB last update: %d", last)
	}
	details := []models.ClairNamespaceTimestamp{}
	if res.OverallUTC > 0 {
		l, err := dao.ListClairVulnTimestamps()
		if err != nil {
			log.Errorf("Failed to list Clair vulnerability timestamps, error:%v", err)
			return res
		}
		m := make(map[string]int64)
		for _, e := range l {
			ns := strings.Split(e.Namespace, ":")
			//only returns the latest time of one distro, i.e. unbuntu:14.04 and ubuntu:15.4 shares one timestamp
			el := e.LastUpdate.UTC().Unix()
			if ts, ok := m[ns[0]]; !ok || ts < el {
				m[ns[0]] = el
			}
		}
		list, err := namespaces.get()
		if err != nil {
			log.Errorf("Failed to get namespace list from Clair, error: %v", err)
		}
		//For namespaces not reported by notifier, the timestamp will be the overall db timestamp.
		for _, n := range list {
			if _, ok := m[n]; !ok {
				m[n] = res.OverallUTC
			}
		}
		for k, v := range m {
			e := models.ClairNamespaceTimestamp{
				Namespace: k,
				Timestamp: v,
			}
			details = append(details, e)
		}
	}
	res.Details = details
	return res
}
