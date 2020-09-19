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
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vicanso/elton"
	"go.uber.org/zap"
)

type (
	// RouterConfig 路由配置信息
	RouterConfig struct {
		Route      string `json:"route,omitempty"`
		Method     string `json:"method,omitempty"`
		Status     int    `json:"status,omitempty"`
		CotentType string `json:"cotentType,omitempty"`
		Response   string `json:"response,omitempty"`
		// Delay 延时，单位秒
		Delay int    `json:"delay,omitempty"`
		URL   string `json:"url,omitempty"`
	}
	// RouterConcurrency 路由并发配置
	RouterConcurrency struct {
		Route   string `json:"route,omitempty"`
		Method  string `json:"method,omitempty"`
		Max     uint32 `json:"max,omitempty"`
		Current uint32 `json:"current,omitempty"`
	}
	// RCLimiter 路由请求限制
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

// IncConcurrency 当前路由处理数+1
func (l *RCLimiter) IncConcurrency(key string) (current uint32, max uint32) {
	r, ok := l.m[key]
	if !ok {
		return
	}
	current = atomic.AddUint32(&r.Current, 1)
	max = r.Max
	return
}

// DecConcurrency 当前路由处理数-1
func (l *RCLimiter) DecConcurrency(key string) {
	r, ok := l.m[key]
	if !ok {
		return
	}
	atomic.AddUint32(&r.Current, ^uint32(0))
}

// GetConcurrency 获取当前路由处理数
func (l *RCLimiter) GetConcurrency(key string) uint32 {
	r, ok := l.m[key]
	if !ok {
		return 0
	}
	return atomic.LoadUint32(&r.Current)
}

// 更新router config配置
func updateRouterConfigs(configs []string) {
	result := make(map[string]*RouterConfig)
	for _, item := range configs {
		v := &RouterConfig{}
		err := json.Unmarshal([]byte(item), v)
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

// RouterGetConfig 获取路由配置
func RouterGetConfig(method, route string) *RouterConfig {
	routerMutex.RLock()
	defer routerMutex.RUnlock()
	return currentRouterConfigs[method+route]
}

// InitRouterConcurrencyLimiter 初始路由并发限制
func InitRouterConcurrencyLimiter(routers []*elton.RouterInfo) {
	m := make(map[string]*RouterConcurrency)
	for _, item := range routers {
		m[item.Method+" "+item.Path] = &RouterConcurrency{}
	}
	rcLimiter.m = m
}

// GetRouterConcurrencyLimiter 获取路由并发限制器
func GetRouterConcurrencyLimiter() *RCLimiter {
	return rcLimiter
}

// ResetRouterConcurrency 重置路由并发数
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
