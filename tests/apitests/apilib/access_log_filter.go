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

type AccessLogFilter struct {

	// Relevant user's name that accessed this project.
	Username string `json:"username,omitempty"`

	// Operation name specified when project created.
	Keywords string `json:"keywords,omitempty"`

	// Begin timestamp for querying access logs.
	BeginTimestamp int64 `json:"begin_timestamp,omitempty"`

	// End timestamp for querying accessl logs.
	EndTimestamp int64 `json:"end_timestamp,omitempty"`
}
