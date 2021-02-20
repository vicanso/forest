// Copyright 2021 tree xie
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

// go routine tracer
// 只允许缓存用户账号与trace id等基本信息，
// 需要注意此缓存使用lru cache，因此有可能丢失，使用时仅用于日志等场景使用，
// 若逻辑上使用到的用户信息等，使用参数形式传递

package tracer

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/huandu/go-tls/g"
)

type TracerInfo struct {
	DeviceID string
	Account  string
	TraceID  string
}

var tracerInfoCache = mustNewTracerCache()

func mustNewTracerCache() *lru.Cache {
	// 设置缓存，根据系统的访问量调整，需要比request limit大
	l, err := lru.New(1024 * 10)
	if err != nil {
		panic(err)
	}
	return l
}

// GetTracerInfo 获取tracer信息
func GetTracerInfo() TracerInfo {
	p := g.G()
	if p == nil {
		return TracerInfo{}
	}
	value, ok := tracerInfoCache.Peek(p)
	if !ok {
		return TracerInfo{}
	}
	info, ok := value.(*TracerInfo)
	if !ok {
		return TracerInfo{}
	}
	return *info
}

// SetTracerInfo 设置tracer信息
func SetTracerInfo(info TracerInfo) {
	p := g.G()
	if p == nil {
		return
	}
	tracerInfoCache.Add(p, &info)
}
