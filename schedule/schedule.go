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

package schedule

import (
	"context"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/email"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/request"
	routerconcurrency "github.com/vicanso/forest/router_concurrency"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/go-performance"
)

type (
	taskFn      func() error
	statsTaskFn func() map[string]any
)

const logCategory = "schedule"

func init() {
	c := cron.New()
	_, _ = c.AddFunc("@every 1m", redisPing)
	_, _ = c.AddFunc("@every 1m", entPing)
	_, _ = c.AddFunc("@every 1m", influxdbPing)
	_, _ = c.AddFunc("@every 1m", configRefresh)
	_, _ = c.AddFunc("@every 1m", redisStats)
	_, _ = c.AddFunc("@every 1m", entStats)
	_, _ = c.AddFunc("@every 1m", influxdbStats)
	_, _ = c.AddFunc("@every 30s", cpuUsageStats)
	_, _ = c.AddFunc("@every 10s", performanceStats)
	_, _ = c.AddFunc("@every 1m", httpInstanceStats)
	_, _ = c.AddFunc("@every 1m", routerConcurrencyStats)
	// 如果是开发环境，则不执行定时任务
	if util.IsDevelopment() {
		return
	}
	c.Start()
}

func doTask(desc string, fn taskFn) {
	startedAt := time.Now()
	err := fn()
	use := time.Since(startedAt)
	if err != nil {
		log.Error(context.Background()).
			Str("category", logCategory).
			Dur("use", use).
			Err(err).
			Msg(desc + " fail")
		email.AlarmError(context.Background(), desc+" fail, "+err.Error())
	} else {
		log.Info(context.Background()).
			Str("category", logCategory).
			Dur("use", use).
			Msg(desc + " success")
	}
}

func doStatsTask(desc string, fn statsTaskFn) {
	startedAt := time.Now()
	stats := fn()
	log.Info(context.Background()).
		Str("category", logCategory).
		Dur("use", time.Since(startedAt)).
		Dict("stats", zerolog.Dict().Fields(stats)).
		Msg("")
}

func redisPing() {
	doTask("redis ping", helper.RedisPing)
}

func configRefresh() {
	configSrv := new(service.ConfigurationSrv)
	doTask("config refresh", func() error {
		return configSrv.Refresh(context.Background())
	})
}

func redisStats() {
	doStatsTask("redis stats", func() map[string]any {
		// 统计中除了redis数据库的统计，还有当前实例的统计指标，因此所有实例都会写入统计
		stats := helper.RedisStats()
		helper.GetInfluxDB().Write(cs.MeasurementRedisStats, nil, stats)
		return stats
	})
}

func entPing() {
	doTask("ent ping", helper.EntPing)
}

// entStats ent的性能统计
func entStats() {
	doStatsTask("ent stats", func() map[string]any {
		stats := helper.EntGetStats()
		helper.GetInfluxDB().Write(cs.MeasurementEntStats, nil, stats)
		return stats
	})
}

// cpuUsageStats cpu使用率
func cpuUsageStats() {
	doTask("update cpu usage", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		return service.UpdateCPUUsage(ctx)
	})
}

// influxdbStats influxdb统计
func influxdbStats() {
	doStatsTask("influxdb stats", func() map[string]any {
		db := helper.GetInfluxDB()
		writeCount := db.GetAndResetWriteCount()
		writingCount := db.GetWritingCount()
		fields := map[string]any{
			cs.FieldProcessing: writingCount,
			cs.FieldCount:      writeCount,
		}
		db.Write(cs.MeasurementInfluxdbStats, nil, fields)
		return fields
	})
}

// 上一次的性能指标
var prevPerformance *performance.Performance

// performanceStats 系统性能
func performanceStats() {
	doStatsTask("performance stats", func() map[string]any {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		data := service.GetPerformance(ctx)
		fields := data.Performance.ToMap(prevPerformance)

		fields[cs.FieldProcessing] = int(data.Concurrency)
		fields[cs.FieldConnProcessing] = int(data.HTTPServerConnStats.ConnProcessing)
		fields[cs.FieldConnAlive] = int(data.HTTPServerConnStats.ConnAlive)
		fields[cs.FieldConnCreatedCount] = int(data.HTTPServerConnStats.ConnCreatedCount)
		// 网络相关
		if data.ConnStat != nil {
			count := make(map[string]string)
			for k, v := range data.ConnStat.Status {
				count[k] = strconv.Itoa(v)
			}
			for k, v := range data.ConnStat.RemoteAddr {
				count[k] = strconv.Itoa(v)
			}
			log.Info(ctx).
				Str("category", "connStat").
				Dict("count", log.Struct(count)).
				Msg("")
		}

		// open files的统计
		if data.OpenFilesStats != nil &&
			len(data.OpenFilesStats.OpenFiles) != 0 {
			log.Info(ctx).
				Str("category", "openFiles").
				Dict("stat", log.Struct(data.OpenFilesStats)).
				Msg("")
		}
		prevPerformance = data.Performance

		helper.GetInfluxDB().Write(cs.MeasurementPerformance, nil, fields)
		return fields
	})
}

// httpInstanceStats http instance stats
func httpInstanceStats() {
	doStatsTask("http instance stats", func() map[string]any {
		fields := make(map[string]any)
		statsList := request.GetHTTPStats()
		for _, stats := range statsList {
			helper.GetInfluxDB().Write(cs.MeasurementHTTPInstanceStats, map[string]string{
				cs.TagService: stats.Name,
			}, map[string]any{
				cs.FieldMaxConcurrency: stats.MaxConcurrency,
				cs.FieldProcessing:     stats.Concurrency,
			})
			fields[stats.Name+":"+cs.FieldMaxConcurrency] = stats.MaxConcurrency
			fields[stats.Name+":"+cs.FieldProcessing] = stats.Concurrency
		}
		return fields
	})
}

// influxdbPing influxdb ping
func influxdbPing() {
	doTask("influxdb ping", helper.GetInfluxDB().Health)
}

// routerConcurrencyStats router concurrency stats
func routerConcurrencyStats() {
	doStatsTask("router concurrency stats", func() map[string]any {
		result := routerconcurrency.GetLimiter().GetStats()
		fields := make(map[string]any)

		influxSrv := helper.GetInfluxDB()
		for key, value := range result {
			// 如果并发为0，则不记录
			if value == 0 {
				continue
			}
			fields[key] = value
			influxSrv.Write(cs.MeasurementRouterConcurrency, map[string]string{
				cs.TagRoute: key,
			}, map[string]any{
				cs.FieldCount: int(value),
			})
		}
		return fields
	})
}
