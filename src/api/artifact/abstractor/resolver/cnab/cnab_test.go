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

package cnab

import (
	"github.com/goharbor/harbor/src/common/models"
	ierror "github.com/goharbor/harbor/src/internal/error"
	"github.com/goharbor/harbor/src/pkg/artifact"
	"github.com/goharbor/harbor/src/testing/api/artifact/abstractor/blob"
	testingartifact "github.com/goharbor/harbor/src/testing/pkg/artifact"
	"github.com/goharbor/harbor/src/testing/pkg/repository"
	"github.com/stretchr/testify/suite"
	"testing"
)

type resolverTestSuite struct {
	suite.Suite
	resolver    *resolver
	repoMgr     *repository.FakeManager
	artMgr      *testingartifact.FakeManager
	blobFetcher *blob.FakeFetcher
}

func (r *resolverTestSuite) SetupTest() {
	r.repoMgr = &repository.FakeManager{}
	r.artMgr = &testingartifact.FakeManager{}
	r.blobFetcher = &blob.FakeFetcher{}
	r.resolver = &resolver{
		repoMgr:     r.repoMgr,
		argMgr:      r.artMgr,
		blobFetcher: r.blobFetcher,
	}

}

func (r *resolverTestSuite) TestResolveMetadata() {
	index := `{
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "digest": "sha256:b9616da7500f8c7c9a5e8d915714cd02d11bcc71ff5b4fd190bb77b1355c8549",
      "size": 193,
      "annotations": {
        "io.cnab.manifest.type": "config"
      }
    },
    {
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "digest": "sha256:a59a4e74d9cc89e4e75dfb2cc7ea5c108e4236ba6231b53081a9e2506d1197b6",
      "size": 942,
      "annotations": {
        "io.cnab.manifest.type": "invocation"
      }
    }
  ],
  "annotations": {
    "io.cnab.keywords": "[\"helloworld\",\"cnab\",\"tutorial\"]",
    "io.cnab.runtime_version": "v1.0.0",
    "org.opencontainers.artifactType": "application/vnd.cnab.manifest.v1",
    "org.opencontainers.image.authors": "[{\"name\":\"Jane Doe\",\"email\":\"jane.doe@example.com\",\"url\":\"https://example.com\"}]",
    "org.opencontainers.image.description": "A short description of your bundle",
    "org.opencontainers.image.title": "helloworld",
    "org.opencontainers.image.version": "0.1.1"
  }
}`

	manifest := `{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:e91b9dfcbbb3b88bac94726f276b89de46e4460b55f6e6d6f876e666b150ec5b",
    "size": 498
  },
  "layers": null
}`
	config := `{
  "description": "A short description of your bundle",
  "invocationImages": [
    {
      "contentDigest": "sha256:a59a4e74d9cc89e4e75dfb2cc7ea5c108e4236ba6231b53081a9e2506d1197b6",
      "image": "cnab/helloworld:0.1.1",
      "imageType": "docker",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "size": 942
    }
  ],
  "keywords": [
    "helloworld",
    "cnab",
    "tutorial"
  ],
  "maintainers": [
    {
      "email": "jane.doe@example.com",
      "name": "Jane Doe",
      "url": "https://example.com"
    }
  ],
  "name": "helloworld",
  "schemaVersion": "v1.0.0",
  "version": "0.1.1"
}`
	art := &artifact.Artifact{}
	r.artMgr.On("GetByDigest").Return(&artifact.Artifact{ID: 1}, nil)
	r.repoMgr.On("Get").Return(&models.RepoRecord{}, nil)
	r.blobFetcher.On("FetchManifest").Return("", []byte(manifest), nil)
	r.blobFetcher.On("FetchLayer").Return([]byte(config), nil)
	err := r.resolver.ResolveMetadata(nil, []byte(index), art)
	r.Require().Nil(err)
	r.Len(art.References, 2)
	r.Equal("0.1.1", art.ExtraAttrs["version"].(string))
	r.Equal("helloworld", art.ExtraAttrs["name"].(string))
}

func (r *resolverTestSuite) TestResolveAddition() {
	_, err := r.resolver.ResolveAddition(nil, nil, "")
	r.Require().NotNil(err)
	r.True(ierror.IsErr(err, ierror.BadRequestCode))
}

func (r *resolverTestSuite) TestGetArtifactType() {
	r.Assert().Equal(ArtifactTypeCNAB, r.resolver.GetArtifactType())
}

func (r *resolverTestSuite) TestListAdditionTypes() {
	r.Nil(r.resolver.ListAdditionTypes())
}

func TestResolverTestSuite(t *testing.T) {
	suite.Run(t, &resolverTestSuite{})
}
