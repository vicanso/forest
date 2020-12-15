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
		GoMaxProcs    int           `json:"goMaxProcs,omitempty"`
		Concurrency   uint32        `json:"concurrency,omitempty"`
		ThreadCount   int32         `json:"threadCount,omitempty"`
		MemSys        int           `json:"memSys,omitempty"`
		MemHeapSys    int           `json:"memHeapSys,omitempty"`
		MemHeapInuse  int           `json:"memHeapInuse,omitempty"`
		MemFrees      uint64        `json:"memFrees,omitempty"`
		RoutineCount  int           `json:"routineCount,omitempty"`
		CPUUsage      uint32        `json:"cpuUsage,omitempty"`
		LastGC        time.Time     `json:"lastGC,omitempty"`
		NumGC         uint32        `json:"numGC,omitempty"`
		RecentPause   string        `json:"recentPause,omitempty"`
		RecentPauseNs time.Duration `json:"recentPauseNs,omitempty"`
		PauseTotal    string        `json:"pauseTotal,omitempty"`
		PauseTotalNs  time.Duration `json:"pauseTotalNs,omitempty"`
		CPUBusy       string        `json:"cpuBusy,omitempty"`
		Uptime        string        `json:"uptime,omitempty"`
		PauseNs       [256]uint64   `json:"pauseNs,omitempty"`
	}
)

var (
	concurrency atomic.Uint32
	cpuUsage    atomic.Uint32
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
		GoMaxProcs:    runtime.GOMAXPROCS(0),
		Concurrency:   GetConcurrency(),
		ThreadCount:   threadCount,
		MemSys:        int(m.Sys / mb),
		MemHeapSys:    int(m.HeapSys / mb),
		MemHeapInuse:  int(m.HeapInuse / mb),
		MemFrees:      m.Frees,
		RoutineCount:  runtime.NumGoroutine(),
		CPUUsage:      cpuUsage.Load(),
		LastGC:        time.Unix(seconds, 0),
		NumGC:         m.NumGC,
		RecentPause:   recentPauseNs.String(),
		RecentPauseNs: recentPauseNs,
		PauseTotal:    pauseTotalNs.String(),
		PauseTotalNs:  pauseTotalNs,
		CPUBusy:       cpuBusy,
		Uptime:        time.Since(startedAt).String(),
		PauseNs:       m.PauseNs,
	}
}

// IncreaseConcurrency 当前并发请求+1
func IncreaseConcurrency() uint32 {
	return concurrency.Inc()
}

// DecreaseConcurrency 当前并发请求-1
func DecreaseConcurrency() uint32 {
	return concurrency.Dec()
}

// GetConcurrency 获取当前并发请求
func GetConcurrency() uint32 {
	return concurrency.Load()
}

// UpdateCPUUsage 更新cpu使用率
func UpdateCPUUsage() error {
	usage, err := currentProcess.Percent(0)
	if err != nil {
		return err
	}
	cpuUsage.Store(uint32(usage))
	return nil
}
