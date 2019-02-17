package main

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/cod/middleware"
	"github.com/vicanso/forest/config"
	log "github.com/vicanso/forest/log"
	"github.com/vicanso/forest/router"

	"go.uber.org/zap"
)

func main() {
	logger := log.Default()
	listen := config.GetListen()

	d := cod.New()

	// 如果出错则会触发此回调（在 ErrorHandler 中会将出错转换为相应的http响应，此类情况不会触发）
	d.OnError(func(c *cod.Context, err error) {
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
	})

	d.Use(middleware.NewRecover())

	// 接口响应统计，在项目中可写入数据库方便统计
	d.Use(middleware.NewStats(middleware.StatsConfig{
		OnStats: func(statsInfo *middleware.StatsInfo, _ *cod.Context) {
			logger.Info("access log",
				zap.String("ip", statsInfo.IP),
				zap.String("method", statsInfo.Method),
				zap.String("uri", statsInfo.URI),
				zap.Int("status", statsInfo.Status),
				zap.String("consuming", statsInfo.Consuming.String()),
			)
		},
	}))

	d.Use(middleware.NewDefaultErrorHandler())

	// 设置所有的请求响应默认都为no cache
	d.Use(func(c *cod.Context) error {
		c.NoCache()
		return c.Next()
	})

	d.Use(middleware.NewDefaultFresh())
	d.Use(middleware.NewDefaultETag())
	d.Use(middleware.NewDefaultCompress())

	d.Use(middleware.NewDefaultResponder())

	// health check
	d.GET("/ping", func(c *cod.Context) (err error) {
		c.Body = "pong"
		return
	})

	groups := router.GetGroups()
	for _, g := range groups {
		d.AddGroup(g)
	}

	logger.Info("server is starting",
		zap.String("listen", listen),
	)
	d.ListenAndServe(listen)
}
