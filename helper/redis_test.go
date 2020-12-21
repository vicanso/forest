// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/forest/util"
)

func TestRedisStats(t *testing.T) {
	assert := assert.New(t)
	m := RedisStats()
	assert.NotNil(m)
}

func TestRedisPing(t *testing.T) {
	assert := assert.New(t)
	err := RedisPing()
	assert.Nil(err)
}

func TestRedisSrv(t *testing.T) {
	assert := assert.New(t)
	srv := new(Redis)
	assert.NotNil(RedisGetClient())
	todoCtx := context.TODO()
	t.Run("lock", func(t *testing.T) {
		key := util.RandomString(10)
		// 首次成功
		ok, err := srv.Lock(todoCtx, key, 5*time.Millisecond)
		assert.Nil(err)
		assert.True(ok)
		// 第二次失败
		ok, err = srv.Lock(todoCtx, key, 5*time.Millisecond)
		assert.Nil(err)
		assert.False(ok)
		// 第三次等待过期后成功
		time.Sleep(10 * time.Millisecond)
		ok, err = srv.Lock(todoCtx, key, 5*time.Millisecond)
		assert.Nil(err)
		assert.True(ok)
	})

	t.Run("lock with done", func(t *testing.T) {
		key := util.RandomString(10)
		// 首次成功
		ok, done, err := srv.LockWithDone(todoCtx, key, 5*time.Millisecond)
		assert.Nil(err)
		assert.True(ok)

		// 第二次失败
		ok, _, err = srv.LockWithDone(todoCtx, key, 5*time.Millisecond)
		assert.Nil(err)
		assert.False(ok)

		// 删除数据后第三次成功
		err = done()
		assert.Nil(err)
		ok, _, err = srv.LockWithDone(todoCtx, key, 5*time.Millisecond)
		assert.Nil(err)
		assert.True(ok)
	})

	t.Run("inc with ttl", func(t *testing.T) {
		key := util.RandomString(10)
		count, err := srv.IncWithTTL(todoCtx, key, time.Minute)
		assert.Nil(err)
		assert.Equal(int64(1), count)

		count, err = srv.IncWithTTL(todoCtx, key, time.Minute, 2)
		assert.Nil(err)
		assert.Equal(int64(3), count)
	})

	t.Run("get/set", func(t *testing.T) {
		key := util.RandomString(10)
		_, err := srv.Get(todoCtx, key)
		assert.True(RedisIsNilError(err))

		result, err := srv.GetIgnoreNilErr(todoCtx, key)
		assert.Nil(err)
		assert.Empty(result)

		value := "abc"
		err = srv.Set(todoCtx, key, value, time.Minute)
		assert.Nil((err))
		result, err = srv.Get(todoCtx, key)
		assert.Nil(err)
		assert.Equal(value, result)
	})

	t.Run("get/set struct", func(t *testing.T) {
		key := util.RandomString(10)
		type T struct {
			Name string `json:"name,omitempty"`
		}
		name := "abc"
		err := srv.SetStruct(todoCtx, key, &T{
			Name: name,
		}, time.Minute)
		assert.Nil(err)
		result := T{}
		err = srv.GetStruct(todoCtx, key, &result)
		assert.Nil(err)
		assert.Equal(name, result.Name)
	})
}

func TestRedisSessionStore(t *testing.T) {
	assert := assert.New(t)

	rs := RedisSessionStore{
		Prefix: "ss:",
	}
	key := util.RandomString(10)
	data := []byte("abc")
	err := rs.Set(key, data, time.Minute)
	assert.Nil(err)

	result, err := rs.Get(key)
	assert.Nil(err)
	assert.Equal(data, result)

	err = rs.Destroy(key)
	assert.Nil(err)

	result, err = rs.Get(key)
	assert.Nil(err)
	assert.Nil(result)
}
