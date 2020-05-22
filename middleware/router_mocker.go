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

package middleware

import (
	"bytes"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/service"
)

// NewRouterMocker create a router mocker
func NewRouterMocker() elton.Handler {
	return func(c *elton.Context) (err error) {
		routerConfig := service.RouterGetConfig(c.Request.Method, c.Route)
		if routerConfig == nil {
			return c.Next()
		}

		// 如果有配置Path，则还要判断path是否相等
		if routerConfig.URL != "" && c.Request.URL.RequestURI() != routerConfig.URL {
			return c.Next()
		}

		// 如果delay大于0，则延时
		if routerConfig.Delay > 0 {
			time.Sleep(time.Second * time.Duration(routerConfig.Delay))
		}

		c.StatusCode = routerConfig.Status
		contentType := routerConfig.CotentType
		if contentType == "" {
			contentType = elton.MIMEApplicationJSON
		}
		c.SetHeader(elton.HeaderContentType, contentType)
		c.BodyBuffer = bytes.NewBufferString(routerConfig.Response)
		return nil
	}
}
