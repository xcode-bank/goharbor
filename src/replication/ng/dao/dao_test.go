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
	"os"
	"testing"

	"github.com/goharbor/harbor/src/common/dao"
)

func TestMain(m *testing.M) {
	dao.PrepareTestForPostgresSQL()

	var code = m.Run()

	// clear test database
	var clearSqls = []string{
		`DROP TABLE "access", "access_log", "admin_job", "alembic_version", "clair_vuln_timestamp",
  "harbor_label", "harbor_resource_label", "harbor_user", "img_scan_job", "img_scan_overview",
  "job_log", "project", "project_member", "project_metadata", "properties", "registry",
  "replication_immediate_trigger", "replication_job", "replication_policy", "replication_policy_ng",
  "replication_target", "repository", "robot", "role", "schema_migrations", "user_group";`,
		`DROP FUNCTION "update_update_time_at_column"();`,
	}
	dao.PrepareTestData(clearSqls, nil)

	os.Exit(code)
}
