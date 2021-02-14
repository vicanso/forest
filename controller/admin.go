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
	"github.com/vicanso/forest/cache"
	"github.com/vicanso/forest/router"
)

type (
	adminCtrl struct{}

	findSessionResp struct {
		Data string `json:"data,omitempty"`
	}
)

func init() {
	ctrl := adminCtrl{}
	g := router.NewGroup("/@admin", loadUserSession, shouldBeAdmin)

	// 查询session数据
	g.GET(
		"/v1/sessions/{id}",
		ctrl.findSessionByID,
	)
	// 清空session数据
	g.DELETE(
		"/v1/sessions/{id}",
		ctrl.cleanSessionByID,
	)
}

// findSessionByID find session by id
func (*adminCtrl) findSessionByID(c *elton.Context) (err error) {
	store := cache.GetRedisSession()
	data, err := store.Get(c.Param("id"))
	if err != nil {
		return
	}
	c.Body = &findSessionResp{
		Data: string(data),
	}
	return
}

// cleanSessionByID clean session by id
func (*adminCtrl) cleanSessionByID(c *elton.Context) (err error) {
	store := cache.GetRedisSession()
	err = store.Destroy(c.Param("id"))
	if err != nil {
		return
	}
	c.NoContent()
	return
}
