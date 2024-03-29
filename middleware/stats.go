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
	"github.com/dustin/go-humanize"
	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/util"
	"go.uber.org/atomic"
)

func NewStats(processingCount *atomic.Int32) elton.Handler {
	return M.NewStats(M.StatsConfig{
		OnStats: func(info *M.StatsInfo, c *elton.Context) {
			// ping 的日志忽略
			if info.URI == "/ping" {
				return
			}
			// TODO 如果需要可以从session中获取账号
			// 此处记录的session id，需要从客户登录记录中获取对应的session id
			// us := service.NewUserSession(c)
			sid := util.GetSessionID(c)
			requestBodySize := len(c.RequestBody)
			processing := processingCount.Load()
			// 由客户端设置的uuid
			// zap.String("uuid", c.GetRequestHeader("X-UUID")),
			log.Info(c.Context()).
				Str("category", "accessLog").
				Str("ip", info.IP).
				Str("sid", sid).
				Str("method", info.Method).
				Str("route", info.Route).
				Str("uri", info.URI).
				Int("status", info.Status).
				Int32("connecting", processing).
				Str("latency", info.Latency.String()).
				Str("requestBodySize", humanize.Bytes(uint64(requestBodySize))).
				Str("size", humanize.Bytes(uint64(info.Size))).
				Int("bytes", info.Size).
				Msg("")

			tags := map[string]string{
				cs.TagMethod: info.Method,
				cs.TagRoute:  info.Route,
			}
			fields := map[string]any{
				cs.FieldIP:         info.IP,
				cs.FieldSID:        sid,
				cs.FieldURI:        info.URI,
				cs.FieldStatus:     info.Status,
				cs.FieldLatency:    int(info.Latency.Milliseconds()),
				cs.FieldBodySize:   requestBodySize,
				cs.FieldSize:       info.Size,
				cs.FieldProcessing: processing,
			}
			helper.GetInfluxDB().Write(cs.MeasurementHTTPStats, tags, fields)
		},
	})
}
