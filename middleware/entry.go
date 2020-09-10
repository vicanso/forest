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

package middleware

import (
	"github.com/vicanso/elton"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

const (
	xResponseID = "X-Response-Id"
)

// NewEntry create an entry middleware
func NewEntry() elton.Handler {
	return func(c *elton.Context) (err error) {
		service.IncreaseConcurrency()
		defer service.DecreaseConcurrency()
		c.SetHeader(xResponseID, c.ID)
		// 非测试环境返回x-forwarded-for，方便确认链路
		if !util.IsProduction() {
			c.SetHeader(elton.HeaderXForwardedFor, c.GetRequestHeader(elton.HeaderXForwardedFor))
		}

		// 设置所有的请求响应默认都为no cache
		c.NoCache()

		return c.Next()
	}
}
