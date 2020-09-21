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
	"database/sql"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/ent"
	"github.com/vicanso/forest/ent/hook"
	"go.uber.org/zap"
)

var (
	client *ent.Client
	driver *entsql.Driver

	initSchemaAndHookOnce sync.Once

	// entProcessing ent中正在处理的请求
	entProcessing uint32
)

func init() {
	postgresConfig := config.GetPostgresConfig()
	c, err := open(postgresConfig.URI)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	if err := c.Schema.Create(ctx); err != nil {
		panic(err)
	}
	client = c
}

// open new connection
func open(databaseUrl string) (*ent.Client, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}

	// Create an ent.Driver from `db`.
	driver = entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(driver)), nil
}

// EntGetClient get ent client
func EntGetClient() *ent.Client {
	return client
}

// EntPing ent driver ping
func EntPing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return driver.DB().PingContext(ctx)
}

// EntInitSchemaAndHook 初始化schema与hook函数
func EntInitSchemaAndHook() (err error) {
	// 只执行一次shcema初始化以及hook
	initSchemaAndHookOnce.Do(func() {
		err = client.Schema.Create(context.Background())
		if err != nil {
			return
		}
		// 禁止删除数据
		client.Use(hook.Reject(ent.OpDelete | ent.OpDeleteOne))
		// 数据库操作统计
		client.Use(func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				processing := atomic.AddUint32(&entProcessing, 1)
				defer atomic.AddUint32(&entProcessing, ^uint32(0))
				schemaType := m.Type()
				op := m.Op().String()

				startedAt := time.Now()
				result := 0
				message := ""
				value, err := next.Mutate(ctx, m)
				// 如果失败，则记录出错信息
				if err != nil {
					result = 1
					message = err.Error()
				}
				data := make(map[string]interface{})
				for _, name := range m.Fields() {
					if name == "updated_at" {
						continue
					}
					value, ok := m.Field(name)
					if !ok {
						continue
					}
					data[name] = value
				}

				d := time.Since(startedAt)
				logger.Info("ent stats",
					zap.String("schema", schemaType),
					zap.String("op", op),
					zap.Int("result", result),
					zap.Uint32("processing", processing),
					zap.String("use", d.String()),
					zap.Any("data", data),
					zap.String("message", message),
				)
				fields := map[string]interface{}{
					"processing": processing,
					"use":        int(d.Milliseconds()),
					"data":       data,
					"message":    message,
				}
				tags := map[string]string{
					"schema": schemaType,
					"op":     op,
					"result": strconv.Itoa(result),
				}
				GetInfluxSrv().Write(cs.MeasurementEntStats, fields, tags)
				return value, err
			})
		})
	})
	return
}
