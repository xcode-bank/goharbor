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

package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/goharbor/harbor/src/lib/cache"
	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
	cache cache.Cache
}

func (suite *CacheTestSuite) SetupSuite() {
	suite.cache, _ = cache.New("redis", cache.Expiration(time.Second*5))
}

func (suite *CacheTestSuite) TestContains() {
	key := "contains"
	suite.False(suite.cache.Contains(key))

	suite.cache.Save(key, "value")
	suite.True(suite.cache.Contains(key))

	suite.cache.Delete(key)
	suite.False(suite.cache.Contains(key))

	suite.cache.Save(key, "value", time.Second*5)
	suite.True(suite.cache.Contains(key))

	time.Sleep(time.Second * 8)
	suite.False(suite.cache.Contains(key))
}

func (suite *CacheTestSuite) TestDelete() {
	key := "delete"

	suite.cache.Save(key, "value")
	suite.True(suite.cache.Contains(key))

	suite.cache.Delete(key)
	suite.False(suite.cache.Contains(key))
}

func (suite *CacheTestSuite) TestFetch() {
	key := "fetch"

	suite.cache.Save(key, map[string]interface{}{"name": "harbor", "version": "1.10"})

	mp := map[string]interface{}{}
	suite.cache.Fetch(key, &mp)
	suite.Len(mp, 2)
	suite.Equal("harbor", mp["name"])
	suite.Equal("1.10", mp["version"])

	var str string
	suite.Error(suite.cache.Fetch(key, &str))
}

func (suite *CacheTestSuite) TestSave() {
	key := "save"

	{
		suite.cache.Save(key, "hello, save")

		var value string
		suite.cache.Fetch(key, &value)
		suite.Equal("hello, save", value)

		time.Sleep(time.Second * 8)

		value = ""
		suite.Error(suite.cache.Fetch(key, &value))
		suite.Equal("", value)
	}

	{
		suite.cache.Save(key, "hello, save", time.Second)

		time.Sleep(time.Second * 2)

		var value string
		suite.Error(suite.cache.Fetch(key, &value))
		suite.Equal("", value)
	}
}

func (suite *CacheTestSuite) TestPing() {
	suite.NoError(suite.cache.Ping())
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func BenchmarkCacheFetchParallel(b *testing.B) {
	key := "benchmark"
	cache, _ := cache.New("redis")
	cache.Save(key, "hello, benchmark")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var value string
			err := cache.Fetch(key, &value)
			if err != nil {
				fmt.Printf("failed, error %v\n", err)
			}
		}
	})
}

func BenchmarkCacheSaveParallel(b *testing.B) {
	key := "benchmark"
	cache, _ := cache.New("redis")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Save(key, "hello, benchmark")
		}
	})
}
