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

// 用户相关的一些路由处理

package controller

import (
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

type (
	userCtrl struct{}

	userInfoResp struct {
		Date string `json:"date,omitempty"`
		service.UserSessionInfo
	}
)

var (
	sessionConfig config.SessionConfig
)

func init() {
	sessionConfig = config.GetSessionConfig()
	g := router.NewGroup("/users", loadUserSession)

	ctrl := userCtrl{}

	// 获取登录token
	g.GET(
		"/v1/me/login",
		shouldBeAnonymous,
		ctrl.getLoginToken,
	)
	// 获取用户信息
	g.GET(
		"/v1/me",
		ctrl.me,
	)
}

// pickUserInfo 获取用户信息
func pickUserInfo(c *elton.Context) (resp userInfoResp, err error) {
	us := getUserSession(c)
	userInfo, err := us.GetInfo()
	if err != nil {
		return
	}
	resp = userInfoResp{
		Date: util.NowString(),
	}
	resp.UserSessionInfo = userInfo
	return
}

// getLoginToken 获取登录的token
func (userCtrl) getLoginToken(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除当前session id，确保每次登录的用户都是新的session
	us.ClearSessionID()
	userInfo := service.UserSessionInfo{
		Token: util.RandomString(8),
	}
	err = us.SetInfo(userInfo)
	if err != nil {
		return
	}
	c.Body = &userInfo
	return
}

// me 获取用户信息
func (ctrl userCtrl) me(c *elton.Context) (err error) {
	cookie, _ := c.Cookie(sessionConfig.TrackKey)
	// ulid的长度为26
	if cookie == nil || len(cookie.Value) != 26 {
		uid := util.GenUlid()
		c.AddCookie(&http.Cookie{
			Name:     sessionConfig.TrackKey,
			Value:    uid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   365 * 24 * 3600,
		})
		// trackRecord := &service.UserTrackRecord{
		// 	UserAgent: c.GetRequestHeader("User-Agent"),
		// 	IP:        c.RealIP(),
		// 	TrackID:   util.GetTrackID(c),
		// }
		// _ = userSrv.AddTrackRecord(trackRecord, c)
	}
	resp, err := pickUserInfo(c)
	if err != nil {
		return
	}
	c.Body = &resp
	return
}
