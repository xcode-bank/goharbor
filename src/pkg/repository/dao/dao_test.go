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

package dao

import (
	"errors"
	"fmt"
	common_dao "github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/models"
	ierror "github.com/goharbor/harbor/src/internal/error"
	"github.com/goharbor/harbor/src/pkg/q"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	repository = fmt.Sprintf("library/%d", time.Now().Unix())
)

type daoTestSuite struct {
	suite.Suite
	dao DAO
	id  int64
}

func (d *daoTestSuite) SetupSuite() {
	d.dao = New()
	common_dao.PrepareTestForPostgresSQL()
}

func (d *daoTestSuite) SetupTest() {
	repository := &models.RepoRecord{
		Name:        repository,
		ProjectID:   1,
		Description: "",
	}
	id, err := d.dao.Create(nil, repository)
	d.Require().Nil(err)
	d.id = id
}

func (d *daoTestSuite) TearDownTest() {
	err := d.dao.Delete(nil, d.id)
	d.Require().Nil(err)
}

func (d *daoTestSuite) TestCount() {
	// nil query
	total, err := d.dao.Count(nil, nil)
	d.Require().Nil(err)
	d.True(total > 0)

	// query by name
	total, err = d.dao.Count(nil, &q.Query{
		Keywords: map[string]interface{}{
			"name": repository,
		},
	})
	d.Require().Nil(err)
	d.Equal(int64(1), total)
}

func (d *daoTestSuite) TestList() {
	// nil query
	repositories, err := d.dao.List(nil, nil)
	d.Require().Nil(err)
	found := false
	for _, repository := range repositories {
		if repository.RepositoryID == d.id {
			found = true
			break
		}
	}
	d.True(found)

	// query by name
	repositories, err = d.dao.List(nil, &q.Query{
		Keywords: map[string]interface{}{
			"name": repository,
		},
	})
	d.Require().Nil(err)
	d.Require().Equal(1, len(repositories))
	d.Equal(d.id, repositories[0].RepositoryID)
}

func (d *daoTestSuite) TestGet() {
	// get the non-exist repository
	_, err := d.dao.Get(nil, 10000)
	d.Require().NotNil(err)
	d.True(ierror.IsErr(err, ierror.NotFoundCode))

	// get the exist repository
	repository, err := d.dao.Get(nil, d.id)
	d.Require().Nil(err)
	d.Require().NotNil(repository)
	d.Equal(d.id, repository.RepositoryID)
}

func (d *daoTestSuite) TestCreate() {
	// the happy pass case is covered in Setup

	// conflict
	repository := &models.RepoRecord{
		Name:      repository,
		ProjectID: 1,
	}
	_, err := d.dao.Create(nil, repository)
	d.Require().NotNil(err)
	d.True(ierror.IsErr(err, ierror.ConflictCode))
}

func (d *daoTestSuite) TestDelete() {
	// the happy pass case is covered in TearDown

	// not exist
	err := d.dao.Delete(nil, 100021)
	d.Require().NotNil(err)
	var e *ierror.Error
	d.Require().True(errors.As(err, &e))
	d.Equal(ierror.NotFoundCode, e.Code)
}

func (d *daoTestSuite) TestUpdate() {
	// pass
	err := d.dao.Update(nil, &models.RepoRecord{
		RepositoryID: d.id,
		PullCount:    1,
	}, "PullCount")
	d.Require().Nil(err)

	repository, err := d.dao.Get(nil, d.id)
	d.Require().Nil(err)
	d.Require().NotNil(repository)
	d.Equal(int64(1), repository.PullCount)

	// not exist
	err = d.dao.Update(nil, &models.RepoRecord{
		RepositoryID: 10000,
	})
	d.Require().NotNil(err)
	var e *ierror.Error
	d.Require().True(errors.As(err, &e))
	d.Equal(ierror.NotFoundCode, e.Code)
}

func TestDaoTestSuite(t *testing.T) {
	suite.Run(t, &daoTestSuite{})
}
