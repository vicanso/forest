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

package service

import (
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/atomic"
)

type (
	// Performance 应用性能指标
	Performance struct {
		GoMaxProcs      int           `json:"goMaxProcs,omitempty"`
		Concurrency     int32         `json:"concurrency,omitempty"`
		ConnConcurrency int32         `json:"connConcurrency,omitempty"`
		ConnAlive       int32         `json:"connAlive,omitempty"`
		ThreadCount     int32         `json:"threadCount,omitempty"`
		MemSys          int           `json:"memSys,omitempty"`
		MemHeapSys      int           `json:"memHeapSys,omitempty"`
		MemHeapInuse    int           `json:"memHeapInuse,omitempty"`
		MemFrees        uint64        `json:"memFrees,omitempty"`
		RoutineCount    int           `json:"routineCount,omitempty"`
		CPUUsage        int32         `json:"cpuUsage,omitempty"`
		LastGC          time.Time     `json:"lastGC,omitempty"`
		NumGC           uint32        `json:"numGC,omitempty"`
		RecentPause     string        `json:"recentPause,omitempty"`
		RecentPauseNs   time.Duration `json:"recentPauseNs,omitempty"`
		PauseTotal      string        `json:"pauseTotal,omitempty"`
		PauseTotalNs    time.Duration `json:"pauseTotalNs,omitempty"`
		CPUBusy         string        `json:"cpuBusy,omitempty"`
		Uptime          string        `json:"uptime,omitempty"`
		PauseNs         [256]uint64   `json:"pauseNs,omitempty"`
	}
)

var (
	// 选择int32可根据如果值少于0，则处理有误
	concurrency atomic.Int32
	cpuUsage    atomic.Int32
	// 连接并发数
	connConcurrency atomic.Int32
	// 存活连接数
	connAlive atomic.Int32
)
var startedAt = time.Now()

var currentProcess *process.Process

func init() {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic(err)
	}
	currentProcess = p
	_ = UpdateCPUUsage()
}

// GetPerformance 获取应用性能指标
func GetPerformance() Performance {
	var mb uint64 = 1024 * 1024
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	seconds := int64(m.LastGC) / int64(time.Second)
	recentPauseNs := time.Duration(int64(m.PauseNs[(m.NumGC+255)%256]))
	pauseTotalNs := time.Duration(int64(m.PauseTotalNs))
	cpuTimes, _ := currentProcess.Times()
	cpuBusy := ""
	if cpuTimes != nil {
		busy := time.Duration(int64(cpuTimes.Total()-cpuTimes.Idle)) * time.Second
		cpuBusy = busy.String()
	}
	threadCount, _ := currentProcess.NumThreads()
	return Performance{
		GoMaxProcs:      runtime.GOMAXPROCS(0),
		Concurrency:     GetConcurrency(),
		ConnConcurrency: GetConnConcurrency(),
		ConnAlive:       GetConnAlive(),
		ThreadCount:     threadCount,
		MemSys:          int(m.Sys / mb),
		MemHeapSys:      int(m.HeapSys / mb),
		MemHeapInuse:    int(m.HeapInuse / mb),
		MemFrees:        m.Frees,
		RoutineCount:    runtime.NumGoroutine(),
		CPUUsage:        cpuUsage.Load(),
		LastGC:          time.Unix(seconds, 0),
		NumGC:           m.NumGC,
		RecentPause:     recentPauseNs.String(),
		RecentPauseNs:   recentPauseNs,
		PauseTotal:      pauseTotalNs.String(),
		PauseTotalNs:    pauseTotalNs,
		CPUBusy:         cpuBusy,
		Uptime:          time.Since(startedAt).String(),
		PauseNs:         m.PauseNs,
	}
}

// IncreaseConcurrency 当前并发请求+1
func IncreaseConcurrency() int32 {
	return concurrency.Inc()
}

// DecreaseConcurrency 当前并发请求-1
func DecreaseConcurrency() int32 {
	return concurrency.Dec()
}

// GetConcurrency 获取当前并发请求
func GetConcurrency() int32 {
	return concurrency.Load()
}

// IncreaseConnConcurrency 当前并发连接数+1
func IncreaseConnConcurrency() int32 {
	return connConcurrency.Inc()
}

// DecreaseConnConcurrency 当前并发连接数-1
func DecreaseConnConcurrency() int32 {
	return connConcurrency.Dec()
}

// GetConnConcurrency 获取当前并发连接数
func GetConnConcurrency() int32 {
	return connConcurrency.Load()
}

// IncreaseConnAlive 存活连接数+1
func IncreaseConnAlive() int32 {
	return connAlive.Inc()
}

// DecreaseConnAlive 存活连接数-1
func DecreaseConnAlive() int32 {
	return connAlive.Dec()
}

// GetConnAlive 获取当前存活连接数
func GetConnAlive() int32 {
	return connAlive.Load()
}

// UpdateCPUUsage 更新cpu使用率
func UpdateCPUUsage() error {
	usage, err := currentProcess.Percent(0)
	if err != nil {
		return err
	}
	cpuUsage.Store(int32(usage))
	return nil
}
