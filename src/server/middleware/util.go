package middleware

import (
	"context"
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/goharbor/harbor/src/api/artifact"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/core/config"
	"github.com/goharbor/harbor/src/core/promgr"
	"github.com/goharbor/harbor/src/pkg/scan/vuln"
	"github.com/goharbor/harbor/src/pkg/scan/whitelist"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
)

type contextKey string

const (
	// RepositorySubexp is the name for sub regex that maps to repository name in the url
	RepositorySubexp = "repository"
	// ReferenceSubexp is the name for sub regex that maps to reference (tag or digest) url
	ReferenceSubexp = "reference"
	// DigestSubexp is the name for sub regex that maps to digest in the url
	DigestSubexp = "digest"
	// ArtifactInfoKey the context key for artifact info
	ArtifactInfoKey = contextKey("artifactInfo")
	// manifestInfoKey the context key for manifest info
	manifestInfoKey = contextKey("ManifestInfo")
	// ScannerPullCtxKey the context key for robot account to bypass the pull policy check.
	ScannerPullCtxKey = contextKey("ScannerPullCheck")
)

var (
	// V2ManifestURLRe is the regular expression for matching request v2 handler to view/delete manifest
	V2ManifestURLRe = regexp.MustCompile(fmt.Sprintf(`^/v2/(?P<%s>%s)/manifests/(?P<%s>%s|%s)$`, RepositorySubexp, reference.NameRegexp.String(), ReferenceSubexp, reference.TagRegexp.String(), digest.DigestRegexp.String()))
	// V2TagListURLRe is the regular expression for matching request to v2 handler to list tags
	V2TagListURLRe = regexp.MustCompile(fmt.Sprintf(`^/v2/(?P<%s>%s)/tags/list`, RepositorySubexp, reference.NameRegexp.String()))
	// V2BlobURLRe is the regular expression for matching request to v2 handler to retrieve delete a blob
	V2BlobURLRe = regexp.MustCompile(fmt.Sprintf(`^/v2/(?P<%s>%s)/blobs/(?P<%s>%s)$`, RepositorySubexp, reference.NameRegexp.String(), DigestSubexp, digest.DigestRegexp.String()))
	// V2BlobUploadURLRe is the regular expression for matching the request to v2 handler to upload a blob, the upload uuid currently is not put into a group
	V2BlobUploadURLRe = regexp.MustCompile(fmt.Sprintf(`^/v2/(?P<%s>%s)/blobs/uploads[/a-zA-Z0-9\-_\.=]*$`, RepositorySubexp, reference.NameRegexp.String()))
	// V2CatalogURLRe is the regular expression for mathing the request to v2 handler to list catalog
	V2CatalogURLRe = regexp.MustCompile(`^/v2/_catalog$`)
)

// ManifestInfo ...
type ManifestInfo struct {
	ProjectID   int64
	ProjectName string
	Repository  string
	Tag         string
	Digest      string

	manifestExist     bool
	manifestExistErr  error
	manifestExistOnce sync.Once
}

// ManifestExists ...
func (info *ManifestInfo) ManifestExists(ctx context.Context) (bool, error) {
	info.manifestExistOnce.Do(func() {
		af, err := artifact.Ctl.GetByReference(ctx, info.Repository, info.Tag, nil)
		if err != nil {
			info.manifestExistErr = err
			return
		}
		info.manifestExist = true
		info.Digest = af.Digest
	})

	return info.manifestExist, info.manifestExistErr
}

// ArtifactInfo ...
type ArtifactInfo struct {
	Repository           string
	Reference            string
	ProjectName          string
	Digest               string
	BlobMountRepository  string
	BlobMountProjectName string
	BlobMountDigest      string
}

// ArtifactInfoFromContext returns the artifact info from context
func ArtifactInfoFromContext(ctx context.Context) (*ArtifactInfo, bool) {
	info, ok := ctx.Value(ArtifactInfoKey).(*ArtifactInfo)
	return info, ok
}

// NewManifestInfoContext returns context with manifest info
func NewManifestInfoContext(ctx context.Context, info *ManifestInfo) context.Context {
	return context.WithValue(ctx, manifestInfoKey, info)
}

// ManifestInfoFromContext returns manifest info from context
func ManifestInfoFromContext(ctx context.Context) (*ManifestInfo, bool) {
	info, ok := ctx.Value(manifestInfoKey).(*ManifestInfo)
	return info, ok
}

// NewScannerPullContext returns context with policy check info
func NewScannerPullContext(ctx context.Context, scannerPull bool) context.Context {
	return context.WithValue(ctx, ScannerPullCtxKey, scannerPull)
}

// ScannerPullFromContext returns whether to bypass policy check
func ScannerPullFromContext(ctx context.Context) (bool, bool) {
	info, ok := ctx.Value(ScannerPullCtxKey).(bool)
	return info, ok
}

// CopyResp ...
func CopyResp(rec *httptest.ResponseRecorder, rw http.ResponseWriter) {
	for k, v := range rec.Header() {
		rw.Header()[k] = v
	}
	rw.WriteHeader(rec.Result().StatusCode)
	rw.Write(rec.Body.Bytes())
}

// PolicyChecker checks the policy of a project by project name, to determine if it's needed to check the image's status under this project.
type PolicyChecker interface {
	// contentTrustEnabled returns whether a project has enabled content trust.
	ContentTrustEnabled(name string) bool
	// vulnerablePolicy  returns whether a project has enabled vulnerable, and the project's severity.
	VulnerablePolicy(name string) (bool, vuln.Severity, models.CVEWhitelist)
}

// PmsPolicyChecker ...
type PmsPolicyChecker struct {
	pm promgr.ProjectManager
}

// ContentTrustEnabled ...
func (pc PmsPolicyChecker) ContentTrustEnabled(name string) bool {
	project, err := pc.pm.Get(name)
	if err != nil {
		log.Errorf("Unexpected error when getting the project, error: %v", err)
		return true
	}
	if project == nil {
		log.Debugf("project %s not found", name)
		return false
	}
	return project.ContentTrustEnabled()
}

// VulnerablePolicy ...
func (pc PmsPolicyChecker) VulnerablePolicy(name string) (bool, vuln.Severity, models.CVEWhitelist) {
	project, err := pc.pm.Get(name)
	wl := models.CVEWhitelist{}
	if err != nil {
		log.Errorf("Unexpected error when getting the project, error: %v", err)
		return true, vuln.Unknown, wl
	}

	mgr := whitelist.NewDefaultManager()
	if project.ReuseSysCVEWhitelist() {
		w, err := mgr.GetSys()
		if err != nil {
			log.Error(errors.Wrap(err, "policy checker: vulnerable policy"))
		} else {
			wl = *w

			// Use the real project ID
			wl.ProjectID = project.ProjectID
		}
	} else {
		w, err := mgr.Get(project.ProjectID)
		if err != nil {
			log.Error(errors.Wrap(err, "policy checker: vulnerable policy"))
		} else {
			wl = *w
		}
	}

	return project.VulPrevented(), vuln.ParseSeverityVersion3(project.Severity()), wl
}

// NewPMSPolicyChecker returns an instance of an pmsPolicyChecker
func NewPMSPolicyChecker(pm promgr.ProjectManager) PolicyChecker {
	return &PmsPolicyChecker{
		pm: pm,
	}
}

// GetPolicyChecker ...
func GetPolicyChecker() PolicyChecker {
	return NewPMSPolicyChecker(config.GlobalProjectMgr)
}
