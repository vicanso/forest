package schedule

import (
	"time"

	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"

	"go.uber.org/zap"
)

func init() {
	if util.IsDevelopment() {
		return
	}
	go initRouteCountTicker()
	go initRedisCheckTicker()
	// go initInfluxdbCheckTicker()
	// go initRouterConfigRefreshTicker()
}

func runTicker(ticker *time.Ticker, message string, do func() error, restart func()) {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			log.Default().DPanic(message+" panic",
				zap.Error(err),
			)
		}
		// 如果退出了，重新启动
		go restart()
	}()
	for range ticker.C {
		err := do()
		// TODO 检测不通过时，发送告警
		if err != nil {
			log.Default().Error(message+" fail",
				zap.Error(err),
			)
		}
	}
}

func initRouteCountTicker() {
	// 每5分钟重置route count
	ticker := time.NewTicker(5 * time.Minute)
	runTicker(ticker, "reset route count", func() error {
		router.ResetRouteCount()
		return nil
	}, initRouteCountTicker)
}

func initRedisCheckTicker() {
	client := service.GetRedisClient()
	// 未使用redis，则不需要检测
	if client == nil {
		return
	}
	// 每一分钟检测一次
	ticker := time.NewTicker(60 * time.Second)
	runTicker(ticker, "redis check", func() error {
		_, err := client.Ping().Result()
		return err
	}, initRedisCheckTicker)
}
