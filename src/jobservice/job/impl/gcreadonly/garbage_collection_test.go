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

package gcreadonly

import (
	"os"
	"testing"

	"github.com/goharbor/harbor/src/common/config"
	"github.com/goharbor/harbor/src/common/models"
	commom_regctl "github.com/goharbor/harbor/src/common/registryctl"
	"github.com/goharbor/harbor/src/controller/project"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/pkg/artifact"
	"github.com/goharbor/harbor/src/pkg/artifactrash/model"
	pkg_blob "github.com/goharbor/harbor/src/pkg/blob/models"
	artifacttesting "github.com/goharbor/harbor/src/testing/controller/artifact"
	projecttesting "github.com/goharbor/harbor/src/testing/controller/project"
	mockjobservice "github.com/goharbor/harbor/src/testing/jobservice"
	"github.com/goharbor/harbor/src/testing/mock"
	trashtesting "github.com/goharbor/harbor/src/testing/pkg/artifactrash"
	"github.com/goharbor/harbor/src/testing/pkg/blob"
	"github.com/goharbor/harbor/src/testing/registryctl"
	"github.com/stretchr/testify/suite"
)

type gcTestSuite struct {
	suite.Suite
	artifactCtl       *artifacttesting.Controller
	artrashMgr        *trashtesting.FakeManager
	registryCtlClient *registryctl.Mockclient
	projectCtl        *projecttesting.Controller
	blobMgr           *blob.Manager

	originalProjectCtl project.Controller

	regCtlInit  func()
	setReadOnly func(cfgMgr *config.CfgManager, switcher bool) error
	getReadOnly func(cfgMgr *config.CfgManager) (bool, error)
}

func (suite *gcTestSuite) SetupTest() {
	suite.artifactCtl = &artifacttesting.Controller{}
	suite.artrashMgr = &trashtesting.FakeManager{}
	suite.registryCtlClient = &registryctl.Mockclient{}
	suite.blobMgr = &blob.Manager{}
	suite.projectCtl = &projecttesting.Controller{}

	suite.originalProjectCtl = project.Ctl
	project.Ctl = suite.projectCtl

	regCtlInit = func() { commom_regctl.RegistryCtlClient = suite.registryCtlClient }
	setReadOnly = func(cfgMgr *config.CfgManager, switcher bool) error { return nil }
	getReadOnly = func(cfgMgr *config.CfgManager) (bool, error) { return true, nil }
}

func (suite *gcTestSuite) TearDownTest() {
	project.Ctl = suite.originalProjectCtl
}

func (suite *gcTestSuite) TestMaxFails() {
	gc := &GarbageCollector{}
	suite.Equal(uint(1), gc.MaxFails())
}

func (suite *gcTestSuite) TestShouldRetry() {
	gc := &GarbageCollector{}
	suite.False(gc.ShouldRetry())
}

func (suite *gcTestSuite) TestValidate() {
	gc := &GarbageCollector{}
	suite.Nil(gc.Validate(nil))
}

func (suite *gcTestSuite) TestDeleteCandidates() {
	ctx := &mockjobservice.MockJobContext{}
	logger := &mockjobservice.MockJobLogger{}
	ctx.On("GetLogger").Return(logger)

	suite.artrashMgr.On("Flush").Return(nil)
	suite.artifactCtl.On("List").Return([]*artifact.Artifact{
		{
			ID:           1,
			RepositoryID: 1,
		},
	}, nil)
	suite.artifactCtl.On("Delete").Return(nil)
	suite.artrashMgr.On("Filter").Return([]model.ArtifactTrash{}, nil)

	gc := &GarbageCollector{
		artCtl:     suite.artifactCtl,
		artrashMgr: suite.artrashMgr,
	}
	suite.Nil(gc.deleteCandidates(ctx))
}

func (suite *gcTestSuite) TestRemoveUntaggedBlobs() {
	ctx := &mockjobservice.MockJobContext{}
	logger := &mockjobservice.MockJobLogger{}
	ctx.On("GetLogger").Return(logger)

	mock.OnAnything(suite.projectCtl, "List").Return([]*models.Project{
		{
			ProjectID: 1234,
			Name:      "test GC",
		},
	}, nil)

	mock.OnAnything(suite.blobMgr, "List").Return([]*pkg_blob.Blob{
		{
			ID:     1234,
			Digest: "sha256:1234",
			Size:   1234,
		},
	}, nil)

	mock.OnAnything(suite.blobMgr, "CleanupAssociationsForProject").Return(nil)

	gc := &GarbageCollector{
		blobMgr: suite.blobMgr,
	}

	suite.NotPanics(func() {
		gc.removeUntaggedBlobs(ctx)
	})
}

