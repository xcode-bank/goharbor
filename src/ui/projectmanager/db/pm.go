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

package db

import (
	"fmt"
	"regexp"
	"time"

	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/models"
	errutil "github.com/vmware/harbor/src/common/utils/error"
	"github.com/vmware/harbor/src/common/utils/log"
)

const dupProjectPattern = `Duplicate entry '\w+' for key 'name'`

// ProjectManager implements pm.PM interface based on database
type ProjectManager struct{}

// Get ...
func (p *ProjectManager) Get(projectIDOrName interface{}) (
	*models.Project, error) {
	switch projectIDOrName.(type) {
	case string:
		return dao.GetProjectByName(projectIDOrName.(string))
	case int64:
		return dao.GetProjectByID(projectIDOrName.(int64))
	default:
		return nil, fmt.Errorf("unsupported type of %v, must be string or int64", projectIDOrName)
	}
}

// Exist ...
func (p *ProjectManager) Exist(projectIDOrName interface{}) (bool, error) {
	project, err := p.Get(projectIDOrName)
	if err != nil {
		return false, err
	}
	return project != nil, nil
}

// IsPublic returns whether the project is public or not
func (p *ProjectManager) IsPublic(projectIDOrName interface{}) (bool, error) {
	project, err := p.Get(projectIDOrName)
	if err != nil {
		return false, err
	}

	if project == nil {
		return false, nil
	}

	return project.Public == 1, nil
}

// GetPublic returns all public projects
func (p *ProjectManager) GetPublic() ([]*models.Project, error) {
	t := true
	return p.GetAll(&models.ProjectQueryParam{
		Public: &t,
	})
}

// Create ...
func (p *ProjectManager) Create(project *models.Project) (int64, error) {
	if project == nil {
		return 0, fmt.Errorf("project is nil")
	}

	if len(project.Name) == 0 {
		return 0, fmt.Errorf("project name is nil")
	}

	if project.OwnerID == 0 {
		if len(project.OwnerName) == 0 {
			return 0, fmt.Errorf("owner ID and owner name are both nil")
		}

		user, err := dao.GetUser(models.User{
			Username: project.OwnerName,
		})
		if err != nil {
			return 0, err
		}
		if user == nil {
			return 0, fmt.Errorf("can not get owner whose name is %s", project.OwnerName)
		}
		project.OwnerID = user.UserID
	}

	t := time.Now()
	pro := &models.Project{
		Name:         project.Name,
		Public:       project.Public,
		OwnerID:      project.OwnerID,
		CreationTime: t,
		UpdateTime:   t,
	}

	id, err := dao.AddProject(*pro)
	if err != nil {
		dup, e := regexp.MatchString(dupProjectPattern, err.Error())
		if e != nil {
			log.Errorf("failed to match duplicate project pattern: %v", e)
		}

		if dup {
			err = errutil.ErrDupProject
		}

		return 0, err
	}

	return id, nil
}

// Delete ...
func (p *ProjectManager) Delete(projectIDOrName interface{}) error {
	id, ok := projectIDOrName.(int64)
	if !ok {
		project, err := p.Get(projectIDOrName)
		if err != nil {
			return err
		}
		id = project.ProjectID
	}

	return dao.DeleteProject(id)
}

// Update ...
func (p *ProjectManager) Update(projectIDOrName interface{},
	project *models.Project) error {
	id, ok := projectIDOrName.(int64)
	if !ok {
		pro, err := p.Get(projectIDOrName)
		if err != nil {
			return err
		}
		id = pro.ProjectID
	}
	return dao.ToggleProjectPublicity(id, project.Public)
}

// GetAll returns a project list according to the query parameters
func (p *ProjectManager) GetAll(query *models.ProjectQueryParam, base ...*models.BaseProjectCollection) (
	[]*models.Project, error) {
	return dao.GetProjects(query, base...)
}

// GetTotal returns the total count according to the query parameters
func (p *ProjectManager) GetTotal(query *models.ProjectQueryParam, base ...*models.BaseProjectCollection) (
	int64, error) {
	return dao.GetTotalOfProjects(query, base...)
}
