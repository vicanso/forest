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
	"time"

	"github.com/vicanso/forest/validate"
	"github.com/vicanso/hes"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

type (
	userCtrl struct{}
	// userInfoResp user info response
	userInfoResp struct {
		Anonymous bool     `json:"anonymous,omitempty"`
		Account   string   `json:"account,omitempty"`
		Roles     []string `json:"roles,omitempty"`
		Date      string   `json:"date,omitempty"`
		UpdatedAt string   `json:"updatedAt,omitempty"`
		IP        string   `json:"ip,omitempty"`
		TrackID   string   `json:"trackId,omitempty"`
		LoginAt   string   `json:"loginAt,omitempty"`
	}

	registerUserParams struct {
		Account  string `json:"account,omitempty" valid:"xUserAccount"`
		Password string `json:"password,omitempty" valid:"xUserPassword"`
	}

	loginUserParams struct {
		Account  string `json:"account,omitempty" valid:"xUserAccount"`
		Password string `json:"password,omitempty" valid:"xUserPassword"`
	}
)

var (
	errLoginTokenNil = hes.New("login token is nil")
)

func init() {
	g := router.NewGroup("/users", loadUserSession)
	ctrl := userCtrl{}

	// 获取用户信息
	g.GET("/v1/me", ctrl.me)

	// 用户注册
	g.POST(
		"/v1/me",
		newTracker(cs.ActionRegister),
		// 限制相同IP在60秒之内只能调用5次
		newIPLimit(5, 60*time.Second, cs.ActionLogin),
		shouldAnonymous,
		ctrl.register,
	)
	// 刷新user session的ttl
	g.PATCH(
		"/v1/me",
		ctrl.refresh,
	)

	// 获取登录token
	g.GET(
		"/v1/me/login",
		shouldAnonymous,
		ctrl.getLoginToken,
	)

	// 用户登录
	// 限制3秒只能登录一次（无论成功还是失败）
	loginLimit := newConcurrentLimit([]string{
		"account",
	}, 3*time.Second, cs.ActionLogin)
	g.POST(
		"/v1/me/login",
		newTracker(cs.ActionLogin),
		shouldAnonymous,
		loginLimit,
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		ctrl.login,
	)
	// 用户退出登录
	g.DELETE(
		"/v1/me/logout",
		newTracker(cs.ActionLogout),
		shouldLogined,
		ctrl.logout,
	)
}

// get user info from session
func pickUserInfo(c *cod.Context) (userInfo *userInfoResp) {
	us := getUserSession(c)
	userInfo = &userInfoResp{
		Anonymous: true,
		Date:      now(),
		IP:        c.RealIP(),
		TrackID:   getTrackID(c),
	}
	account := us.GetAccount()
	if account != "" {
		userInfo.Account = account
		userInfo.Roles = us.GetRoles()
		userInfo.Anonymous = false
	}
	return
}

// get user info
func (ctrl userCtrl) me(c *cod.Context) (err error) {
	key := config.GetTrackKey()
	cookie, _ := c.Cookie(key)
	// ulid的长度为26
	if cookie == nil || len(cookie.Value) != 26 {
		c.AddCookie(&http.Cookie{
			Name:     key,
			Value:    util.GenUlid(),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   365 * 24 * 3600,
		})
	}
	c.Body = pickUserInfo(c)
	return
}

// getLoginToken get login token
func (ctrl userCtrl) getLoginToken(c *cod.Context) (err error) {
	us := getUserSession(c)
	// 清除当前session id，确保每次登录的用户都是新的session
	us.ClearSessionID()
	token := util.RandomString(8)
	err = us.SetLoginToken(token)
	if err != nil {
		return
	}
	c.Body = &struct {
		Token string `json:"token,omitempty"`
	}{
		token,
	}
	return
}

func omitUserInfo(u *service.User) {
	u.Password = ""
}

// register user register
func (ctrl userCtrl) register(c *cod.Context) (err error) {
	params := &registerUserParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	u := &service.User{
		Account:  params.Account,
		Password: params.Password,
	}
	err = service.UserAdd(u)
	if err != nil {
		return
	}
	omitUserInfo(u)
	c.Created(u)
	return
}

// login user login
func (ctrl userCtrl) login(c *cod.Context) (err error) {
	params := &registerUserParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	token := us.GetLoginToken()
	if token == "" {
		err = errLoginTokenNil
		return
	}
	u, err := service.UserLogin(params.Account, params.Password, token)
	if err != nil {
		return
	}
	service.UserLoginRecordAdd(&service.UserLoginRecord{
		Account:   params.Account,
		UserAgent: c.GetRequestHeader("User-Agent"),
		IP:        c.RealIP(),
		TrackID:   util.GetTrackID(c),
		SessionID: util.GetSessionID(c),
	})
	omitUserInfo(u)
	us.SetAccount(u.Account)
	us.SetRoles(u.Roles)
	c.Body = u
	return
}

// logout user logout
func (ctrl userCtrl) logout(c *cod.Context) (err error) {
	us := getUserSession(c)
	if us != nil {
		err = us.Destroy()
	}
	c.NoContent()
	return
}

// refresh user refresh
func (ctrl userCtrl) refresh(c *cod.Context) (err error) {
	us := getUserSession(c)
	if us == nil {
		c.NoContent()
		return
	}

	scf := config.GetSessionConfig()
	cookie, _ := c.SignedCookie(scf.Key)
	// 如果认证的cookie已过期，则不做刷新
	if cookie == nil {
		c.NoContent()
		return
	}

	err = us.Refresh()
	if err != nil {
		return
	}
	// 更新session
	err = c.AddSignedCookie(&http.Cookie{
		Name:     scf.Key,
		Value:    cookie.Value,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})
	if err != nil {
		return
	}

	c.NoContent()
	return
}
