//  Copyright Project Harbor Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package repoproxy

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/security"
	"github.com/goharbor/harbor/src/common/security/proxycachesecret"
	"github.com/goharbor/harbor/src/controller/project"
	"github.com/goharbor/harbor/src/controller/proxy"
	"github.com/goharbor/harbor/src/lib"
	"github.com/goharbor/harbor/src/lib/errors"
	httpLib "github.com/goharbor/harbor/src/lib/http"
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/replication/model"
	"github.com/goharbor/harbor/src/replication/registry"
	"github.com/goharbor/harbor/src/server/middleware"
)

var registryMgr = registry.NewDefaultManager()

// BlobGetMiddleware handle get blob request
func BlobGetMiddleware() func(http.Handler) http.Handler {
	return middleware.New(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		if err := handleBlob(w, r, next); err != nil {
			httpLib.SendError(w, err)
		}
	})
}

func handleBlob(w http.ResponseWriter, r *http.Request, next http.Handler) error {
	ctx := r.Context()
	art, p, proxyCtl, err := preCheck(ctx)
	if err != nil {
		return err
	}
	if !canProxy(p) || proxyCtl.UseLocalBlob(ctx, art) {
		next.ServeHTTP(w, r)
		return nil
	}
	size, reader, err := proxyCtl.ProxyBlob(ctx, p, art)
	if err != nil {
		return err
	}
	defer reader.Close()
	// Use io.CopyN to avoid out of memory when pulling big blob
	written, err := io.CopyN(w, reader, size)
	if err != nil {
		return err
	}
	if written != size {
		return errors.Errorf("The size mismatch, actual:%d, expected: %d", written, size)
	}
	setHeaders(w, size, "", art.Digest)
	return nil
}

func preCheck(ctx context.Context) (art lib.ArtifactInfo, p *models.Project, ctl proxy.Controller, err error) {
	none := lib.ArtifactInfo{}
	art = lib.GetArtifactInfo(ctx)
	if art == none {
		return none, nil, nil, errors.New("artifactinfo is not found").WithCode(errors.NotFoundCode)
	}
	ctl = proxy.ControllerInstance()
	p, err = project.Ctl.GetByName(ctx, art.ProjectName, project.Metadata(false))
	return
}

// ManifestMiddleware middleware handle request for get or head manifest
func ManifestMiddleware() func(http.Handler) http.Handler {
	return middleware.New(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		if err := handleManifest(w, r, next); err != nil {
			httpLib.SendError(w, err)
		}
	})
}

func handleManifest(w http.ResponseWriter, r *http.Request, next http.Handler) error {
	ctx := r.Context()
	art, p, proxyCtl, err := preCheck(ctx)
	if err != nil {
		return err
	}
	if !canProxy(p) {
		next.ServeHTTP(w, r)
		return nil
	}
	remote, err := proxy.NewRemoteHelper(p.RegistryID)
	if err != nil {
		return err
	}
	useLocal, err := proxyCtl.UseLocalManifest(ctx, art, remote)
	if err != nil {
		return err
	}
	if useLocal {
		next.ServeHTTP(w, r)
		return nil
	}
	log.Debugf("the tag is %v, digest is %v", art.Tag, art.Digest)
	if r.Method == http.MethodHead {
		err = proxyManifestHead(ctx, w, proxyCtl, p, art, remote)
	} else if r.Method == http.MethodGet {
		err = proxyManifestGet(ctx, w, proxyCtl, p, art, remote)
	}
	if err != nil {
		if errors.IsNotFoundErr(err) {
			return err
		}
		log.Warningf("Proxy to remote failed, fallback to local repo, error: %v", err)
		next.ServeHTTP(w, r)
	}
	return nil
}

func proxyManifestGet(ctx context.Context, w http.ResponseWriter, ctl proxy.Controller, p *models.Project, art lib.ArtifactInfo, remote proxy.RemoteInterface) error {
	man, err := ctl.ProxyManifest(ctx, art, remote)
	if err != nil {
		return err
	}
	ct, payload, err := man.Payload()
	if err != nil {
		return err
	}
	setHeaders(w, int64(len(payload)), ct, art.Digest)
	if _, err = w.Write(payload); err != nil {
		return err
	}
	return nil
}

func canProxy(p *models.Project) bool {
	if p.RegistryID < 1 {
		return false
	}
	reg, err := registryMgr.Get(p.RegistryID)
	if err != nil {
		log.Errorf("failed to get registry, error:%v", err)
		return false
	}
	if reg.Status != model.Healthy {
		log.Errorf("current registry is unhealthy, regID:%v, Name:%v, Status: %v", reg.ID, reg.Name, reg.Status)
	}
	return reg.Status == model.Healthy
}

func setHeaders(w http.ResponseWriter, size int64, mediaType string, dig string) {
	h := w.Header()
	h.Set("Content-Length", fmt.Sprintf("%v", size))
	if len(mediaType) > 0 {
		h.Set("Content-Type", mediaType)
	}
	h.Set("Docker-Content-Digest", dig)
	h.Set("Etag", dig)
}

// isProxySession check if current security context is proxy session
func isProxySession(ctx context.Context) bool {
	sc, ok := security.FromContext(ctx)
	if !ok {
		log.Error("Failed to get security context")
		return false
	}
	if sc.GetUsername() == proxycachesecret.ProxyCacheService {
		return true
	}
	return false
}

// DisableBlobAndManifestUploadMiddleware disable push artifact to a proxy project with a non-proxy session
func DisableBlobAndManifestUploadMiddleware() func(http.Handler) http.Handler {
	return middleware.New(func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		ctx := r.Context()
		art := lib.GetArtifactInfo(ctx)
		p, err := project.Ctl.GetByName(ctx, art.ProjectName)
		if err != nil {
			httpLib.SendError(w, err)
			return
		}
		if p.IsProxy() && !isProxySession(ctx) {
			httpLib.SendError(w,
				errors.DeniedError(
					errors.Errorf("can not push artifact to a proxy project: %v", p.Name)))
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func proxyManifestHead(ctx context.Context, w http.ResponseWriter, ctl proxy.Controller, p *models.Project, art lib.ArtifactInfo, remote proxy.RemoteInterface) error {
	exist, dig, err := ctl.HeadManifest(ctx, art, remote)
	if err != nil {
		return err
	}
	if !exist {
		return errors.NotFoundError(fmt.Errorf("The tag %v:%v is not found", art.Repository, art.Tag))
	}
	w.Header().Set("Docker-Content-Digest", dig)
	w.Header().Set("Etag", dig)
	return nil
}
