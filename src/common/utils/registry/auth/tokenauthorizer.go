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
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/common/utils/registry"
	token_util "github.com/vmware/harbor/src/ui/service/token"
)

const (
	latency int = 10 //second, the network latency when token is received
	scheme      = "bearer"
)

// Scope ...
type Scope struct {
	Type    string
	Name    string
	Actions []string
}

func (s *Scope) string() string {
	return fmt.Sprintf("%s:%s:%s", s.Type, s.Name, strings.Join(s.Actions, ","))
}

type tokenGenerator interface {
	generate(scopes []*Scope, endpoint string) (*models.Token, error)
}

// tokenAuthorizer implements registry.Modifier interface. It parses scopses
// from the request, generates authentication token and modifies the requset
// by adding the token
type tokenAuthorizer struct {
	registryURL  *url.URL // used to filter request
	generator    tokenGenerator
	client       *http.Client
	cachedTokens map[string]*models.Token
	sync.RWMutex
}

// add token to the request
func (t *tokenAuthorizer) Modify(req *http.Request) error {
	//only handle requests sent to registry
	goon, err := t.filterReq(req)
	if err != nil {
		return err
	}

	if !goon {
		log.Debugf("the request %s is not sent to registry, skip", req.URL.String())
		return nil
	}

	// parse scopes from request
	scopes := parseScopes(req)

	var token *models.Token
	// try to get token from cache if the request is for empty scope(login)
	// or single scope
	if len(scopes) <= 1 {
		key := ""
		if len(scopes) == 1 {
			key = scopes[0].string()
		}
		token = t.getCachedToken(key)
	}

	// request a new token if the token is null
	if token == nil {
		token, err = t.generator.generate(scopes, t.registryURL.String())
		if err != nil {
			return err
		}
		if token == nil {
			return nil
		}
		// only cache the token for empty scope(login) or single scope request
		if len(scopes) <= 1 {
			key := ""
			if len(scopes) == 1 {
				key = scopes[0].string()
			}
			t.updateCachedToken(key, token)
		}
	}

	req.Header.Add(http.CanonicalHeaderKey("Authorization"), fmt.Sprintf("Bearer %s", token.Token))

	return nil
}

// some requests are sent to backend storage, such as s3, this method filters
// the requests only sent to registry
func (t *tokenAuthorizer) filterReq(req *http.Request) (bool, error) {
	// the registryURL is nil when the first request comes, init it with
	// the scheme and host of the request which must be sent to the registry
	if t.registryURL == nil {
		u, err := url.Parse(buildPingURL(req.URL.Scheme + "://" + req.URL.Host))
		if err != nil {
			return false, err
		}
		t.registryURL = u
	}

	v2Index := strings.Index(req.URL.Path, "/v2/")
	if v2Index == -1 {
		return false, nil
	}

	if req.URL.Host != t.registryURL.Host || req.URL.Scheme != t.registryURL.Scheme ||
		req.URL.Path[:v2Index+4] != t.registryURL.Path {
		return false, nil
	}

	return true, nil
}

// parse scopes from the request according to its method, path and query string
func parseScopes(req *http.Request) []*Scope {
	scopes := []*Scope{}

	from := req.URL.Query().Get("from")
	if len(from) != 0 {
		scopes = append(scopes, &Scope{
			Type:    "repository",
			Name:    from,
			Actions: []string{"pull"},
		})
	}

	var scope *Scope
	path := strings.TrimRight(req.URL.Path, "/")
	repository := parseRepository(path)
	if len(repository) > 0 {
		// pull, push, delete blob/manifest
		scope = &Scope{
			Type: "repository",
			Name: repository,
		}
		switch req.Method {
		case http.MethodGet:
			scope.Actions = []string{"pull"}
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			scope.Actions = []string{"push"}
		case http.MethodDelete:
			scope.Actions = []string{"*"}
		default:
			scope = nil
			log.Warningf("unsupported method: %s", req.Method)
		}
	} else if catalog.MatchString(path) {
		// catalog
		scope = &Scope{
			Type:    "registry",
			Name:    "catalog",
			Actions: []string{"*"},
		}
	} else if base.MatchString(path) {
		// base
		scope = nil
	} else {
		// unknow
		log.Warningf("can not parse scope from the request: %s %s", req.Method, req.URL.Path)
	}

	if scope != nil {
		scopes = append(scopes, scope)
	}

	strs := []string{}
	for _, s := range scopes {
		strs = append(strs, s.string())
	}
	log.Debugf("scopses parsed from request: %s", strings.Join(strs, " "))

	return scopes
}

