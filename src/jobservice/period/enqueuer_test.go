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
package period

import (
	"context"
	"fmt"
	"github.com/goharbor/harbor/src/jobservice/common/rds"
	"github.com/goharbor/harbor/src/jobservice/common/utils"
	"github.com/goharbor/harbor/src/jobservice/env"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/jobservice/lcm"
	"github.com/goharbor/harbor/src/jobservice/tests"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

// EnqueuerTestSuite tests functions of enqueuer
type EnqueuerTestSuite struct {
	suite.Suite

	enqueuer  *enqueuer
	namespace string
	pool      *redis.Pool
	cancel    context.CancelFunc
}

// TestEnqueuerTestSuite is entry of go test
func TestEnqueuerTestSuite(t *testing.T) {
	suite.Run(t, new(EnqueuerTestSuite))
}

// SetupSuite prepares the test suite
func (suite *EnqueuerTestSuite) SetupSuite() {
	suite.namespace = tests.GiveMeTestNamespace()
	suite.pool = tests.GiveMeRedisPool()

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), utils.NodeID, "fake_node_ID"))
	suite.cancel = cancel

	envCtx := &env.Context{
		SystemContext: ctx,
		WG:            new(sync.WaitGroup),
	}

	lcmCtl := lcm.NewController(
		envCtx,
		suite.namespace,
		suite.pool,
		func(hookURL string, change *job.StatusChange) error { return nil },
	)
	suite.enqueuer = newEnqueuer(ctx, suite.namespace, suite.pool, lcmCtl)

	suite.prepare()
}

// TearDownSuite clears the test suite
func (suite *EnqueuerTestSuite) TearDownSuite() {
	suite.cancel()

	conn := suite.pool.Get()
	defer conn.Close()

	tests.ClearAll(suite.namespace, conn)
}

// TestEnqueuer tests enqueuer
func (suite *EnqueuerTestSuite) TestEnqueuer() {
	go func() {
		defer func() {
			suite.enqueuer.stopChan <- true
		}()

		<-time.After(1 * time.Second)

		key := rds.RedisKeyScheduled(suite.namespace)
		conn := suite.pool.Get()
		defer conn.Close()

		count, err := redis.Int(conn.Do("ZCARD", key))
		require.Nil(suite.T(), err, "count scheduled: nil error expected but got %s", err)
		assert.Condition(suite.T(), func() bool {
			return count > 0
		}, "count of scheduled jobs should be greater than 0 but got %d", count)
	}()

	err := suite.enqueuer.start()
	require.Nil(suite.T(), err, "enqueuer start: nil error expected but got %s", err)
}

func (suite *EnqueuerTestSuite) prepare() {
	now := time.Now()
	minute := now.Minute()

	coreSpec := fmt.Sprintf("30,50 %d * * * *", minute+2)

	// Prepare one
	p := &Policy{
		ID:       "fake_policy",
		JobName:  job.SampleJob,
		CronSpec: coreSpec,
	}
	rawData, err := p.Serialize()
	assert.Nil(suite.T(), err, "prepare data: nil error expected but got %s", err)
	key := rds.KeyPeriodicPolicy(suite.namespace)

	conn := suite.pool.Get()
	defer conn.Close()

	_, err = conn.Do("ZADD", key, time.Now().Unix(), rawData)
	assert.Nil(suite.T(), err, "prepare policy: nil error expected but got %s", err)
}
