package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/util"
)

func TestRedisSrv(t *testing.T) {
	assert := assert.New(t)
	srv := new(Redis)
	assert.NotNil(helper.RedisGetClient())
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
		assert.True(helper.RedisIsNilError(err))

		result, err := srv.GetIgnoreNilErr(todoCtx, key)
		assert.Nil(err)
		assert.Empty(result)

		value := "abc"
		err = srv.Set(todoCtx, key, value, time.Minute)
		assert.Nil((err))
		result, err = srv.Get(todoCtx, key)
		assert.Nil(err)
		assert.Equal(value, string(result))
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

	t.Run("get/set struct snappy", func(t *testing.T) {
		key := util.RandomString(10)
		type T struct {
			Name string `json:"name,omitempty"`
		}
		name := "Snappy 是一个 C++ 的用来压缩和解压缩的开发包。其目标不是最大限度压缩或者兼容其他压缩格式，而是旨在提供高速压缩速度和合理的压缩率。Snappy 比 zlib 更快，但文件相对要大 20% 到 100%。在 64位模式的 Core i7 处理器上，可达每秒 250~500兆的压缩速度。"
		err := srv.SetStructSnappy(todoCtx, key, &T{
			Name: name,
		}, time.Minute)
		assert.Nil(err)
		result := T{}
		err = srv.GetStructSnappy(todoCtx, key, &result)
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
