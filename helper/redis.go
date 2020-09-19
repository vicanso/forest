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
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
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
	errRedisNil               = hes.New("key is not exists or expired")
	redisSrv                  = new(Redis)
	rh                        *redisHook
	ErrRedisTooManyProcessing = &hes.Error{
		Message:    "too many processing",
		StatusCode: http.StatusInternalServerError,
		Category:   "redis",
	}
)

type (

	// redisHook redis的hook配置
	redisHook struct {
		maxProcessing  uint32
		slow           time.Duration
		processing     uint32
		pipeProcessing uint32
		total          uint64
	}
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

// 对于慢或出错请求输出日志并写入influxdb
func (rh *redisHook) logSlowOrError(ctx context.Context, cmd, err string) {
	t := ctx.Value(startedAtKey).(*time.Time)
	d := time.Since(*t)
	if d > rh.slow || err != "" {
		logger.Info("redis process slow or error",
			zap.String("cmd", cmd),
			zap.String("use", d.String()),
			zap.String("error", err),
		)
		tags := map[string]string{
			"cmd": cmd,
		}
		fields := map[string]interface{}{
			"use":   int(d.Milliseconds()),
			"error": err,
		}
		GetInfluxSrv().Write(cs.MeasurementRedisStats, fields, tags)
	}
}

// BeforeProcess redis处理命令前的hook函数
func (rh *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, startedAtKey, &t)
	atomic.AddUint32(&rh.processing, 1)
	atomic.AddUint64(&rh.total, 1)
	return ctx, nil
}

// AfterProcess redis处理命令后的hook函数
func (rh *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	message := ""
	err := cmd.Err()
	if err != nil {
		message = err.Error()
	}
	rh.logSlowOrError(ctx, cmd.Name(), message)
	atomic.AddUint32(&rh.processing, ^uint32(0))
	return nil
}

// BeforeProcessPipeline redis pipeline命令前的hook函数
func (rh *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, startedAtKey, &t)
	atomic.AddUint32(&rh.pipeProcessing, 1)
	atomic.AddUint64(&rh.total, 1)
	return ctx, nil
}

// AfterProcessPipeline redis pipeline命令后的hook函数
func (rh *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
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
	rh.logSlowOrError(ctx, cmdSb.String(), message)
	atomic.AddUint32(&rh.pipeProcessing, ^uint32(0))
	return nil
}

// getProcessingAndTotal 获取正在处理中的请求与总请求量
func (rh *redisHook) getProcessingAndTotal() (uint32, uint32, uint64) {
	processing := atomic.LoadUint32(&rh.processing)
	pipeProcessing := atomic.LoadUint32(&rh.pipeProcessing)
	total := atomic.LoadUint64(&rh.total)
	return processing, pipeProcessing, total
}

// Allow 是否允许继续执行redis
func (rh *redisHook) Allow() error {
	// 如果处理请求量超出，则不允许继续请求
	if atomic.LoadUint32(&rh.processing) > rh.maxProcessing {
		return ErrRedisTooManyProcessing
	}
	return nil
}

// ReportResult 记录结果
func (*redisHook) ReportResult(result error) {
	// TODO 对于nil error另外统计
	if result != nil && !IsRedisNilError(result) {
		logger.Error("redis process fail",
			zap.Error(result),
		)
	}
}

func init() {
	redisConfig := config.GetRedisConfig()
	logger.Info("connect to redis",
		zap.String("addr", redisConfig.Addr),
		zap.Int("db", redisConfig.DB),
	)
	rh = &redisHook{
		slow:          redisConfig.Slow,
		maxProcessing: redisConfig.MaxProcessing,
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		Limiter:  rh,
	})
	redisClient.AddHook(rh)
}

// RedisGetClient 获取redis client
func RedisGetClient() *redis.Client {
	return redisClient
}

// IsRedisNilError 判断是否redis的nil error
func IsRedisNilError(err error) bool {
	return err == errRedisNil || err == redis.Nil
}

// RedisStats 获取redis的性能统计
func RedisStats() map[string]interface{} {
	stats := redisClient.PoolStats()
	processing, pipeProcessing, total := rh.getProcessingAndTotal()
	return map[string]interface{}{
		"hits":           stats.Hits,
		"missed":         stats.Misses,
		"timeouts":       stats.Timeouts,
		"totalConns":     stats.TotalConns,
		"idleConns":      stats.IdleConns,
		"staleConns":     stats.StaleConns,
		"processing":     processing,
		"pipeProcessing": pipeProcessing,
		"total":          total,
	}
}

// RedisPing ping操作
func RedisPing() (err error) {
	_, err = redisClient.Ping().Result()
	return
}

// Lock 将key锁定ttl的时间
func (srv *Redis) Lock(key string, ttl time.Duration) (bool, error) {
	return redisClient.SetNX(key, true, ttl).Result()
}

// Del 从缓存中删除key
func (srv *Redis) Del(key string) (err error) {
	_, err = redisClient.Del(key).Result()
	return
}

// LockWithDone 将key锁定ttl的时间，并提供done(删除)函数
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

// IncWithTTL 增加key对应的值，并设置ttl
func (srv *Redis) IncWithTTL(key string, ttl time.Duration, value ...int64) (count int64, err error) {
	pipe := redisClient.TxPipeline()
	// 保证只有首次会设置ttl
	pipe.SetNX(key, 0, ttl)
	var incr *redis.IntCmd
	if len(value) != 0 {
		incr = pipe.IncrBy(key, value[0])
	} else {
		incr = pipe.Incr(key)
	}
	_, err = pipe.Exec()
	if err != nil {
		return
	}
	count = incr.Val()
	return
}

// Get 获取key的值
func (srv *Redis) Get(key string) (result string, err error) {
	result, err = redisClient.Get(key).Result()
	if err == redis.Nil {
		err = errRedisNil
	}
	return
}

// GetIgnoreNilErr 获取key的值并忽略nil error
func (srv *Redis) GetIgnoreNilErr(key string) (result string, err error) {
	result, err = srv.Get(key)
	if IsRedisNilError(err) {
		err = nil
	}
	return
}

// GetAndDel 获取key的值之后并删除它
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

// Set 设置key的值并添加ttl
func (srv *Redis) Set(key string, value interface{}, ttl time.Duration) (err error) {
	redisClient.Set(key, value, ttl)
	return
}

// GetStruct 获取缓存并转换为struct
func (srv *Redis) GetStruct(key string, value interface{}) (err error) {
	result, err := srv.Get(key)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(result), value)
	return
}

// SetStruct 将struct转换为字符串后保存并设置ttl
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

// Get 从redis中获取缓存的session
func (rs *RedisSessionStore) Get(key string) ([]byte, error) {
	result, err := redisSrv.Get(rs.getKey(key))
	if IsRedisNilError(err) {
		return nil, nil
	}
	return []byte(result), err
}

// Set 设置session至redis中
func (rs *RedisSessionStore) Set(key string, data []byte, ttl time.Duration) error {
	return redisSrv.Set(rs.getKey(key), data, ttl)
}

// Destroy 从redis中删除session
func (rs *RedisSessionStore) Destroy(key string) error {
	return redisSrv.Del(rs.getKey(key))
}
