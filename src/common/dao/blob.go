package dao

import (
	"fmt"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils/log"
	"strings"
	"time"
)

// AddBlob ...
func AddBlob(blob *models.Blob) (int64, error) {
	now := time.Now()
	blob.CreationTime = now
	id, err := GetOrmer().Insert(blob)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return 0, ErrDupRows
		}
		return 0, err
	}
	return id, nil
}

// GetBlob ...
func GetBlob(digest string) (*models.Blob, error) {
	o := GetOrmer()
	qs := o.QueryTable(&models.Blob{})
	qs = qs.Filter("Digest", digest)
	b := []*models.Blob{}
	_, err := qs.All(&b)
	if err != nil {
		return nil, fmt.Errorf("failed to get blob for digest %s, error: %v", digest, err)
	}
	if len(b) == 0 {
		log.Infof("No blob found for digest %s, returning empty.", digest)
		return &models.Blob{}, nil
	} else if len(b) > 1 {
		log.Infof("Multiple blob found for digest %s", digest)
		return &models.Blob{}, fmt.Errorf("Multiple blob found for digest %s", digest)
	}
	return b[0], nil
}

// DeleteBlob ...
func DeleteBlob(digest string) error {
	o := GetOrmer()
	_, err := o.QueryTable("blob").Filter("digest", digest).Delete()
	return err
}
