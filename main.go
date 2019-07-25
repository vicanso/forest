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

package main

import (
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/vicanso/cod"
	bodyparser "github.com/vicanso/cod-body-parser"
	compress "github.com/vicanso/cod-compress"
	errorHandler "github.com/vicanso/cod-error-handler"
	etag "github.com/vicanso/cod-etag"
	fresh "github.com/vicanso/cod-fresh"
	recover "github.com/vicanso/cod-recover"
	responder "github.com/vicanso/cod-responder"
	stats "github.com/vicanso/cod-stats"

	"go.uber.org/zap"

	"github.com/vicanso/forest/config"
	_ "github.com/vicanso/forest/controller"
	_ "github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/router"
	_ "github.com/vicanso/forest/schedule"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

// 相关依赖服务的校验，主要是数据库等
func dependServiceCheck() (err error) {
	err = service.RedisPing()
	if err != nil {
		return
	}
	configSrv := new(service.ConfigurationSrv)
	err = configSrv.Refresh()
	if err != nil {
		return
	}
	return
}

func main() {
	logger := log.Default()
	d := cod.New()
	d.SignedKeys = service.GetSignedKeys()

	// 未处理的error才会触发
	d.OnError(func(c *cod.Context, err error) {
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
	})
	// TODO 对于404的请求，不会执行中间件，一般都是因为攻击之类才会导致大量出现404，
	// 因此可在此处汇总出错IP，针对较频繁出错IP，增加告警信息
	d.NotFoundHandler = func(resp http.ResponseWriter, req *http.Request) {
		logger.Info("404",
			zap.String("ip", cod.GetRealIP(req)),
			zap.String("uri", req.RequestURI),
		)
		resp.Header().Set(cod.HeaderContentType, cod.MIMEApplicationJSON)
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"statusCode": 404,"message": "Not found"}`))
	}

	// 捕捉panic异常，避免程序崩溃
	d.Use(recover.New())

	d.Use(middleware.NewEntry())

	// 接口相关统计信息
	d.Use(stats.New(stats.Config{
		OnStats: func(info *stats.Info, c *cod.Context) {
			// ping 的日志忽略
			if info.URI == "/ping" {
				return
			}
			logger.Info("access log",
				zap.String("id", info.CID),
				zap.String("ip", info.IP),
				zap.String("sid", util.GetSessionID(c)),
				zap.String("method", info.Method),
				zap.String("uri", info.URI),
				zap.Int("status", info.Status),
				zap.String("consuming", info.Consuming.String()),
				zap.String("size", humanize.Bytes(uint64(info.Size))),
			)
		},
	}))

	// 根据应用配置限制路由
	d.Use(middleware.NewRouterController())

	// 错误处理，将错误转换为json响应
	d.Use(errorHandler.NewDefault())

	// 压缩响应数据
	d.Use(compress.NewDefault())

	// etag与fresh的处理
	d.Use(fresh.NewDefault())
	d.Use(etag.NewDefault())

	// 对响应数据 c.Body 转换为相应的json响应
	d.Use(responder.NewDefault())

	// 读取读取body的数的，转换为json bytes
	d.Use(bodyparser.NewDefault())

	// 初始化路由
	router.Init(d)

	err := dependServiceCheck()
	if err != nil {
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.Error(err),
		)
		panic(err)
	}
	d.ListenAndServe(config.GetListen())
}
