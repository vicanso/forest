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

// 多层缓存，缓存时将缓存写入redis与lru中，
// 读取缓存时，优先读取lru，如果lru缓存不存在，则读取redis，
// 二级缓存主要用于多次使用(同一次请求的处理函数中）且数据量较少的缓存，
// 建议使用时针对lru size设置较小的值（如1000），避免过多占用应用内存。
// 使用时建议添加prefix方便区别缓存，避免冲突

package cache

import (
	"context"
	"encoding/json"
	"time"

	lruttl "github.com/vicanso/lru-ttl"
)

type multiCache struct {
	ttl    time.Duration
	lru    *lruttl.Cache
	prefix string
}

// NewMultiCache create a new multi cache
func NewMultiCache(lruSize int, ttl time.Duration) *multiCache {
	return &multiCache{
		ttl: ttl,
		lru: lruttl.New(lruSize, ttl),
	}
}

// SetPrefix set prefix for cache
func (mc *multiCache) SetPrefix(prefix string) *multiCache {
	mc.prefix = prefix
	return mc
}

func (mc *multiCache) getKey(key string) string {
	return mc.prefix + key
}

// SetStruct set struct to cache
func (mc *multiCache) SetStruct(ctx context.Context, key string, value interface{}) (err error) {
	buf, err := json.Marshal(value)
	if err != nil {
		return
	}
	key = mc.getKey(key)
	// 先保存至redis
	err = redisSrv.Set(ctx, key, buf, mc.ttl)
	if err != nil {
		return
	}
	// 保存至lru
	mc.lru.Add(key, buf)
	return
}

// GetStruct get struct from cache
func (mc *multiCache) GetStruct(ctx context.Context, key string, value interface{}) (err error) {
	var buf []byte
	key = mc.getKey(key)
	data, ok := mc.lru.Get(key)
	if ok {
		buf, _ = data.([]byte)
	}
	// 如果数据为空，则从redis中拉取
	if len(buf) == 0 {
		buf, err = redisSrv.Get(ctx, key)
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(buf, value)
	if err != nil {
		return
	}
	return
}
