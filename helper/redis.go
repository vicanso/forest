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
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"

	"github.com/vicanso/forest/config"
)

var (
	redisClient *redis.Client
)

type (
	contextKey int
	// redisStats redis stats
	redisStats struct {
		Slow time.Duration
	}
)

const (
	// 记录命令开始时间
	startedAtKey contextKey = iota
)

func (rs *redisStats) logSlowOrError(ctx context.Context, cmd, message string) {
	t := ctx.Value(startedAtKey).(*time.Time)
	d := time.Since(*t)
	if d > rs.Slow || message != "" {
		// TODO 写入influxdb
		logger.Info("redis process slow or error",
			zap.String("cmd", cmd),
			zap.String("use", d.String()),
			zap.String("message", message),
		)
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