func (t *tokenAuthorizer) getCachedToken(scope string) *models.Token {
	t.RLock()
	defer t.RUnlock()
	token := t.cachedTokens[scope]
	if token == nil {
		return nil
	}

	issueAt, err := time.Parse(time.RFC3339, token.IssuedAt)
	if err != nil {
		log.Errorf("failed parse %s: %v", token.IssuedAt, err)
		return nil
	}

	if issueAt.Add(time.Duration(token.ExpiresIn-latency) * time.Second).Before(time.Now().UTC()) {
		return nil
	}

	log.Debug("get token from cache")
	return token
}

func (t *tokenAuthorizer) updateCachedToken(scope string, token *models.Token) {
	t.Lock()
	defer t.Unlock()
	t.cachedTokens[scope] = token
}

// ping returns the realm, service and error
func ping(client *http.Client, endpoint string) (string, string, error) {
	resp, err := client.Get(endpoint)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	challenges := ParseChallengeFromResponse(resp)
	for _, challenge := range challenges {
		if scheme == challenge.Scheme {
			realm := challenge.Parameters["realm"]
			service := challenge.Parameters["service"]
			return realm, service, nil
		}
	}

	log.Warningf("schemes %v are unsupportted", challenges)
	return "", "", nil
}

// NewStandardTokenAuthorizer returns a standard token authorizer. The authorizer will request a token
// from token server and add it to the origin request
// If customizedTokenService is set, the token request will be sent to it instead of the server get from authorizer
func NewStandardTokenAuthorizer(credential Credential, insecure bool,
	customizedTokenService ...string) registry.Modifier {
	client := &http.Client{
		Transport: registry.GetHTTPTransport(insecure),
		Timeout:   30 * time.Second,
	}

	generator := &standardTokenGenerator{
		credential: credential,
		client:     client,
	}

	// when the registry client is used inside Harbor, the token request
	// can be posted to token service directly rather than going through nginx.
	// If realm is set as the internal url of token service, this can resolve
	// two problems:
	// 1. performance issue
	// 2. the realm field returned by registry is an IP which can not reachable
	// inside Harbor
	if len(customizedTokenService) > 0 {
		generator.realm = customizedTokenService[0]
	}

	return &tokenAuthorizer{
		cachedTokens: make(map[string]*models.Token),
		generator:    generator,
		client:       client,
	}
}

// standardTokenGenerator implements interface tokenGenerator
type standardTokenGenerator struct {
	realm      string
	service    string
	credential Credential
	client     *http.Client
}

// get token from token service
func (s *standardTokenGenerator) generate(scopes []*Scope, endpoint string) (*models.Token, error) {
	// ping first if the realm or service is null
	if len(s.realm) == 0 || len(s.service) == 0 {
		realm, service, err := ping(s.client, endpoint)
		if err != nil {
			return nil, err
		}
		if len(realm) == 0 {
			log.Warning("empty realm, skip")
			return nil, nil
		}
		if len(s.realm) == 0 {
			s.realm = realm
		}
		s.service = service
	}

	return getToken(s.client, s.credential, s.realm, s.service, scopes)
}

// NewRawTokenAuthorizer returns a token authorizer which calls method to create
// token directly
func NewRawTokenAuthorizer(username, service string) registry.Modifier {
	generator := &rawTokenGenerator{
		service:  service,
		username: username,
	}

	return &tokenAuthorizer{
		cachedTokens: make(map[string]*models.Token),
		generator:    generator,
	}
}

// rawTokenGenerator implements interface tokenGenerator
type rawTokenGenerator struct {
	service  string
	username string
}

// generate token directly
func (r *rawTokenGenerator) generate(scopes []*Scope, endpoint string) (*models.Token, error) {
	strs := []string{}
	for _, scope := range scopes {
		strs = append(strs, scope.string())
	}
	token, expiresIn, issuedAt, err := token_util.RegistryTokenForUI(r.username, r.service, strs)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		Token:     token,
		ExpiresIn: expiresIn,
		IssuedAt:  issuedAt.Format(time.RFC3339),
	}, nil
}

func buildPingURL(endpoint string) string {
	return fmt.Sprintf("%s/v2/", endpoint)
}
