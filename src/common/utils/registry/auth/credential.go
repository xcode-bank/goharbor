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

package auth

import (
	"net/http"

	"github.com/vmware/harbor/src/common/http/modifier"
)

// Credential ...
type Credential modifier.Modifier

// Implements interface Credential
type basicAuthCredential struct {
	username string
	password string
}

// NewBasicAuthCredential ...
func NewBasicAuthCredential(username, password string) Credential {
	return &basicAuthCredential{
		username: username,
		password: password,
	}
}

func (b *basicAuthCredential) AddAuthorization(req *http.Request) {
	req.SetBasicAuth(b.username, b.password)
}

// implement github.com/vmware/harbor/src/common/http/modifier.Modifier
func (b *basicAuthCredential) Modify(req *http.Request) error {
	b.AddAuthorization(req)
	return nil
}

type cookieCredential struct {
	cookie *http.Cookie
}

// NewCookieCredential initialize a cookie based crendential handler, the cookie in parameter will be added to request to registry
// if this crendential is attached to a registry client.
func NewCookieCredential(c *http.Cookie) Credential {
	return &cookieCredential{
		cookie: c,
	}
}

func (c *cookieCredential) AddAuthorization(req *http.Request) {
	req.AddCookie(c.cookie)
}

// implement github.com/vmware/harbor/src/common/http/modifier.Modifier
func (c *cookieCredential) Modify(req *http.Request) error {
	c.AddAuthorization(req)
	return nil
}
