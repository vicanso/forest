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

package controller

import (
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/hes"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	tracker "github.com/vicanso/elton-tracker"
)

var (
	errShouldLogin  = hes.New("should login first")
	errLoginAlready = hes.New("login already, please logout first")
	errForbidden    = &hes.Error{
		StatusCode: http.StatusForbidden,
		Message:    "acccess forbidden",
	}
)

var (
	logger       = log.Default()
	now          = util.NowString
	getTrackID   = util.GetTrackID
	standardJSON = jsoniter.ConfigCompatibleWithStandardLibrary

	// 服务列表
	// 配置服务
	configSrv = new(service.ConfigurationSrv)
	// 用户服务
	userSrv = new(service.UserSrv)

	// 创建新的并发控制中间件
	newConcurrentLimit = middleware.NewConcurrentLimit
	// 创建IP限制中间件
	newIPLimit = middleware.NewIPLimit
	// 创建出错限制中间件
	newErrorLimit = middleware.NewErrorLimit

	getUserSession = service.NewUserSession
	// 加载用户session
	loadUserSession = middleware.NewSession()
	// 判断用户是否登录
	shouldLogined = elton.Compose(loadUserSession, checkLogin)
	// 判断用户是否未登录
	shouldAnonymous = elton.Compose(loadUserSession, checkAnonymous)
	// 判断用户是否admin权限
	shouldBeAdmin = elton.Compose(loadUserSession, isAdmin)

	// 图形验证码校验
	captchaValidate elton.Handler
)

func init() {
	magicalValue := ""
	if !util.IsProduction() {
		magicalValue = cs.MagicalCaptcha
	}
	captchaValidate = middleware.ValidateCaptch(magicalValue)
}

func newTracker(action string) elton.Handler {
	return tracker.New(tracker.Config{
		OnTrack: func(info *tracker.Info, c *elton.Context) {
			account := ""
			us := service.NewUserSession(c)
			if us != nil {
				account = us.GetAccount()
			}
			logger.Info("tracker",
				zap.String("action", action),
				zap.String("cid", info.CID),
				zap.String("account", account),
				zap.String("ip", c.RealIP()),
				zap.String("sid", util.GetSessionID(c)),
				zap.Int("result", info.Result),
				zap.Any("query", info.Query),
				zap.Any("params", info.Params),
				zap.Any("form", info.Form),
				zap.Error(info.Err),
			)
		},
	})
}

func isLogin(c *elton.Context) bool {
	us := service.NewUserSession(c)
	if us == nil || us.GetAccount() == "" {
		return false
	}
	return true
}

func checkLogin(c *elton.Context) (err error) {
	if !isLogin(c) {
		err = errShouldLogin
		return
	}
	return c.Next()
}

func checkAnonymous(c *elton.Context) (err error) {
	if isLogin(c) {
		err = errLoginAlready
		return
	}
	return c.Next()
}

func newCheckRoles(validRoles []string) elton.Handler {
	return func(c *elton.Context) (err error) {
		if !isLogin(c) {
			err = errShouldLogin
			return
		}
		us := service.NewUserSession(c)
		roles := us.GetRoles()
		valid := util.UserRoleIsValid(validRoles, roles)
		if valid {
			return c.Next()
		}
		err = errForbidden
		return
	}
}

func isAdmin(c *elton.Context) (err error) {
	if !isLogin(c) {
		err = errShouldLogin
		return
	}
	us := service.NewUserSession(c)
	if us.IsAdmin() {
		return c.Next()
	}
	err = errForbidden
	return
}
