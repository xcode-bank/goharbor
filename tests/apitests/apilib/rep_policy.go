/* 
 * Harbor API
 *
 * These APIs provide services for manipulating Harbor project.
 *
 * OpenAPI spec version: 0.3.0
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package apilib

type RepPolicy struct {

	// The policy ID.
	Id int64 `json:"id,omitempty"`

	// The project ID.
	ProjectId int64 `json:"project_id,omitempty"`

	// The project name.
	ProjectName string `json:"project_name,omitempty"`

	// The target ID.
	TargetId int64 `json:"target_id,omitempty"`

	// The target name.
	TargetName string `json:"target_name,omitempty"`

	// The policy name.
	Name string `json:"name,omitempty"`

	// The policy's enabled status.
	Enabled int32 `json:"enabled,omitempty"`

	// The description of the policy.
	Description string `json:"description,omitempty"`

	// The cron string for schedule job.
	CronStr string `json:"cron_str,omitempty"`

	// The start time of the policy.
	StartTime string `json:"start_time,omitempty"`

	// The create time of the policy.
	CreationTime string `json:"creation_time,omitempty"`

	// The update time of the policy.
	UpdateTime string `json:"update_time,omitempty"`
}
