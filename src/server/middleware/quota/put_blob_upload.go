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

package quota

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/goharbor/harbor/src/api/blob"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/pkg/distribution"
	"github.com/goharbor/harbor/src/pkg/types"
)

// PutBlobUploadMiddleware middleware to request storage resource for the project
func PutBlobUploadMiddleware() func(http.Handler) http.Handler {
	return RequestMiddleware(RequestConfig{
		ReferenceObject: projectReferenceObject,
		Resources:       putBlobUploadResources,
	})
}

func putBlobUploadResources(r *http.Request, reference, referenceID string) (types.ResourceList, error) {
	logPrefix := fmt.Sprintf("[middleware][%s][quota]", r.URL.Path)

	size, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil || size == 0 {
		size, err = blobController.GetAcceptedBlobSize(distribution.ParseSessionID(r.URL.Path))
	}
	if err != nil {
		log.Errorf("%s: get blob size failed, error: %v", logPrefix, err)
		return nil, err
	}

	projectID, _ := strconv.ParseInt(referenceID, 10, 64)

	digest := r.URL.Query().Get("digest")
	exist, err := blobController.Exist(r.Context(), digest, blob.IsAssociatedWithProject(projectID))
	if err != nil {
		log.Errorf("%s: checking blob %s is associated with project %d failed, error: %v", logPrefix, digest, projectID, err)
		return nil, err
	}

	if exist {
		return nil, nil
	}

	return types.ResourceList{types.ResourceStorage: size}, nil
}
