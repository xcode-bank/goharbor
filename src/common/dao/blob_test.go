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

package dao

import (
	"github.com/goharbor/harbor/src/common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddBlob(t *testing.T) {
	blob := &models.Blob{
		Digest:      "1234abcd",
		ContentType: "v2.blob",
		Size:        1523,
	}

	// add
	_, err := AddBlob(blob)
	require.Nil(t, err)
}

func TestGetBlob(t *testing.T) {
	blob := &models.Blob{
		Digest:      "12345abcde",
		ContentType: "v2.blob",
		Size:        453,
	}

	// add
	id, err := AddBlob(blob)
	require.Nil(t, err)
	blob.ID = id

	blob2, err := GetBlob("12345abcde")
	require.Nil(t, err)
	assert.Equal(t, blob.Digest, blob2.Digest)

}

func TestDeleteBlob(t *testing.T) {
	blob := &models.Blob{
		Digest:      "123456abcdef",
		ContentType: "v2.blob",
		Size:        4543,
	}
	id, err := AddBlob(blob)
	require.Nil(t, err)
	blob.ID = id
	err = DeleteBlob(blob.Digest)
	require.Nil(t, err)
}

func TestHasBlobInProject(t *testing.T) {
	af := &models.Artifact{
		PID:    1,
		Repo:   "TestHasBlobInProject",
		Tag:    "latest",
		Digest: "tttt",
		Kind:   "image",
	}

	// add
	_, err := AddArtifact(af)
	require.Nil(t, err)

	afnb1 := &models.ArtifactAndBlob{
		DigestAF:   "tttt",
		DigestBlob: "zzza",
	}
	afnb2 := &models.ArtifactAndBlob{
		DigestAF:   "tttt",
		DigestBlob: "zzzb",
	}
	afnb3 := &models.ArtifactAndBlob{
		DigestAF:   "tttt",
		DigestBlob: "zzzc",
	}

	var afnbs []*models.ArtifactAndBlob
	afnbs = append(afnbs, afnb1)
	afnbs = append(afnbs, afnb2)
	afnbs = append(afnbs, afnb3)

	// add
	err = AddArtifactNBlobs(afnbs)
	require.Nil(t, err)

	has, err := HasBlobInProject(1, "zzzb")
	require.Nil(t, err)
	assert.True(t, has)
}
