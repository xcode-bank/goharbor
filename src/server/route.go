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

package server

import (
	"github.com/astaxie/beego"
	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/core/api"
	"github.com/goharbor/harbor/src/core/config"
	"github.com/goharbor/harbor/src/core/controllers"
	"github.com/goharbor/harbor/src/core/service/notifications/admin"
	"github.com/goharbor/harbor/src/core/service/notifications/jobs"
	"github.com/goharbor/harbor/src/core/service/notifications/scheduler"
	"github.com/goharbor/harbor/src/core/service/token"
)

func registerRoutes() {
	// Controller API:
	beego.Router("/c/login", &controllers.CommonController{}, "post:Login")
	beego.Router("/c/log_out", &controllers.CommonController{}, "get:LogOut")
	beego.Router("/c/userExists", &controllers.CommonController{}, "post:UserExists")
	beego.Router(common.OIDCLoginPath, &controllers.OIDCController{}, "get:RedirectLogin")
	beego.Router("/c/oidc/onboard", &controllers.OIDCController{}, "post:Onboard")
	beego.Router(common.OIDCCallbackPath, &controllers.OIDCController{}, "get:Callback")

	beego.Router("/api/internal/configurations", &api.ConfigAPI{}, "get:GetInternalConfig;put:Put")
	beego.Router("/api/internal/renameadmin", &api.InternalAPI{}, "post:RenameAdmin")
	beego.Router("/api/internal/switchquota", &api.InternalAPI{}, "put:SwitchQuota")
	beego.Router("/api/internal/syncquota", &api.InternalAPI{}, "post:SyncQuota")

	beego.Router("/service/notifications/jobs/adminjob/:id([0-9]+)", &admin.Handler{}, "post:HandleAdminJob")
	beego.Router("/service/notifications/jobs/replication/:id([0-9]+)", &jobs.Handler{}, "post:HandleReplicationScheduleJob")
	beego.Router("/service/notifications/jobs/replication/task/:id([0-9]+)", &jobs.Handler{}, "post:HandleReplicationTask")
	beego.Router("/service/notifications/jobs/webhook/:id([0-9]+)", &jobs.Handler{}, "post:HandleNotificationJob")
	beego.Router("/service/notifications/jobs/retention/task/:id([0-9]+)", &jobs.Handler{}, "post:HandleRetentionTask")
	beego.Router("/service/notifications/schedules/:id([0-9]+)", &scheduler.Handler{}, "post:Handle")
	beego.Router("/service/notifications/jobs/scan/:uuid", &jobs.Handler{}, "post:HandleScan")

	beego.Router("/service/token", &token.Handler{})

	// chart repository services
	if config.WithChartMuseum() {
		chartRepositoryAPIType := &api.ChartRepositoryAPI{}
		beego.Router("/chartrepo/:repo/index.yaml", chartRepositoryAPIType, "get:GetIndexByRepo")
		beego.Router("/chartrepo/index.yaml", chartRepositoryAPIType, "get:GetIndex")
		beego.Router("/chartrepo/:repo/charts/:filename", chartRepositoryAPIType, "get:DownloadChart")
	}

	// Error pages
	beego.ErrorController(&controllers.ErrorController{})
}
