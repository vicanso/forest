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

package helper

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
	redisNoop   = func() error {
		return nil
	}
	errRedisNil = hes.New("key is not exists or expired")
	redisSrv    = new(Redis)
)

type (

	// redisStats redis stats
	redisStats struct {
		Slow time.Duration
	}
)

type (
	// RedisDone redis done function
	RedisDone func() error
	// Redis redis service
	Redis struct{}

	// RedisSessionStore redis session store
	RedisSessionStore struct {
		Prefix string
	}
)

func (rs *redisStats) logSlowOrError(ctx context.Context, cmd, err string) {
	t := ctx.Value(startedAtKey).(*time.Time)
	d := time.Since(*t)
	if d > rs.Slow || err != "" {
		logger.Info("redis process slow or error",
			zap.String("cmd", cmd),
			zap.String("use", d.String()),
			zap.String("error", err),
		)
		tags := map[string]string{
			"cmd": cmd,
		}
		fields := map[string]interface{}{
			"use":   d.Milliseconds(),
			"error": err,
		}
		GetInfluxSrv().Write(cs.MeasurementRedis, fields, tags)
	}
}

// BeforeProcess before process
func (rs *redisStats) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, startedAtKey, &t)
	return ctx, nil
}

// AfterProcess after process
func (rs *redisStats) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	message := ""
	err := cmd.Err()
	if err != nil {
		message = err.Error()
	}
	rs.logSlowOrError(ctx, cmd.Name(), message)
	return nil
}

// BeforeProcessPipeline before process pipeline
func (rs *redisStats) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, startedAtKey, &t)
	return ctx, nil
}

// AfterProcessPipeline after process pipeline
func (rs *redisStats) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	cmdSb := new(strings.Builder)
	message := ""
	for index, cmd := range cmds {
		if index != 0 {
			cmdSb.WriteString(",")
		}
		cmdSb.WriteString(cmd.Name())
		err := cmd.Err()
		if err != nil {
			message += err.Error()
		}
	}
	rs.logSlowOrError(ctx, cmdSb.String(), message)
	return nil
}

func init() {
	options, err := config.GetRedisConfig()
	if err != nil {
		panic(err)
	}
	logger.Info("connect to redis",
		zap.String("addr", options.Addr),
		zap.Int("db", options.DB),
	)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
		OnConnect: func(_ *redis.Conn) error {
			logger.Info("redis new connection is established")
			return nil
		},
	})
	redisClient.AddHook(&redisStats{
		Slow: 300 * time.Millisecond,
	})
}

// RedisGetClient get redis client
func RedisGetClient() *redis.Client {
	return redisClient
}

// IsRedisNilError is redis nil errror
func IsRedisNilError(err error) bool {
	return err == errRedisNil
}

// RedisPing redis ping
func RedisPing() (err error) {
	_, err = redisClient.Ping().Result()
	return
}

// Lock lock the key for ttl seconds
func (srv *Redis) Lock(key string, ttl time.Duration) (bool, error) {
	return redisClient.SetNX(key, true, ttl).Result()
}

// Del del the key of redis
func (srv *Redis) Del(key string) (err error) {
	_, err = redisClient.Del(key).Result()
	return
}

// LockWithDone lock the key for ttl, and with done function
func (srv *Redis) LockWithDone(key string, ttl time.Duration) (bool, RedisDone, error) {
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
func (srv *Redis) IncWithTTL(key string, ttl time.Duration) (count int64, err error) {
	pipe := redisClient.TxPipeline()
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
func (srv *Redis) Get(key string) (result string, err error) {
	result, err = redisClient.Get(key).Result()
	if err == redis.Nil {
		err = errRedisNil
	}
	return
}

// GetIgnoreNilErr get value ignore nil error
func (srv *Redis) GetIgnoreNilErr(key string) (result string, err error) {
	result, err = srv.Get(key)
	if IsRedisNilError(err) {
		err = nil
	}
	return
}

// GetAndDel get value and del
func (srv *Redis) GetAndDel(key string) (result string, err error) {
	pipe := redisClient.TxPipeline()
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
func (srv *Redis) Set(key string, value interface{}, ttl time.Duration) (err error) {
	redisClient.Set(key, value, ttl)
	return
}

// GetStruct get struct
func (srv *Redis) GetStruct(key string, value interface{}) (err error) {
	result, err := srv.Get(key)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(result), value)
	return
}

// SetStruct redis set struct with ttl
func (srv *Redis) SetStruct(key string, value interface{}, ttl time.Duration) (err error) {
	buf, err := json.Marshal(value)
	if err != nil {
		return
	}
	return srv.Set(key, string(buf), ttl)
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
