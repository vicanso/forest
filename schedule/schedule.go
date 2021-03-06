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
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/request"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

type (
	taskFn      func() error
	statsTaskFn func() map[string]interface{}
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
	_, _ = c.AddFunc("@every 30s", cpuUsageStats)
	_, _ = c.AddFunc("@every 1m", performanceStats)
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
		log.Default().Error().
			Str("category", logCategory).
			Dur("use", use).
			Err(err).
			Msg(desc + " fail")
		service.AlarmError(desc + " fail, " + err.Error())
	} else {
		log.Default().Info().
			Str("category", logCategory).
			Dur("use", use).
			Msg(desc + " success")
	}
}

func doStatsTask(desc string, fn statsTaskFn) {
	startedAt := time.Now()
	stats := fn()
	log.Default().Info().
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
	doTask("config refresh", configSrv.Refresh)
}

func redisStats() {
	doStatsTask("redis stats", func() map[string]interface{} {
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
	doStatsTask("ent stats", func() map[string]interface{} {
		stats := helper.EntGetStats()
		helper.GetInfluxDB().Write(cs.MeasurementEntStats, nil, stats)
		return stats
	})
}

// cpuUsageStats cpu使用率
func cpuUsageStats() {
	doTask("update cpu usage", service.UpdateCPUUsage)
}

// prevMemFrees 上一次 memory objects 释放的数量
var prevMemFrees uint64

// prevNumGC 上一次 gc 的次数
var prevNumGC uint32

// prevPauseTotal 上一次 pause 的总时长
var prevPauseTotal time.Duration

// performanceStats 系统性能
func performanceStats() {
	doStatsTask("performance stats", func() map[string]interface{} {
		data := service.GetPerformance()
		fields := map[string]interface{}{
			cs.FieldGoMaxProcs:       data.GoMaxProcs,
			cs.FieldProcessing:       int(data.Concurrency),
			cs.FieldThreadCount:      int(data.ThreadCount),
			cs.FieldMemSys:           data.MemSys,
			cs.FieldMemHeapSys:       data.MemHeapSys,
			cs.FieldMemHeapInuse:     data.MemHeapInuse,
			cs.FieldMemFrees:         int(data.MemFrees - prevMemFrees),
			cs.FieldRoutineCount:     data.RoutineCount,
			cs.FieldCpuUsage:         int(data.CPUUsage),
			cs.FieldNumGC:            int(data.NumGC - prevNumGC),
			cs.FieldPauseNS:          int((data.PauseTotalNs - prevPauseTotal).Milliseconds()),
			cs.FieldConnProcessing:   int(data.ConnProcessing),
			cs.FieldConnAlive:        int(data.ConnAlive),
			cs.FieldConnCreatedCount: int(data.ConnCreatedCount),
		}
		prevMemFrees = data.MemFrees
		prevNumGC = data.NumGC
		prevPauseTotal = data.PauseTotalNs

		helper.GetInfluxDB().Write(cs.MeasurementPerformance, nil, fields)
		return fields
	})
}

// httpInstanceStats http instance stats
func httpInstanceStats() {
	doStatsTask("http instance stats", func() map[string]interface{} {
		fields := make(map[string]interface{})
		statsList := request.GetHTTPStats()
		for _, stats := range statsList {
			helper.GetInfluxDB().Write(cs.MeasurementHTTPInstanceStats, map[string]string{
				cs.TagService: stats.Name,
			}, map[string]interface{}{
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
	doStatsTask("router concurrency stats", func() map[string]interface{} {
		result := service.GetRouterConcurrencyLimiter().GetStats()
		fields := make(map[string]interface{})

		influxSrv := helper.GetInfluxDB()
		for key, value := range result {
			// 如果并发为0，则不记录
			if value == 0 {
				continue
			}
			fields[key] = value
			influxSrv.Write(cs.MeasurementRouterConcurrency, map[string]string{
				cs.TagRoute: key,
			}, map[string]interface{}{
				cs.FieldCount: int(value),
			})
		}
		return fields
	})
}
