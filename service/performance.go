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
	"runtime"

	"go.uber.org/atomic"
)

type (
	// Performance 应用性能指标
	Performance struct {
		GoMaxProcs   int    `json:"goMaxProcs,omitempty"`
		Concurrency  uint32 `json:"concurrency,omitempty"`
		Sys          int    `json:"sys,omitempty"`
		HeapSys      int    `json:"heapSys,omitempty"`
		HeapInuse    int    `json:"heapInuse,omitempty"`
		RoutineCount int    `json:"routineCount,omitempty"`
	}
)

var (
	concurrency atomic.Uint32
)

// GetPerformance 获取应用性能指标
func GetPerformance() Performance {
	var mb uint64 = 1024 * 1024
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return Performance{
		GoMaxProcs:   runtime.GOMAXPROCS(0),
		Concurrency:  GetConcurrency(),
		Sys:          int(m.Sys / mb),
		HeapSys:      int(m.HeapSys / mb),
		HeapInuse:    int(m.HeapInuse / mb),
		RoutineCount: runtime.NumGoroutine(),
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
