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

// 管理员的相关一些功能

package controller

import (
	"github.com/vicanso/elton"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/router"
)

type adminCtrl struct{}

type (
	listCacheParams struct {
		Keyword string `json:"keyword" validate:"required,xKeyword"`
		Offset  int    `json:"offset"`
		Limit   int    `json:"limit" default:"1000"`
	}
)

type (
	findCacheResp struct {
		Data string `json:"data"`
	}
	listCacheResp struct {
		Keys []string `json:"keys"`
	}
)

func init() {
	ctrl := adminCtrl{}
	g := router.NewGroup("/@admin", loadUserSession, shouldBeAdmin)

	g.GET(
		"/v1/caches",
		ctrl.listCache,
	)
	// 查询缓存数据
	g.GET(
		"/v1/caches/{key}",
		ctrl.findCacheByKey,
	)
	// 清空session数据
	g.DELETE(
		"/v1/caches/{key}",
		newTrackerMiddleware(cs.ActionAdminCleanCache),
		ctrl.cleanCacheByKey,
	)
}

func (*adminCtrl) listCache(c *elton.Context) error {
	params := listCacheParams{}
	err := validateQuery(c, &params)
	if err != nil {
		return err
	}
	keys, _, err := helper.RedisGetClient().Scan(c.Context(), uint64(params.Offset), params.Keyword, int64(params.Limit)).Result()
	if err != nil {
		return err
	}
	c.Body = &listCacheResp{
		Keys: keys,
	}
	return nil
}

// findCacheByKey find cache by key
func (*adminCtrl) findCacheByKey(c *elton.Context) error {
	data, err := helper.RedisGetClient().Get(c.Context(), c.Param("key")).Result()
	if err != nil {
		return err
	}
	c.Body = &findCacheResp{
		Data: data,
	}
	return nil
}

// cleanCacheByKey clean cache by key
func (*adminCtrl) cleanCacheByKey(c *elton.Context) error {
	_, err := helper.RedisGetClient().Del(c.Context(), c.Param("key")).Result()
	if err != nil {
		return err
	}
	c.NoContent()
	return nil
}
