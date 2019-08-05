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

	"github.com/go-redis/redis"
)

var (
	redisNoop = func() error {
		return nil
	}
)

type (
	// RedisDone redis done function
	RedisDone func() error
)

// RedisPing redis ping
func RedisPing() (err error) {
	_, err = redisGetClient().Ping().Result()
	return
}

// RedisLock lock the key for ttl seconds
func RedisLock(key string, ttl time.Duration) (bool, error) {
	return redisGetClient().SetNX(key, true, ttl).Result()
}

// RedisLockWithDone lock the key for ttl, and with done function
func RedisLockWithDone(key string, ttl time.Duration) (bool, RedisDone, error) {
	success, err := RedisLock(key, ttl)
	// 如果lock失败，则返回no op 的done function
	if err != nil || !success {
		return false, redisNoop, err
	}
	done := func() error {
		_, err := redisGetClient().Del(key).Result()
		return err
	}
	return true, done, nil
}

// RedisIncWithTTL inc value with ttl
func RedisIncWithTTL(key string, ttl time.Duration) (count int64, err error) {
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

// RedisGet get value
func RedisGet(key string) (result string, err error) {
	result, err = redisGetClient().Get(key).Result()
	// key不存在则不返回出错
	if err == redis.Nil {
		err = nil
	}
	return
}
