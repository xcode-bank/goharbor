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

package auth

import (
	"fmt"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/goharbor/harbor/src/common/http/modifier"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// NewAuthorizer creates an authorizer that can handle different auth schemes
func NewAuthorizer(credential Credential, client ...*http.Client) modifier.Modifier {
	authorizer := &authorizer{
		credential: credential,
	}
	if len(client) > 0 {
		authorizer.client = client[0]
	}
	if authorizer.client == nil {
		authorizer.client = http.DefaultClient
	}

	return authorizer
}

// authorizer authorizes the request with the provided credential.
// It determines the auth scheme of registry automatically and calls
// different underlying authorizers to do the auth work
type authorizer struct {
	sync.Mutex
	client     *http.Client
	url        *url.URL          // registry URL
	authorizer modifier.Modifier // the underlying authorizer
	credential Credential
}

func (a *authorizer) Modify(req *http.Request) error {
	// Nil URL means this is the first time the authorizer is called
	// Try to connect to the registry and determine the auth method
	if a.url == nil {
		// to avoid concurrent issue
		a.Lock()
		defer a.Unlock()
		if err := a.initialize(req.URL); err != nil {
			return err
		}
	}

	// check whether the request targets the registry
	if !a.isTarget(req) {
		return nil
	}

	return a.authorizer.Modify(req)
}

func (a *authorizer) initialize(u *url.URL) error {
	if a.url != nil {
		return nil
	}
	url, err := url.Parse(u.Scheme + "://" + u.Host + "/v2/")
	if err != nil {
		return err
	}
	a.url = url
	resp, err := a.client.Get(a.url.String())
	if err != nil {
		return err
	}
	challenges := ParseChallengeFromResponse(resp)
	// no challenge, mean no auth
	if len(challenges) == 0 {
		a.authorizer = &nullAuthorizer{}
		return nil
	}
	cm := map[string]challenge.Challenge{}
	for _, challenge := range challenges {
		cm[challenge.Scheme] = challenge
	}
	if _, exist := cm["basic"]; exist {
		a.authorizer = a.credential
		return nil
	}

	if _, exist := cm["bearer"]; exist {
		// TODO clean up the code of "StandardTokenAuthorizer"
		// TODO Currently, the checking of auth scheme is done twice, this can be avoided
		a.authorizer = NewStandardTokenAuthorizer(a.client, a.credential)
		return nil
	}
	return fmt.Errorf("unspported auth scheme: %v", challenges)
}

// Check whether the request targets to the registry.
// If doesn't, the request shouldn't be handled by the authorizer.
// e.g. the requests sent to backend storage(s3, etc.)
func (a *authorizer) isTarget(req *http.Request) bool {
	index := strings.Index(req.URL.Path, "/v2/")
	if index == -1 {
		return false
	}
	if req.URL.Host != a.url.Host || req.URL.Scheme != a.url.Scheme ||
		req.URL.Path[:index+4] != a.url.Path {
		return false
	}
	return true
}
