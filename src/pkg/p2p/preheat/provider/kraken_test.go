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

package provider

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	cm "github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/pkg/p2p/preheat/models/provider"
	"github.com/goharbor/harbor/src/pkg/p2p/preheat/provider/auth"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// KrakenTestSuite is a test suite of testing Kraken driver.
type KrakenTestSuite struct {
	suite.Suite

	kraken *httptest.Server
	driver *KrakenDriver
}

// TestKraken is the entry method of running KrakenTestSuite.
func TestKraken(t *testing.T) {
	suite.Run(t, &KrakenTestSuite{})
}

// SetupSuite prepares the env for KrakenTestSuite.
func (suite *KrakenTestSuite) SetupSuite() {
	suite.kraken = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case krakenHealthPath:
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}

			w.WriteHeader(http.StatusOK)
		case krakenPreheatPath:
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}

			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			var payload = &cm.Notification{
				Events: []cm.Event{},
			}

			if err := json.Unmarshal(data, payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			if len(payload.Events) > 0 {
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))

	suite.kraken.StartTLS()

	suite.driver = &KrakenDriver{
		instance: &provider.Instance{
			ID:       2,
			Name:     "test-instance2",
			Vendor:   DriverKraken,
			Endpoint: suite.kraken.URL,
			AuthMode: auth.AuthModeNone,
			Enabled:  true,
			Default:  true,
			Insecure: true,
			Status:   DriverStatusHealthy,
		},
		digestFetcher: func(repoName, tag string) (s string, e error) {
			return "image@digest", nil
		},
	}
}

// TearDownSuite clears the env for KrakenTestSuite.
func (suite *KrakenTestSuite) TearDownSuite() {
	suite.kraken.Close()
}

// TestSelf tests Self method.
func (suite *KrakenTestSuite) TestSelf() {
	m := suite.driver.Self()
	suite.Equal(DriverKraken, m.ID, "self metadata")
}

// TestGetHealth tests GetHealth method.
func (suite *KrakenTestSuite) TestGetHealth() {
	st, err := suite.driver.GetHealth()
	require.NoError(suite.T(), err, "get health")
	suite.Equal(DriverStatusHealthy, st.Status, "healthy status")
}

// TestPreheat tests Preheat method.
func (suite *KrakenTestSuite) TestPreheat() {
	st, err := suite.driver.Preheat(&PreheatImage{
		Type:      "image",
		ImageName: "busybox",
		Tag:       "latest",
		URL:       "https://harbor.com",
	})
	require.NoError(suite.T(), err, "preheat image")
	suite.Equal(provider.PreheatingStatusSuccess, st.Status, "preheat image result")
	suite.NotEmptyf(st.FinishTime, "finish time")
}

// TestCheckProgress tests CheckProgress method.
func (suite *KrakenTestSuite) TestCheckProgress() {
	st, err := suite.driver.CheckProgress("kraken-id")
	require.NoError(suite.T(), err, "get preheat status")
	suite.Equal(provider.PreheatingStatusSuccess, st.Status, "preheat status")
	suite.NotEmptyf(st.FinishTime, "finish time")
}
