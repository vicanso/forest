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

package service

import (
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vicanso/elton"

	"go.uber.org/zap"
)

type (
	// RouterConfig route config
	RouterConfig struct {
		Route      string `json:"route"`
		Method     string `json:"method"`
		Status     int    `json:"status"`
		CotentType string `json:"cotentType"`
		Response   string `json:"response"`
		// Delay 延时，单位秒
		Delay int    `json:"delay"`
		URL   string `json:"url"`
	}
	// RouterConcurrency router concurrency
	RouterConcurrency struct {
		Route   string `json:"route"`
		Method  string `json:"method"`
		Max     uint32 `json:"max"`
		Current uint32 `json:"current"`
	}
	// RCLimiter
	RCLimiter struct {
		m map[string]*RouterConcurrency
	}
)

var (
	routerMutex          = new(sync.RWMutex)
	currentRouterConfigs map[string]*RouterConfig
	rcLimiter            *RCLimiter
)

func init() {
	rcLimiter = &RCLimiter{}
}

// IncConcurrency inc concurrency
func (l *RCLimiter) IncConcurrency(key string) (current uint32, max uint32) {
	r, ok := l.m[key]
	if !ok {
		return
	}
	current = atomic.AddUint32(&r.Current, 1)
	max = r.Max
	return
}

// DecConcurrency dec concurrency
func (l *RCLimiter) DecConcurrency(key string) {
	r, ok := l.m[key]
	if !ok {
		return
	}
	atomic.AddUint32(&r.Current, ^uint32(0))
}

// GetConcurrency get concurrency
func (l *RCLimiter) GetConcurrency(key string) uint32 {
	r, ok := l.m[key]
	if !ok {
		return 0
	}
	return atomic.LoadUint32(&r.Current)
}

// 更新router config配置
func updateRouterConfigs(configs []*Configuration) {
	result := make(map[string]*RouterConfig)
	for _, item := range configs {
		v := &RouterConfig{}
		err := json.Unmarshal([]byte(item.Data), v)
		if err != nil {
			logger.Error("router config is invalid",
				zap.Error(err),
			)
			AlarmError("router config is invalid:" + err.Error())
			continue
		}
		// 如果未配置Route或者method的则忽略
		if v.Route == "" || v.Method == "" {
			continue
		}
		result[v.Method+v.Route] = v
	}
	routerMutex.Lock()
	defer routerMutex.Unlock()
	currentRouterConfigs = result
}

// RouterGetConfig get router config
func RouterGetConfig(method, route string) *RouterConfig {
	routerMutex.RLock()
	defer routerMutex.RUnlock()
	return currentRouterConfigs[method+route]
}

// InitRouterConcurrencyLimiter init router concurrency limiter
func InitRouterConcurrencyLimiter(routers []*elton.RouterInfo) {
	m := make(map[string]*RouterConcurrency)
	for _, item := range routers {
		m[item.Method+" "+item.Path] = &RouterConcurrency{}
	}
	rcLimiter.m = m
}

// GetRouterConcurrencyLimiter get router concurrency limiter
func GetRouterConcurrencyLimiter() *RCLimiter {
	return rcLimiter
}

// ResetRouterConcurrency reset router councurrency
func ResetRouterConcurrency(arr []string) {
	concurrencyList := make([]*RouterConcurrency, 0)
	for _, str := range arr {
		v := &RouterConcurrency{}
		err := json.Unmarshal([]byte(str), v)
		if err != nil {
			logger.Error("router concurrency config is invalid",
				zap.Error(err),
			)
			AlarmError("router concurrency config is invalid:" + err.Error())
			continue
		}
		concurrencyList = append(concurrencyList, v)
	}
	for key, r := range rcLimiter.m {
		keys := strings.Split(key, " ")
		if len(keys) != 2 {
			continue
		}
		found := false
		for _, item := range concurrencyList {
			if item.Method == keys[0] && item.Route == keys[1] {
				found = true
				// 设置并发请求量
				atomic.StoreUint32(&r.Max, item.Max)
			}
		}
		// 如果未配置，则设置为限制0（无限制）
		if !found {
			atomic.StoreUint32(&r.Max, 0)
		}
	}
}
