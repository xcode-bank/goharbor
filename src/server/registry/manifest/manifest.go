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

package manifest

import (
	"github.com/goharbor/harbor/src/api/artifact"
	"github.com/goharbor/harbor/src/api/repository"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/internal"
	"github.com/goharbor/harbor/src/server/registry/error"
	"github.com/goharbor/harbor/src/server/router"
	"github.com/gorilla/mux"
	"github.com/opencontainers/go-digest"
	"net/http"
	"net/http/httputil"
	"strings"
)

// NewHandler returns the handler to handler manifest requests
func NewHandler(proxy *httputil.ReverseProxy) http.Handler {
	return &handler{
		repoCtl: repository.Ctl,
		artCtl:  artifact.Ctl,
		proxy:   proxy,
	}
}

type handler struct {
	repoCtl repository.Controller
	artCtl  artifact.Controller
	proxy   *httputil.ReverseProxy
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodHead, http.MethodGet:
		h.get(w, req)
	case http.MethodDelete:
		h.delete(w, req)
	case http.MethodPut:
		h.put(w, req)
	}
}

// make sure the artifact exist before proxying the request to the backend registry
func (h *handler) get(w http.ResponseWriter, req *http.Request) {
	// check the existence in the database first
	vars := mux.Vars(req)
	reference := vars["reference"]
	artifact, err := h.artCtl.GetByReference(req.Context(), vars["name"], reference, nil)
	if err != nil {
		error.Handle(w, req, err)
		return
	}

	// the reference is tag, replace it with digest
	if _, err = digest.Parse(reference); err != nil {
		req = req.Clone(req.Context())
		req.URL.Path = strings.TrimSuffix(req.URL.Path, reference) + artifact.Digest
		req.URL.RawPath = ""
		req.URL.RawPath = req.URL.EscapedPath()
	}
	h.proxy.ServeHTTP(w, req)

	// TODO fire event(only for GET method), add access log in the event handler
}

func (h *handler) delete(w http.ResponseWriter, req *http.Request) {
	// just delete the artifact from database
	vars := mux.Vars(req)
	artifact, err := h.artCtl.GetByReference(req.Context(), vars["name"], vars["reference"], nil)
	if err != nil {
		error.Handle(w, req, err)
		return
	}
	if err = h.artCtl.Delete(req.Context(), artifact.ID); err != nil {
		error.Handle(w, req, err)
		return
	}

	// TODO fire event, add access log in the event handler
}

func (h *handler) put(w http.ResponseWriter, req *http.Request) {
	repository, err := router.Param(req.Context(), ":splat")
	if err != nil {
		error.Handle(w, req, err)
		return
	}
	reference, err := router.Param(req.Context(), ":reference")
	if err != nil {
		error.Handle(w, req, err)
		return
	}

	// make sure the repository exist before pushing the manifest
	_, repositoryID, err := h.repoCtl.Ensure(req.Context(), repository)
	if err != nil {
		error.Handle(w, req, err)
		return
	}

	buffer := internal.NewResponseBuffer(w)
	// proxy the req to the backend docker registry
	h.proxy.ServeHTTP(buffer, req)
	if !buffer.Success() {
		if _, err := buffer.Flush(); err != nil {
			log.Errorf("failed to flush: %v", err)
		}
		return
	}

	// When got the response from the backend docker registry, the manifest and
	// tag are both ready, so we don't need to handle the issue anymore:
	// https://github.com/docker/distribution/issues/2625

	var tags []string
	dgt := reference
	// the reference is tag, get the digest from the response header
	if _, err = digest.Parse(reference); err != nil {
		dgt = buffer.Header().Get("Docker-Content-Digest")
		tags = append(tags, reference)
	}

	_, _, err = h.artCtl.Ensure(req.Context(), repositoryID, dgt, tags...)
	if err != nil {
		error.Handle(w, req, err)
		return
	}

	// flush the origin response from the docker registry to the underlying response writer
	if _, err := buffer.Flush(); err != nil {
		log.Errorf("failed to flush: %v", err)
	}

	// TODO fire event, add access log in the event handler
}
