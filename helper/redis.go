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
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/hes"
	"go.uber.org/atomic"
)

var (
	defaultRedisClient, defaultRedisHook = mustNewRedisClient()

	// ErrRedisTooManyProcessing 处理请求太多时的出错
	ErrRedisTooManyProcessing = &hes.Error{
		Message:    "too many processing",
		StatusCode: http.StatusInternalServerError,
		Category:   "redis",
	}
)

type (

	// redisHook redis的hook配置
	redisHook struct {
		// 连接池大小
		poolSize int
		// 最大正在处理数量
		maxProcessing uint32
		// 慢请求阀值
		slow time.Duration
		// 正在处理数
		processing atomic.Uint32
		// pipe的正在处理数
		pipeProcessing atomic.Uint32
		// 总的处理请求数
		total atomic.Uint64
	}
)

func init() {
	redis.SetLogger(log.NewRedisLogger())
}
func mustNewRedisClient() (redis.UniversalClient, *redisHook) {
	redisConfig := config.MustGetRedisConfig()
	log.Default().Info().
		Strs("addr", redisConfig.Addrs).
		Msg("connect to redis")
	hook := &redisHook{
		slow:          redisConfig.Slow,
		maxProcessing: redisConfig.MaxProcessing,
	}
	opts := &redis.UniversalOptions{
		Addrs:            redisConfig.Addrs,
		Username:         redisConfig.Username,
		Password:         redisConfig.Password,
		SentinelPassword: redisConfig.Password,
		MasterName:       redisConfig.Master,
		PoolSize:         redisConfig.PoolSize,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Default().Info().Msg("redis new connection is established")
			GetInfluxDB().Write(cs.MeasurementRedisConn, nil, map[string]interface{}{
				cs.FieldCount: 1,
			})
			return nil
		},
		MinIdleConns: 2,
	}
	var c redis.UniversalClient
	// 需要对增加limiter，因此单独判断处理
	if opts.MasterName != "" {
		// TODO 确认有无可能增加limiter
		failoverOpts := opts.Failover()
		c = redis.NewFailoverClient(failoverOpts)
		hook.poolSize = failoverOpts.PoolSize
	} else if len(opts.Addrs) > 1 {
		clusterOpts := opts.Cluster()
		clusterOpts.NewClient = func(opt *redis.Options) *redis.Client {
			// 对每个client的增加limiter
			opt.Limiter = hook
			return redis.NewClient(opt)
		}
		c = redis.NewClusterClient(clusterOpts)
		hook.poolSize = clusterOpts.PoolSize
	} else {
		simpleOpts := opts.Simple()
		simpleOpts.Limiter = hook
		c = redis.NewClient(simpleOpts)
		hook.poolSize = simpleOpts.PoolSize
	}
	c.AddHook(hook)
	return c, hook
}

// 对于慢或出错请求输出日志并写入influxdb
func (rh *redisHook) logSlowOrError(ctx context.Context, cmd, err string) {
	t := ctx.Value(startedAtKey).(*time.Time)
	d := time.Since(*t)
	if d > rh.slow || err != "" {
		log.Default().Info().
			Str("category", "redisSlowOrErr").
			Str("cmd", cmd).
			Str("use", d.String()).
			Str("error", err).
			Msg("")
		tags := map[string]string{
			cs.TagOP: cmd,
		}
		fields := map[string]interface{}{
			cs.FieldUse:   int(d.Milliseconds()),
			cs.FieldError: err,
		}
		GetInfluxDB().Write(cs.MeasurementRedisOP, tags, fields)
	}
}

// BeforeProcess redis处理命令前的hook函数
func (rh *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, startedAtKey, &t)
	rh.processing.Inc()
	rh.total.Inc()
	return ctx, nil
}

// AfterProcess redis处理命令后的hook函数
func (rh *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	// allow返回error时也触发
	message := ""
	err := cmd.Err()
	if err != nil {
		message = err.Error()
	}
	rh.logSlowOrError(ctx, cmd.FullName(), message)
	rh.processing.Dec()
	if log.DebugEnabled() {
		// 由于redis是较频繁的操作
		// 由于cmd string的执行也有耗时，因此判断是否启用debug再输出
		log.Default().Debug().Msg(cmd.String())
	}
	return nil
}

// BeforeProcessPipeline redis pipeline命令前的hook函数
func (rh *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	// allow返回error时也触发
	t := time.Now()
	ctx = context.WithValue(ctx, startedAtKey, &t)
	rh.pipeProcessing.Inc()
	rh.total.Inc()
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
	rh.pipeProcessing.Dec()
	return nil
}

// getProcessingAndTotal 获取正在处理中的请求与总请求量
func (rh *redisHook) getProcessingAndTotal() (uint32, uint32, uint64) {
	processing := rh.processing.Load()
	pipeProcessing := rh.pipeProcessing.Load()
	total := rh.total.Load()
	return processing, pipeProcessing, total
}

// Allow 是否允许继续执行redis
func (rh *redisHook) Allow() error {
	// 如果处理请求量超出，则不允许继续请求
	if rh.processing.Load()+rh.pipeProcessing.Load() > rh.maxProcessing {
		return ErrRedisTooManyProcessing
	}
	return nil
}

// ReportResult 记录结果
func (*redisHook) ReportResult(result error) {
	// 仅是调用redis完成时触发
	// 如allow返回出错的不会触发
	if result != nil && !RedisIsNilError(result) {
		log.Default().Error().
			Str("category", "redisProcessFail").
			Err(result).
			Msg("")
		GetInfluxDB().Write(cs.MeasurementRedisError, nil, map[string]interface{}{
			cs.FieldError: result.Error(),
		})
	}
}

// RedisGetClient 获取redis client
func RedisGetClient() redis.UniversalClient {
	return defaultRedisClient
}

// RedisIsNilError 判断是否redis的nil error
func RedisIsNilError(err error) bool {
	return err == redis.Nil
}

// RedisStats 获取redis的性能统计
func RedisStats() map[string]interface{} {
	stats := RedisGetClient().PoolStats()
	processing, pipeProcessing, total := defaultRedisHook.getProcessingAndTotal()
	return map[string]interface{}{
		cs.FieldHits:          int(stats.Hits),
		cs.FieldMisses:        int(stats.Misses),
		cs.FieldTimeouts:      int(stats.Timeouts),
		cs.FieldTotalConns:    int(stats.TotalConns),
		cs.FieldIdleConns:     int(stats.IdleConns),
		cs.FieldStaleConns:    int(stats.StaleConns),
		cs.FieldProcessing:    int(processing),
		cs.FilePipeProcessing: int(pipeProcessing),
		cs.FieldTotal:         int(total),
		cs.FieldPoolSize:      defaultRedisHook.poolSize,
	}
}

// RedisPing ping操作
func RedisPing() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = RedisGetClient().Ping(ctx).Result()
	return
}
