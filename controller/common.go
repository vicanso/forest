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

// 公共的处理函数，包括程序基本信息、性能指标等

package controller

import (
	"bytes"
	"net/http"
	"runtime"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/hes"
)

type (
	commonCtrl struct{}

	applicationInfoResp struct {
		// 版本号
		Version string `json:"version,omitempty"`
		// 构建时间
		BuildedAt string `json:"buildedAt,omitempty"`
		// 运行时长
		Uptime string `json:"uptime,omitempty"`
		// os类型
		OS string `json:"os,omitempty"`
		// go版本
		GO string `json:"go,omitempty"`
		// 架构类型
		ARCH string `json:"arch,omitempty"`
		// 运行环境配置
		ENV string `json:"env,omitempty"`
	}
	routersResp struct {
		// 路由信息
		Routers []*elton.RouterInfo `json:"routers,omitempty"`
	}
)

var (
	// applicationStartedAt 应用启动时间
	applicationStartedAt = time.Now()

	// errAppIsNotRunning 应用非运行状态
	errAppIsNotRunning = &hes.Error{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "应用服务不可用",
	}
)

func init() {
	ctrl := commonCtrl{}
	router.NewGroup("").GET("/ping", ctrl.ping)
	g := router.NewGroup("/commons")

	g.GET("/application", ctrl.getApplicationInfo)
	g.GET("/routers", ctrl.routers)
	g.GET("/captcha", ctrl.captcha)
	g.GET("/performance", ctrl.getPerformance)
}

// ping 用于检测服务是否可用
func (commonCtrl) ping(c *elton.Context) error {
	if !config.ApplicationIsRunning() {
		return errAppIsNotRunning
	}
	c.BodyBuffer = bytes.NewBufferString("pong")
	return nil
}

// getApplicationInfo 获取应用信息
func (commonCtrl) getApplicationInfo(c *elton.Context) (err error) {
	c.CacheMaxAge("1m")
	c.Body = &applicationInfoResp{
		config.GetApplicationVersion(),
		config.GetApplicationBuildedAt(),
		time.Since(applicationStartedAt).String(),
		runtime.GOOS,
		runtime.Version(),
		runtime.GOARCH,
		config.GetENV(),
	}
	return
}

// routers 获取系统的路由
func (commonCtrl) routers(c *elton.Context) (err error) {
	c.CacheMaxAge("1m")
	c.Body = &routersResp{
		Routers: c.Elton().Routers,
	}
	return
}

// captcha 获取图形验证码
func (commonCtrl) captcha(c *elton.Context) (err error) {
	bgColor := c.QueryParam("bg")
	fontColor := c.QueryParam("color")
	if bgColor == "" {
		bgColor = "255,255,255"
	}
	if fontColor == "" {
		fontColor = "102,102,102"
	}
	info, err := service.GetCaptcha(fontColor, bgColor)
	if err != nil {
		return
	}
	// c.SetContentTypeByExt(".jpeg")
	// c.Body = info.Data
	c.NoStore()
	c.Body = info
	return
}

// getPerformance 获取应用性能指标
func (commonCtrl) getPerformance(c *elton.Context) (err error) {
	p := service.GetPerformance()
	c.Body = &p
	return
}
