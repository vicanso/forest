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

package controller

import (
	"encoding/json"
	"net"
	"net/http"
	"regexp"
	"strconv"

	"github.com/vicanso/elton"
	M "github.com/vicanso/elton/middleware"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/ent/schema"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/hes"
	intranetip "github.com/vicanso/intranet-ip"
	"go.uber.org/zap"
)

type listParams = helper.EntListParams

var (
	getEntClient = helper.EntGetClient
	now          = util.NowString

	getUserSession = service.NewUserSession
	// 加载用户session
	loadUserSession = elton.Compose(middleware.NewSession(), sessionInterceptorMiddleware)
	// 判断用户是否登录
	shouldBeLogin = checkLoginMiddleware
	// 判断用户是否未登录
	shouldBeAnonymous = checkAnonymousMiddleware
	// 判断用户是否admin权限
	shouldBeAdmin = newCheckRolesMiddleware([]string{
		schema.UserRoleSu,
		schema.UserRoleAdmin,
	})
	// shouldBeSu 判断用户是否su权限
	shouldBeSu = newCheckRolesMiddleware([]string{
		schema.UserRoleSu,
	})

	// 创建新的并发控制中间件
	newConcurrentLimit = middleware.NewConcurrentLimit
	// 创建IP限制中间件
	newIPLimit = middleware.NewIPLimit
	// 创建出错限制中间件
	newErrorLimit = middleware.NewErrorLimit
	// noCacheIfRequestNoCache 请求参数指定no cache，则设置no-cache
	noCacheIfRequestNoCache = middleware.NewNoCacheWithCondition("cacheControl", "no-cache")

	// 图形验证码校验
	captchaValidate = newMagicalCaptchaValidate()
	// 获取influx service
	getInfluxSrv = helper.GetInfluxSrv
	// 文件服务
	fileSrv = &service.FileSrv{}
	// prof服务
	profSrv = &service.ProfSrv{}
)

func newMagicalCaptchaValidate() elton.Handler {
	magicValue := ""
	if !util.IsProduction() {
		magicValue = "0145"
	}
	return middleware.ValidateCaptcha(magicValue)
}

// isLogin 判断是否登录状态
func isLogin(c *elton.Context) bool {
	us := service.NewUserSession(c)
	return us.IsLogin()
}

func validateLogin(c *elton.Context) (err error) {
	if !isLogin(c) {
		err = hes.New("请先登录", errUserCategory)
		return
	}
	return
}

// checkLoginMiddleware 校验是否登录中间件
func checkLoginMiddleware(c *elton.Context) (err error) {
	err = validateLogin(c)
	if err != nil {
		return
	}
	return c.Next()
}

// checkAnonymousMiddleware 判断是匿名状态
func checkAnonymousMiddleware(c *elton.Context) (err error) {
	if isLogin(c) {
		err = hes.New("已是登录状态，请先退出登录", errUserCategory)
		return
	}
	return c.Next()
}

// newCheckRolesMiddleware 创建用户角色校验中间件
func newCheckRolesMiddleware(validRoles []string) elton.Handler {
	return func(c *elton.Context) (err error) {
		err = validateLogin(c)
		if err != nil {
			return
		}
		us := service.NewUserSession(c)
		userInfo, err := us.GetInfo()
		if err != nil {
			return
		}
		valid := util.ContainsAny(validRoles, userInfo.Roles)
		if valid {
			return c.Next()
		}
		err = hes.NewWithStatusCode("禁止使用该功能", http.StatusForbidden, errUserCategory)
		return
	}
}

// newTrackerMiddleware 初始化用户行为跟踪中间件
func newTrackerMiddleware(action string) elton.Handler {
	marshalString := func(data interface{}) string {
		buf, _ := json.Marshal(data)
		return string(buf)
	}
	return M.NewTracker(M.TrackerConfig{
		Mask: regexp.MustCompile(`(?i)password`),
		OnTrack: func(info *M.TrackerInfo, c *elton.Context) {
			account := ""
			tid := util.GetTrackID(c)
			us := service.NewUserSession(c)
			if us != nil && us.IsLogin() {
				account = us.MustGetInfo().Account
			}
			ip := c.RealIP()
			sid := util.GetSessionID(c)
			zapFields := make([]zap.Field, 0, 10)
			zapFields = append(
				zapFields,
				zap.String("action", action),
				zap.String("account", account),
				zap.String("ip", ip),
				zap.String("sid", sid),
				zap.String("tid", tid),
				zap.Int("result", info.Result),
			)
			fields := map[string]interface{}{
				cs.FieldAccount: account,
				cs.FieldIP:      ip,
				cs.FieldSID:     sid,
				cs.FieldTID:     tid,
			}
			if len(info.Query) != 0 {
				zapFields = append(zapFields, zap.Any("query", info.Query))
				fields[cs.FieldQuery] = marshalString(info.Query)
			}
			if len(info.Params) != 0 {
				zapFields = append(zapFields, zap.Any("params", info.Params))
				fields[cs.FieldParams] = marshalString(info.Params)
			}
			if len(info.Form) != 0 {
				zapFields = append(zapFields, zap.Any("form", info.Form))
				fields[cs.FieldForm] = marshalString(info.Form)
			}
			if info.Err != nil {
				zapFields = append(zapFields, zap.Error(info.Err))
				fields[cs.FieldError] = info.Err.Error()
			}
			log.Default().Info("tracker", zapFields...)
			getInfluxSrv().Write(cs.MeasurementUserTracker, map[string]string{
				cs.TagAction: action,
				cs.TagResult: strconv.Itoa(info.Result),
			}, fields)
		},
	})
}

// getIDFromParams get id form context params
func getIDFromParams(c *elton.Context) (id int, err error) {
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		he := hes.Wrap(err)
		he.Category = "parseInt"
		err = he
		return
	}
	return
}

// sessionInterceptorMiddleware session的拦截
func sessionInterceptorMiddleware(c *elton.Context) error {
	interData, _ := service.GetSessionInterceptorData()
	// 如果无配置，则直接跳过
	if interData == nil {
		return c.Next()
	}
	us := service.NewUserSession(c)
	account := ""
	if us.IsLogin() {
		account = us.MustGetInfo().Account
	}
	// 如果配置该账号允许
	if account != "" && util.ContainsString(interData.AllowAccounts, account) {
		return c.Next()
	}
	// 如果路由配置允许
	if util.ContainsString(interData.AllowRoutes, c.Route) {
		return c.Next()
	}

	// 如果有配置拦截信息，则以出错返回
	he := hes.New(interData.Message)
	he.Category = "sessionInterceptorMiddleware"
	return he
}

// isIntranet 判断是否内网访问
func isIntranet(c *elton.Context) error {
	if intranetip.Is(net.ParseIP(c.ClientIP())) {
		return c.Next()
	}
	return hes.NewWithStatusCode("Forbidden", 403)
}
