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

package api

import (
	"net/http"
)

// InternalAPI handles request of harbor admin...
type InternalAPI struct {
	BaseController
}

// Prepare validates the URL and parms
func (ia *InternalAPI) Prepare() {
	ia.BaseController.Prepare()
	if !ia.SecurityCtx.IsAuthenticated() {
		ia.HandleUnauthorized()
		return
	}
	if !ia.SecurityCtx.IsSysAdmin() {
		ia.HandleForbidden(ia.SecurityCtx.GetUsername())
		return
	}
}

// SyncRegistry ...
func (ia *InternalAPI) SyncRegistry() {
	err := SyncRegistry()
	if err != nil {
		ia.CustomAbort(http.StatusInternalServerError, "internal error")
	}
}
