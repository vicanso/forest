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
	"time"

	lruttl "github.com/vicanso/lru-ttl"
)

var multiCacheDefaultTimeout = 3 * time.Second

type redisCache struct{}

func (sc *redisCache) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), multiCacheDefaultTimeout)
	defer cancel()
	return redisSrv.Get(ctx, key)
}

func (sc *redisCache) Set(key string, value []byte, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), multiCacheDefaultTimeout)
	defer cancel()
	return redisSrv.Set(ctx, key, value, ttl)
}

// NewMultiCache create a new multi cache
func NewMultiCache(lruSize int, ttl time.Duration) *lruttl.L2Cache {
	return lruttl.NewL2Cache(&redisCache{}, lruSize, ttl)
}

// NewMultiCacheWithPrefix create a new multi cache and set prefix key
func NewMultiCacheWithPrefix(lruSize int, ttl time.Duration, prefix string) *lruttl.L2Cache {
	l2 := NewMultiCache(lruSize, ttl)
	l2.SetPrefix(prefix)
	return l2
}
