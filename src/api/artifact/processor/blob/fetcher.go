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
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/goharbor/harbor/src/pkg/registry"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io/ioutil"
)

var (
	// Fcher is a global blob fetcher instance
	Fcher = NewFetcher()

	accept = []string{
		schema1.MediaTypeSignedManifest,
		schema2.MediaTypeManifest,
		v1.MediaTypeImageManifest,
		manifestlist.MediaTypeManifestList,
		v1.MediaTypeImageIndex,
	}
)

// TODO use the registry.Client directly? then the Fetcher can be deleted

// Fetcher fetches the content of blob
type Fetcher interface {
	// FetchManifest the content of manifest under the repository
	FetchManifest(repository, digest string) (mediaType string, content []byte, err error)
	// FetchLayer the content of layer under the repository
	FetchLayer(repository, digest string) (content []byte, err error)
}

// NewFetcher returns an instance of the default blob fetcher
func NewFetcher() Fetcher {
	return &fetcher{
		client: registry.Cli,
	}
}

type fetcher struct {
	client registry.Client
}

func (f *fetcher) FetchManifest(repository, digest string) (string, []byte, error) {
	// TODO read from cache first
	manifest, _, err := f.client.PullManifest(repository, digest)
	if err != nil {
		return "", nil, err
	}
	mediaType, payload, err := manifest.Payload()
	if err != nil {
		return "", nil, err
	}
	return mediaType, payload, err
}

func (f *fetcher) FetchLayer(repository, digest string) ([]byte, error) {
	// TODO read from cache first
	_, reader, err := f.client.PullBlob(repository, digest)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}
