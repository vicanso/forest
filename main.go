package main

import (
	"net/http"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	_ "github.com/vicanso/forest/controller"
	"github.com/vicanso/forest/global"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/router"
	_ "github.com/vicanso/forest/schedule"
	"github.com/vicanso/hes"

	"go.uber.org/zap"

	bodyparser "github.com/vicanso/cod-body-parser"
	compress "github.com/vicanso/cod-compress"
	errorHandler "github.com/vicanso/cod-error-handler"
	etag "github.com/vicanso/cod-etag"
	fresh "github.com/vicanso/cod-fresh"
	recover "github.com/vicanso/cod-recover"
	responder "github.com/vicanso/cod-responder"
	stats "github.com/vicanso/cod-stats"
)

func main() {
	logger := log.Default()
	listen := config.GetListen()

	d := cod.New()

	d.Keys = config.GetStringSlice("keys")

	// 如果出错则会触发此回调（在 ErrorHandler 中会将出错转换为相应的http响应，此类情况不会触发）
	d.OnError(func(c *cod.Context, err error) {
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
	})

	d.Use(recover.New())

	d.Use(middleware.NewEntry())

	// 接口响应统计，在项目中可写入数据库方便统计
	d.Use(stats.New(stats.Config{
		OnStats: func(statsInfo *stats.Info, _ *cod.Context) {
			// 增加从session中获取当前账号
			logger.Info("access log",
				zap.String("cid", statsInfo.CID),
				zap.String("ip", statsInfo.IP),
				zap.String("method", statsInfo.Method),
				zap.String("uri", statsInfo.URI),
				zap.Int("status", statsInfo.Status),
				zap.String("consuming", statsInfo.Consuming.String()),
			)
		},
	}))

	d.Use(errorHandler.NewDefault())

	d.Use(middleware.NewLimiter())

	d.Use(bodyparser.NewDefault())

	d.Use(fresh.NewDefault())
	d.Use(etag.NewDefault())
	d.Use(compress.NewDefault())

	d.Use(responder.NewDefault())

	// health check
	d.GET("/ping", func(c *cod.Context) (err error) {
		if !global.IsApplicationRunning() {
			err = hes.NewWithStatusCode("application is not running", http.StatusServiceUnavailable)
			return
		}
		c.Body = "pong"
		return
	})

	groups := router.GetGroups()
	for _, g := range groups {
		d.AddGroup(g)
	}

	router.InitRouteCounter(d.Routers)

	logger.Info("server is starting",
		zap.String("listen", listen),
	)

	// 设置应用状态为运行中
	global.StartApplication()
	err := d.ListenAndServe(listen)
	if err != nil {
		panic(err)
	}
}
