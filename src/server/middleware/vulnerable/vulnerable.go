package vulnerable

import (
	"github.com/goharbor/harbor/src/api/project"
	"github.com/goharbor/harbor/src/common/rbac"
	"github.com/goharbor/harbor/src/common/security"
	"net/http"
	"net/http/httptest"

	"github.com/goharbor/harbor/src/api/artifact"
	sc "github.com/goharbor/harbor/src/api/scan"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils/log"
	internal_errors "github.com/goharbor/harbor/src/internal/error"
	"github.com/goharbor/harbor/src/pkg/scan/report"
	v1 "github.com/goharbor/harbor/src/pkg/scan/rest/v1"
	"github.com/goharbor/harbor/src/pkg/scan/vuln"
	serror "github.com/goharbor/harbor/src/server/error"
	"github.com/goharbor/harbor/src/server/middleware"
	"github.com/pkg/errors"
)

// Middleware handle docker pull vulnerable check
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			doVulCheck, img, projectVulnerableSeverity, wl := validate(req)
			if !doVulCheck {
				next.ServeHTTP(rw, req)
				return
			}
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, req)
			// only enable vul policy check the response 200
			if rec.Result().StatusCode == http.StatusOK {
				// Invalid project ID
				if wl.ProjectID == 0 {
					err := errors.Errorf("project verification error: project %s", img.ProjectName)
					pkgE := internal_errors.New(err).WithCode(internal_errors.PROJECTPOLICYVIOLATION)
					serror.SendError(rw, pkgE)
					return
				}

				ctx := req.Context()
				art, err := artifact.Ctl.GetByReference(ctx, img.Repository, img.Digest, nil)
				if err != nil {
					// TODO: error handle
					return
				}

				cve := report.CVESet(wl.CVESet())
				summaries, err := sc.DefaultController.GetSummary(ctx, art, []string{v1.MimeTypeNativeReport}, report.WithCVEWhitelist(&cve))

				if err != nil {
					err = errors.Wrap(err, "middleware: vulnerable handler")
					pkgE := internal_errors.New(err).WithCode(internal_errors.PROJECTPOLICYVIOLATION)
					serror.SendError(rw, pkgE)
					return
				}

				rawSummary, ok := summaries[v1.MimeTypeNativeReport]
				// No report yet?
				if !ok {
					err = errors.Errorf("no scan report existing for the artifact: %s:%s@%s", img.Repository, img.Tag, img.Digest)
					pkgE := internal_errors.New(err).WithCode(internal_errors.PROJECTPOLICYVIOLATION)
					serror.SendError(rw, pkgE)
					return
				}

				summary := rawSummary.(*vuln.NativeReportSummary)

				// Do judgement
				if summary.Severity.Code() >= projectVulnerableSeverity.Code() {
					err = errors.Errorf("current image with '%q vulnerable' cannot be pulled due to configured policy in 'Prevent images with vulnerability severity of %q from running.' "+
						"Please contact your project administrator for help'", summary.Severity, projectVulnerableSeverity)
					pkgE := internal_errors.New(err).WithCode(internal_errors.PROJECTPOLICYVIOLATION)
					serror.SendError(rw, pkgE)
					return
				}

				// Print scannerPull CVE list
				if len(summary.CVEBypassed) > 0 {
					for _, cve := range summary.CVEBypassed {
						log.Infof("Vulnerable policy check: scannerPull CVE %s", cve)
					}
				}
			}
			middleware.CopyResp(rec, rw)
		})
	}
}

func validate(req *http.Request) (bool, middleware.ArtifactInfo, vuln.Severity, models.CVEWhitelist) {
	var vs vuln.Severity
	var wl models.CVEWhitelist
	var af middleware.ArtifactInfo
	err := middleware.EnsureArtifactDigest(req.Context())
	if err != nil {
		return false, af, vs, wl
	}
	af, ok := middleware.ArtifactInfoFromContext(req.Context())
	if !ok {
		return false, af, vs, wl
	}

	pro, err := project.Ctl.GetByName(req.Context(), af.ProjectName)
	if err != nil {
		return false, af, vs, wl
	}
	resource := rbac.NewProjectNamespace(pro.ProjectID).Resource(rbac.ResourceRepository)
	securityCtx, ok := security.FromContext(req.Context())
	if !ok {
		return false, af, vs, wl
	}
	if !securityCtx.Can(rbac.ActionScannerPull, resource) {
		return false, af, vs, wl
	}
	// Is vulnerable policy set?
	projectVulnerableEnabled, projectVulnerableSeverity, wl := middleware.GetPolicyChecker().VulnerablePolicy(af.ProjectName)
	if !projectVulnerableEnabled {
		return false, af, vs, wl
	}
	return true, af, projectVulnerableSeverity, wl
}
