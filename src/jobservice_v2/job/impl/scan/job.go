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

package scan

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/job"
	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils/clair"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/jobservice_v2/env"
	"github.com/vmware/harbor/src/jobservice_v2/job/impl/utils"
)

// ClairJob is the struct to scan Harbor's Image with Clair
type ClairJob struct {
}

// MaxFails implements the interface in job/Interface
func (cj *ClairJob) MaxFails() uint {
	return 1
}

// ShouldRetry implements the interface in job/Interface
func (cj *ClairJob) ShouldRetry() bool {
	return false
}

// Validate implements the interface in job/Interface
func (cj *ClairJob) Validate(params map[string]interface{}) error {
	return nil
}

// Run implements the interface in job/Interface
func (cj *ClairJob) Run(ctx env.JobContext, params map[string]interface{}) error {
	// TODO: get logger from ctx?
	logger := log.DefaultLogger()

	jobParms, err := transformParam(params)
	if err != nil {
		logger.Errorf("Failed to prepare parms for scan job, error: %v", err)
		return err
	}

	repoClient, err := utils.NewRepositoryClientForJobservice(jobParms.Repository, jobParms.RegistryURL, jobParms.Secret, jobParms.TokenEndpoint)
	if err != nil {
		return err
	}
	imgDigest, _, payload, err := repoClient.PullManifest(jobParms.Tag, []string{schema2.MediaTypeManifest})
	if err != nil {
		logger.Errorf("Error pulling manifest for image %s:%s :%v", jobParms.Repository, jobParms.Tag, err)
		return err
	}
	token, err := utils.GetTokenForRepo(jobParms.Repository, jobParms.Secret, jobParms.TokenEndpoint)
	if err != nil {
		logger.Errorf("Failed to get token, error: %v", err)
		return err
	}
	layers, err := prepareLayers(payload, jobParms.RegistryURL, jobParms.Repository, token)
	if err != nil {
		logger.Errorf("Failed to prepare layers, error: %v", err)
		return err
	}
	clairClient := clair.NewClient(jobParms.ClairEndpoint, logger)

	for _, l := range layers {
		logger.Infof("Scanning Layer: %s, path: %s", l.Name, l.Path)
		if err := clairClient.ScanLayer(l); err != nil {
			logger.Errorf("Failed to scan layer: %s, error: %v", l.Name, err)
			return err
		}
	}

	layerName := layers[len(layers)-1].Name
	res, err := clairClient.GetResult(layerName)
	if err != nil {
		logger.Errorf("Failed to get result from Clair, error: %v", err)
		return err
	}
	compOverview, sev := clair.TransformVuln(res)
	err = dao.UpdateImgScanOverview(imgDigest, layerName, sev, compOverview)
	return err
}

func transformParam(params map[string]interface{}) (*job.ScanJobParms, error) {
	res := job.ScanJobParms{}
	parmsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(parmsBytes, &res)
	return &res, err
}

func prepareLayers(payload []byte, registryURL, repo, tk string) ([]models.ClairLayer, error) {
	layers := []models.ClairLayer{}
	manifest, _, err := distribution.UnmarshalManifest(schema2.MediaTypeManifest, payload)
	if err != nil {
		return layers, err
	}
	tokenHeader := map[string]string{"Connection": "close", "Authorization": fmt.Sprintf("Bearer %s", tk)}
	// form the chain by using the digests of all parent layers in the image, such that if another image is built on top of this image the layer name can be re-used.
	shaChain := ""
	for _, d := range manifest.References() {
		if d.MediaType == schema2.MediaTypeConfig {
			continue
		}
		shaChain += string(d.Digest) + "-"
		l := models.ClairLayer{
			Name:    fmt.Sprintf("%x", sha256.Sum256([]byte(shaChain))),
			Headers: tokenHeader,
			Format:  "Docker",
			Path:    utils.BuildBlobURL(registryURL, repo, string(d.Digest)),
		}
		if len(layers) > 0 {
			l.ParentName = layers[len(layers)-1].Name
		}
		layers = append(layers, l)
	}
	return layers, nil
}
