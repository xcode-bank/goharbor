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
	"testing"

	"github.com/goharbor/harbor/src/common/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddArtifact(t *testing.T) {
	af := &models.Artifact{
		PID:    1,
		Repo:   "hello-world",
		Tag:    "latest",
		Digest: "1234abcd",
		Kind:   "image",
	}

	// add
	id, err := AddArtifact(af)
	require.Nil(t, err)
	af.ID = id
	assert.Equal(t, id, int64(1))

}

func TestDeleteArtifact(t *testing.T) {
	af := &models.Artifact{
		PID:    1,
		Repo:   "hello-world",
		Tag:    "v1.0",
		Digest: "1234abcd",
		Kind:   "image",
	}
	// add
	id, err := AddArtifact(af)
	require.Nil(t, err)

	// delete
	err = DeleteArtifact(id)
	require.Nil(t, err)
}

func TestDeleteArtifactByDigest(t *testing.T) {
	af := &models.Artifact{
		PID:    1,
		Repo:   "hello-world",
		Tag:    "v1.1",
		Digest: "TestDeleteArtifactByDigest",
		Kind:   "image",
	}
	// add
	_, err := AddArtifact(af)
	require.Nil(t, err)

	// delete
	err = DeleteArtifactByDigest(af.Digest)
	require.Nil(t, err)
}

func TestDeleteArtifactByTag(t *testing.T) {
	af := &models.Artifact{
		PID:    1,
		Repo:   "hello-world",
		Tag:    "v1.2",
		Digest: "TestDeleteArtifactByTag",
		Kind:   "image",
	}
	// add
	_, err := AddArtifact(af)
	require.Nil(t, err)

	// delete
	err = DeleteByTag(1, "hello-world", "v1.2")
	require.Nil(t, err)
}
