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

package blob

import (
	"context"
	"fmt"

	"github.com/docker/distribution"
	"github.com/garyburd/redigo/redis"
	"github.com/goharbor/harbor/src/common/utils/log"
	util "github.com/goharbor/harbor/src/common/utils/redis"
	ierror "github.com/goharbor/harbor/src/internal/error"
	"github.com/goharbor/harbor/src/internal/orm"
	"github.com/goharbor/harbor/src/pkg/blob"
)

var (
	// Ctl is a global blob controller instance
	Ctl = NewController()
)

// Controller defines the operations related with blobs
type Controller interface {
	// AssociateWithArtifact associate blobs with manifest.
	AssociateWithArtifact(ctx context.Context, blobDigests []string, artifactDigest string) error

	// AssociateWithProjectByID associate blob with project by blob id
	AssociateWithProjectByID(ctx context.Context, blobID int64, projectID int64) error

	// AssociateWithProjectByDigest associate blob with project by blob digest
	AssociateWithProjectByDigest(ctx context.Context, blobDigest string, projectID int64) error

	// Ensure create blob when it not exist.
	Ensure(ctx context.Context, digest string, contentType string, size int64) (int64, error)

	// Exist check blob exist by digest,
	// it check the blob associated with the artifact when `IsAssociatedWithArtifact` option provided,
	// and also check the blob associated with the project when `IsAssociatedWithProject` option provied.
	Exist(ctx context.Context, digest string, options ...Option) (bool, error)

	// Get get the blob by digest,
	// it check the blob associated with the artifact when `IsAssociatedWithArtifact` option provided,
	// and also check the blob associated with the project when `IsAssociatedWithProject` option provied.
	Get(ctx context.Context, digest string, options ...Option) (*blob.Blob, error)

	// Sync create blobs from `References` when they are not exist
	// and update the blob content type when they are exist,
	Sync(ctx context.Context, references []distribution.Descriptor) error

	// SetAcceptedBlobSize update the accepted size of stream upload blob.
	SetAcceptedBlobSize(sessionID string, size int64) error

	// GetAcceptedBlobSize returns the accepted size of stream upload blob.
	GetAcceptedBlobSize(sessionID string) (int64, error)
}

// NewController creates an instance of the default repository controller
func NewController() Controller {
	return &controller{
		blobMgr:   blob.Mgr,
		logPrefix: "[controller][blob]",
	}
}

type controller struct {
	blobMgr   blob.Manager
	logPrefix string
}

func (c *controller) AssociateWithArtifact(ctx context.Context, blobDigests []string, artifactDigest string) error {
	exist, err := c.blobMgr.IsAssociatedWithArtifact(ctx, artifactDigest, artifactDigest)
	if err != nil {
		return err
	}

	if exist {
		log.Infof("%s: artifact digest %s already exist, skip to associate blobs with the artifact", c.logPrefix, artifactDigest)
		return nil
	}

	for _, blobDigest := range blobDigests {
		_, err := c.blobMgr.AssociateWithArtifact(ctx, blobDigest, artifactDigest)
		if err != nil {
			return err
		}
	}

	// process manifest as blob
	_, err = c.blobMgr.AssociateWithArtifact(ctx, artifactDigest, artifactDigest)
	return err
}

func (c *controller) AssociateWithProjectByID(ctx context.Context, blobID int64, projectID int64) error {
	_, err := c.blobMgr.AssociateWithProject(ctx, blobID, projectID)
	return err
}

func (c *controller) AssociateWithProjectByDigest(ctx context.Context, blobDigest string, projectID int64) error {
	blob, err := c.blobMgr.Get(ctx, blobDigest)
	if err != nil {
		return err
	}

	_, err = c.blobMgr.AssociateWithProject(ctx, blob.ID, projectID)
	return err
}

