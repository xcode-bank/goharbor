package immutable

import (
	"errors"
	"fmt"
	"net/http"

	common_util "github.com/goharbor/harbor/src/common/utils"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/controller/artifact"
	"github.com/goharbor/harbor/src/controller/tag"
	"github.com/goharbor/harbor/src/lib"
	internal_errors "github.com/goharbor/harbor/src/lib/error"
	serror "github.com/goharbor/harbor/src/server/error"
)

// Middleware ...
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if err := handlePush(req); err != nil {
				var e *ErrImmutable
				if errors.As(err, &e) {
					pkgE := internal_errors.New(e).WithCode(internal_errors.PreconditionCode)
					serror.SendError(rw, pkgE)
					return
				}
				pkgE := internal_errors.New(fmt.Errorf("error occurred when to handle request in immutable handler: %v", err)).WithCode(internal_errors.GeneralCode)
				serror.SendError(rw, pkgE)
				return
			}
			next.ServeHTTP(rw, req)
		})
	}
}

// handlePush ...
// If the pushing image matched by any of immutable rule, will have to whether it is the first time to push it,
// as the immutable rule only impacts the existing tag.
func handlePush(req *http.Request) error {
	none := lib.ArtifactInfo{}
	art := lib.GetArtifactInfo(req.Context())
	if art == none {
		return errors.New("cannot get the manifest information from request context")
	}

	af, err := artifact.Ctl.GetByReference(req.Context(), art.Repository, art.Tag, &artifact.Option{
		WithTag:   true,
		TagOption: &tag.Option{WithImmutableStatus: true},
	})
	if err != nil {
		log.Debugf("failed to list artifact, %v", err.Error())
		return nil
	}

	_, repoName := common_util.ParseRepository(art.Repository)
	for _, tag := range af.Tags {
		// push a existing immutable tag, reject th e request
		if tag.Name == art.Tag && tag.Immutable {
			return NewErrImmutable(repoName, art.Tag)
		}
	}

	return nil
}