func (suite *gcTestSuite) TestInit() {
	ctx := &mockjobservice.MockJobContext{}
	logger := &mockjobservice.MockJobLogger{}
	mock.OnAnything(ctx, "Get").Return("core url", true)
	ctx.On("GetLogger").Return(logger)
	ctx.On("OPCommand").Return(job.NilCommand, true)

	gc := &GarbageCollector{
		registryCtlClient: suite.registryCtlClient,
	}
	params := map[string]interface{}{
		"delete_untagged": true,
		"redis_url_reg":   "redis url",
	}
	suite.Nil(gc.init(ctx, params))
	suite.True(gc.deleteUntagged)

	params = map[string]interface{}{
		"delete_untagged": "unsupported",
		"redis_url_reg":   "redis url",
	}
	suite.Nil(gc.init(ctx, params))
	suite.True(gc.deleteUntagged)

	params = map[string]interface{}{
		"delete_untagged": false,
		"redis_url_reg":   "redis url",
	}
	suite.Nil(gc.init(ctx, params))
	suite.False(gc.deleteUntagged)

	params = map[string]interface{}{
		"redis_url_reg": "redis url",
	}
	suite.Nil(gc.init(ctx, params))
	suite.True(gc.deleteUntagged)
}

func (suite *gcTestSuite) TestStop() {
	ctx := &mockjobservice.MockJobContext{}
	logger := &mockjobservice.MockJobLogger{}
	mock.OnAnything(ctx, "Get").Return("core url", true)
	ctx.On("GetLogger").Return(logger)
	ctx.On("OPCommand").Return(job.StopCommand, true)

	gc := &GarbageCollector{
		registryCtlClient: suite.registryCtlClient,
	}
	params := map[string]interface{}{
		"delete_untagged": true,
		"redis_url_reg":   "redis url",
	}
	suite.Nil(gc.init(ctx, params))

	ctx = &mockjobservice.MockJobContext{}
	mock.OnAnything(ctx, "Get").Return("core url", true)
	ctx.On("OPCommand").Return(job.StopCommand, false)
	suite.Nil(gc.init(ctx, params))

	ctx = &mockjobservice.MockJobContext{}
	mock.OnAnything(ctx, "Get").Return("core url", true)
	ctx.On("OPCommand").Return(job.NilCommand, true)
	suite.Nil(gc.init(ctx, params))

}

func (suite *gcTestSuite) TestRun() {
	ctx := &mockjobservice.MockJobContext{}
	logger := &mockjobservice.MockJobLogger{}
	ctx.On("GetLogger").Return(logger)
	ctx.On("OPCommand").Return(job.NilCommand, true)
	mock.OnAnything(ctx, "Get").Return("core url", true)

	suite.artrashMgr.On("Flush").Return(nil)
	suite.artifactCtl.On("List").Return([]*artifact.Artifact{
		{
			ID:           1,
			RepositoryID: 1,
		},
	}, nil)
	suite.artifactCtl.On("Delete").Return(nil)
	suite.artrashMgr.On("Filter").Return([]model.ArtifactTrash{}, nil)

	mock.OnAnything(suite.projectCtl, "List").Return([]*models.Project{
		{
			ProjectID: 12345,
			Name:      "test GC",
		},
	}, nil)

	mock.OnAnything(suite.blobMgr, "List").Return([]*pkg_blob.Blob{
		{
			ID:     12345,
			Digest: "sha256:12345",
			Size:   12345,
		},
	}, nil)

	mock.OnAnything(suite.blobMgr, "CleanupAssociationsForProject").Return(nil)

	gc := &GarbageCollector{
		artCtl:            suite.artifactCtl,
		artrashMgr:        suite.artrashMgr,
		cfgMgr:            config.NewInMemoryManager(),
		blobMgr:           suite.blobMgr,
		registryCtlClient: suite.registryCtlClient,
	}
	params := map[string]interface{}{
		"delete_untagged": false,
		// ToDo add a redis testing pkg, we do have a 'localhost' redis server in UT
		"redis_url_reg": "redis://localhost:6379",
	}

	suite.Nil(gc.Run(ctx, params))
}

func TestGCTestSuite(t *testing.T) {
	os.Setenv("UTTEST", "true")
	suite.Run(t, &gcTestSuite{})
}
