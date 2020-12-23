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

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vicanso/forest/helper"
)

var (
	// redisSessionDefaultTimeout 默认超时
	redisSessionDefaultTimeout = 3 * time.Second
	// redisNoop 无操作的空函数
	redisNoop = func() error {
		return nil
	}

	// redisSrv 默认的redis服务
	redisSrv = new(Redis)
)

type (
	// RedisDone redis的done函数
	RedisDone func() error
	// Redis redis service
	Redis struct{}

	// RedisSessionStore session的redis缓存
	RedisSessionStore struct {
		Prefix string
	}
)

var getRedisClient = helper.RedisGetClient
var isRedisNilError = helper.RedisIsNilError

// Lock 将key锁定ttl的时间
func (srv *Redis) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return getRedisClient().SetNX(ctx, key, true, ttl).Result()
}

// Del 从缓存中删除key
func (srv *Redis) Del(ctx context.Context, key string) (err error) {
	_, err = getRedisClient().Del(ctx, key).Result()
	return
}

// LockWithDone 将key锁定ttl的时间，并提供done(删除)函数
func (srv *Redis) LockWithDone(ctx context.Context, key string, ttl time.Duration) (bool, RedisDone, error) {
	success, err := srv.Lock(ctx, key, ttl)
	// 如果lock失败，则返回no op 的done function
	if err != nil || !success {
		return false, redisNoop, err
	}
	done := func() error {
		err := srv.Del(ctx, key)
		return err
	}
	return true, done, nil
}

// IncWithTTL 增加key对应的值，并设置ttl
func (srv *Redis) IncWithTTL(ctx context.Context, key string, ttl time.Duration, value ...int64) (count int64, err error) {
	pipe := getRedisClient().TxPipeline()
	// 保证只有首次会设置ttl
	pipe.SetNX(ctx, key, 0, ttl)
	var incr *redis.IntCmd
	if len(value) != 0 {
		incr = pipe.IncrBy(ctx, key, value[0])
	} else {
		incr = pipe.Incr(ctx, key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}
	count = incr.Val()
	return
}

// Get 获取key的值
func (srv *Redis) Get(ctx context.Context, key string) (result []byte, err error) {
	result, err = getRedisClient().Get(ctx, key).Bytes()
	if err == redis.Nil {
		err = helper.ErrRedisNil
	}
	return
}

// GetIgnoreNilErr 获取key的值并忽略nil error
func (srv *Redis) GetIgnoreNilErr(ctx context.Context, key string) (result []byte, err error) {
	result, err = srv.Get(ctx, key)
	if helper.RedisIsNilError(err) {
		err = nil
	}
	return
}

// GetAndDel 获取key的值之后并删除它
func (srv *Redis) GetAndDel(ctx context.Context, key string) (result string, err error) {
	pipe := getRedisClient().TxPipeline()
	cmd := pipe.Get(ctx, key)
	pipe.Del(ctx, key)
	_, err = pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			err = helper.ErrRedisNil
		}
		return
	}
	result = cmd.Val()
	return
}

// Set 设置key的值并添加ttl
func (srv *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	getRedisClient().Set(ctx, key, value, ttl)
	return
}

// GetStruct 获取缓存并转换为struct
func (srv *Redis) GetStruct(ctx context.Context, key string, value interface{}) (err error) {
	result, err := srv.Get(ctx, key)
	if err != nil {
		return
	}
	err = json.Unmarshal(result, value)
	return
}

// SetStruct 将struct转换为json后保存并设置ttl
func (srv *Redis) SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	buf, err := json.Marshal(value)
	if err != nil {
		return
	}
	return srv.Set(ctx, key, buf, ttl)
}

// GetStruct 获取缓存解压后并转换为struct
func (srv *Redis) GetStructSnappy(ctx context.Context, key string, value interface{}) (err error) {
	result, err := srv.Get(ctx, key)
	if err != nil {
		return
	}
	result, err = helper.SnappyDecode(result)
	if err != nil {
		return
	}
	err = json.Unmarshal(result, value)
	return
}

// SetStructSnappy 将struct转换为json，压缩后保存，并设置ttl
// 需要注意，如果不是数据量较大（如10KB以上的，不建议使用）
func (srv *Redis) SetStructSnappy(ctx context.Context, key string, value interface{}, ttl time.Duration) (err error) {
	buf, err := json.Marshal(value)
	if err != nil {
		return
	}
	buf = helper.SnappyEncode(buf)
	return srv.Set(ctx, key, buf, ttl)
}

func (rs *RedisSessionStore) getKey(key string) string {
	return rs.Prefix + key
}

// Get 从redis中获取缓存的session
func (rs *RedisSessionStore) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisSessionDefaultTimeout)
	defer cancel()
	result, err := redisSrv.Get(ctx, rs.getKey(key))
	if isRedisNilError(err) {
		return nil, nil
	}
	return []byte(result), err
}

// Set 设置session至redis中
func (rs *RedisSessionStore) Set(key string, data []byte, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisSessionDefaultTimeout)
	defer cancel()
	return redisSrv.Set(ctx, rs.getKey(key), data, ttl)
}

// Destroy 从redis中删除session
func (rs *RedisSessionStore) Destroy(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisSessionDefaultTimeout)
	defer cancel()
	return redisSrv.Del(ctx, rs.getKey(key))
}
