// Copyright 2019 tree xie
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

package service

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/vicanso/hes"
)

var (
	redisNoop = func() error {
		return nil
	}
	errRedisNil = hes.New("key is not exists or expired")
)

type (
	// RedisDone redis done function
	RedisDone func() error
	// RedisSrv redis service
	RedisSrv struct{}

	// RedisSessionStore redis session store
	RedisSessionStore struct {
		Prefix string
	}
)

// IsRedisNilError is redis nil errror
func IsRedisNilError(err error) bool {
	return err == errRedisNil
}

// RedisPing redis ping
func RedisPing() (err error) {
	_, err = redisGetClient().Ping().Result()
	return
}

// Lock lock the key for ttl seconds
func (srv *RedisSrv) Lock(key string, ttl time.Duration) (bool, error) {
	return redisGetClient().SetNX(key, true, ttl).Result()
}

// Del del the key of redis
func (srv *RedisSrv) Del(key string) (err error) {
	_, err = redisGetClient().Del(key).Result()
	return
}

// LockWithDone lock the key for ttl, and with done function
func (srv *RedisSrv) LockWithDone(key string, ttl time.Duration) (bool, RedisDone, error) {
	success, err := srv.Lock(key, ttl)
	// 如果lock失败，则返回no op 的done function
	if err != nil || !success {
		return false, redisNoop, err
	}
	done := func() error {
		err := srv.Del(key)
		return err
	}
	return true, done, nil
}

// IncWithTTL inc value with ttl
func (srv *RedisSrv) IncWithTTL(key string, ttl time.Duration) (count int64, err error) {
	pipe := redisGetClient().TxPipeline()
	// 保证只有首次会设置ttl
	pipe.SetNX(key, 0, ttl)
	incr := pipe.Incr(key)
	_, err = pipe.Exec()
	if err != nil {
		return
	}
	count = incr.Val()
	return
}

// Get get value
func (srv *RedisSrv) Get(key string) (result string, err error) {
	result, err = redisGetClient().Get(key).Result()
	if err == redis.Nil {
		err = errRedisNil
	}
	return
}

// GetIgnoreNilErr get value ignore nil error
func (srv *RedisSrv) GetIgnoreNilErr(key string) (result string, err error) {
	result, err = srv.Get(key)
	if IsRedisNilError(err) {
		err = nil
	}
	return
}

// GetAndDel get value and del
func (srv *RedisSrv) GetAndDel(key string) (result string, err error) {
	pipe := redisGetClient().TxPipeline()
	cmd := pipe.Get(key)
	pipe.Del(key)
	_, err = pipe.Exec()
	if err != nil {
		if err == redis.Nil {
			err = errRedisNil
		}
		return
	}
	result = cmd.Val()
	return
}

// Set redis set with ttl
func (srv *RedisSrv) Set(key string, value interface{}, ttl time.Duration) (err error) {
	redisGetClient().Set(key, value, ttl)
	return
}

// GetStruct get struct
func (srv *RedisSrv) GetStruct(key string, value interface{}) (err error) {
	result, err := srv.Get(key)
	if err != nil {
		return
	}
	err = fastestJSON.UnmarshalFromString(result, value)
	return
}

// SetStruct redis set struct with ttl
func (srv *RedisSrv) SetStruct(key string, value interface{}, ttl time.Duration) (err error) {
	str, err := fastestJSON.MarshalToString(value)
	if err != nil {
		return
	}
	return srv.Set(key, str, ttl)
}

func (rs *RedisSessionStore) getKey(key string) string {
	return rs.Prefix + key
}

// Get get the session from redis
func (rs *RedisSessionStore) Get(key string) ([]byte, error) {
	result, err := redisSrv.Get(rs.getKey(key))
	if IsRedisNilError(err) {
		return nil, nil
	}
	return []byte(result), err
}

// Set set the session to redis
func (rs *RedisSessionStore) Set(key string, data []byte, ttl time.Duration) error {
	return redisSrv.Set(rs.getKey(key), data, ttl)
}

// Destroy remove the session from redis
func (rs *RedisSessionStore) Destroy(key string) error {
	return redisSrv.Del(rs.getKey(key))
}