func (c *controller) Get(ctx context.Context, digest string, options ...Option) (*blob.Blob, error) {
	if digest == "" {
		return nil, ierror.New(nil).WithCode(ierror.BadRequestCode).WithMessage("require Digest")
	}

	blob, err := c.blobMgr.Get(ctx, digest)
	if err != nil {
		return nil, err
	}

	opts := &Options{}
	for _, f := range options {
		f(opts)
	}

	if opts.ProjectID != 0 {
		exist, err := c.blobMgr.IsAssociatedWithProject(ctx, digest, opts.ProjectID)
		if err != nil {
			return nil, err
		}

		if !exist {
			return nil, ierror.NotFoundError(nil).WithMessage("blob %s is not associated with the project %d", digest, opts.ProjectID)
		}
	}

	if opts.ArtifactDigest != "" {
		exist, err := c.blobMgr.IsAssociatedWithArtifact(ctx, digest, opts.ArtifactDigest)
		if err != nil {
			return nil, err
		}

		if !exist {
			return nil, ierror.NotFoundError(nil).WithMessage("blob %s is not associated with the artifact %s", digest, opts.ArtifactDigest)
		}
	}

	return blob, nil
}

func (c *controller) Ensure(ctx context.Context, digest string, contentType string, size int64) (blobID int64, err error) {
	blob, err := c.blobMgr.Get(ctx, digest)
	if err == nil {
		return blob.ID, nil
	}

	if !ierror.IsNotFoundErr(err) {
		return 0, err
	}

	return c.blobMgr.Create(ctx, digest, contentType, size)
}

func (c *controller) Exist(ctx context.Context, digest string, options ...Option) (bool, error) {
	if digest == "" {
		return false, ierror.BadRequestError(nil).WithMessage("exist blob require digest")
	}

	_, err := c.Get(ctx, digest, options...)
	if err != nil {
		if ierror.IsNotFoundErr(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *controller) Sync(ctx context.Context, references []distribution.Descriptor) error {
	if len(references) == 0 {
		return nil
	}

	var digests []string
	for _, reference := range references {
		digests = append(digests, reference.Digest.String())
	}

	blobs, err := c.blobMgr.List(ctx, blob.ListParams{BlobDigests: digests})
	if err != nil {
		return err
	}

	mp := make(map[string]*blob.Blob, len(blobs))
	for _, blob := range blobs {
		mp[blob.Digest] = blob
	}

	var missing, updating []*blob.Blob
	for _, reference := range references {
		if exist, found := mp[reference.Digest.String()]; found {
			if exist.ContentType != reference.MediaType {
				exist.ContentType = reference.MediaType
				updating = append(updating, exist)
			}
		} else {
			missing = append(missing, &blob.Blob{
				Digest:      reference.Digest.String(),
				ContentType: reference.MediaType,
				Size:        reference.Size,
			})
		}
	}

	if len(updating) > 0 {
		orm.WithTransaction(func(ctx context.Context) error {
			for _, blob := range updating {
				if err := c.blobMgr.Update(ctx, blob); err != nil {
					log.Warningf("Failed to update blob %s, error: %v", blob.Digest, err)
					return err
				}
			}

			return nil
		})(ctx)
	}

	if len(missing) > 0 {
		for _, blob := range missing {
			if _, err := c.blobMgr.Create(ctx, blob.Digest, blob.ContentType, blob.Size); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *controller) SetAcceptedBlobSize(sessionID string, size int64) error {
	conn := util.DefaultPool().Get()
	defer conn.Close()

	key := fmt.Sprintf("upload:%s:size", sessionID)
	reply, err := redis.String(conn.Do("SET", key, size))
	if err != nil {
		return err
	}

	if reply != "OK" {
		return fmt.Errorf("bad reply value")
	}

	return nil
}

func (c *controller) GetAcceptedBlobSize(sessionID string) (int64, error) {
	conn := util.DefaultPool().Get()
	defer conn.Close()

	key := fmt.Sprintf("upload:%s:size", sessionID)
	size, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return size, nil
}
