package robot

import (
	"github.com/goharbor/harbor/src/pkg/permission/types"
	"github.com/goharbor/harbor/src/pkg/robot2/model"
)

const (
	// LEVELSYSTEM ...
	LEVELSYSTEM = "system"
	// LEVELPROJECT ...
	LEVELPROJECT = "project"

	// SCOPESYSTEM ...
	SCOPESYSTEM = "/system"
	// SCOPEALLPROJECT ...
	SCOPEALLPROJECT = "/project/*"

	// ROBOTTYPE ...
	ROBOTTYPE = "robotaccount"
)

// Robot ...
type Robot struct {
	model.Robot
	ProjectName string
	Level       string
	Permissions []*Permission `json:"permissions"`
}

// setLevel, 0 is a system level robot, others are project level.
func (r *Robot) setLevel() {
	if r.ProjectID == 0 {
		r.Level = LEVELSYSTEM
	} else {
		r.Level = LEVELPROJECT
	}
}

// Permission ...
type Permission struct {
	Kind      string          `json:"kind"`
	Namespace string          `json:"namespace"`
	Access    []*types.Policy `json:"access"`
}

// Option ...
type Option struct {
	WithPermission bool
}
