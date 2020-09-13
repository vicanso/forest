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
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/ent"
	"github.com/vicanso/forest/ent/schema"
	"github.com/vicanso/forest/ent/user"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/forest/validate"
	"github.com/vicanso/hes"
)

type (
	userCtrl struct{}

	// userInfoResp 用户信息响应
	userInfoResp struct {
		Date string `json:"date,omitempty"`
		service.UserSessionInfo
	}

	// userListResp 用户列表响应
	userListResp struct {
		Users []*ent.User `json:"users,omitempty"`
	}

	// listUserParams 用户查询参数
	listUserParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		Role    string `json:"role,omitempty" validate:"omitempty,xUserRole"`
		Group   string `json:"group,omitempty" validate:"omitempty,xUserGroup"`
		Status  string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}

	// registerLoginUserParams 注册与登录参数
	registerLoginUserParams struct {
		// 账户
		Account string `json:"account,omitempty" validate:"xUserAccount"`
		// 密码，密码为sha256后的加密串
		Password string `json:"password,omitempty" validate:"xUserPassword"`
	}

	// updateMeParams 用户信息更新参数
	updateMeParams struct {
		Name        string `json:"name,omitempty" validate:"omitempty,xUserName"`
		Email       string `json:"email,omitempty" validate:"omitempty,xUserEmail"`
		Password    string `json:"password,omitempty" validate:"omitempty,xUserPassword"`
		NewPassword string `json:"newPassword,omitempty" validate:"omitempty,xUserPassword"`
	}
)

var (
	sessionConfig config.SessionConfig
)

const (
	errUserCategory = "user"
)

var (
	errLoginTokenNil = &hes.Error{
		Message:    "登录令牌不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errAccountOrPasswordInvalid = &hes.Error{
		Message:    "账户或者密码错误",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errOldPasswordWrong = &hes.Error{
		Message:    "旧密码错误，请重新输入",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errUserStatusInvalid = &hes.Error{
		Message:    "该账户不允许登录",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errUserAccountExists = &hes.Error{
		Message:    "该账户已注册",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
	errUserRcmderNotExists = &hes.Error{
		Message:    "推荐人编号不存在，请重新填写",
		StatusCode: http.StatusBadRequest,
		Category:   errUserCategory,
	}
)

func init() {
	sessionConfig = config.GetSessionConfig()
	g := router.NewGroup("/users", loadUserSession)

	ctrl := userCtrl{}

	// 获取用户列表
	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

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

	// 用户注册
	g.POST(
		"/v1/me",
		newTracker(cs.ActionRegister),
		captchaValidate,
		// 限制相同IP在60秒之内只能调用5次
		newIPLimit(5, 60*time.Second, cs.ActionRegister),
		shouldBeAnonymous,
		ctrl.register,
	)

	// 用户登录
	// 限制3秒只能登录一次（无论成功还是失败）
	loginLimit := newConcurrentLimit([]string{
		"account",
	}, 3*time.Second, cs.ActionLogin)
	g.POST(
		"/v1/me/login",
		middleware.WaitFor(time.Second, true),
		newTracker(cs.ActionLogin),
		captchaValidate,
		shouldBeAnonymous,
		loginLimit,
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		// 限制10分钟内，相同的账号只允许出错5次
		newErrorLimit(5, 10*time.Minute, func(c *elton.Context) string {
			return gjson.GetBytes(c.RequestBody, "account").String()
		}),
		ctrl.login,
	)

	// 刷新user session的ttl
	g.PATCH(
		"/v1/me",
		newTracker(cs.ActionUserMeUpdate),
		ctrl.updateMe,
	)

	// 用户退出登录
	g.DELETE(
		"/v1/me",
		newTracker(cs.ActionLogout),
		shouldBeLogined,
		ctrl.logout,
	)
}

// save 创建用户
func (params *registerLoginUserParams) save(ctx context.Context) (*ent.User, error) {
	return getEntClient().User.Create().
		SetAccount(params.Account).
		SetPassword(params.Password).
		Save(ctx)
}

// login 登录
func (params *registerLoginUserParams) login(ctx context.Context, token string) (u *ent.User, err error) {
	u, err = getEntClient().User.Query().
		Where(user.Account(params.Account)).
		First(ctx)
	if err != nil {
		return
	}
	pwd := util.Sha256(u.Password + token)
	// 用于自动化测试使用
	if util.IsDevelopment() && params.Password == "fEqNCco3Yq9h5ZUglD3CZJT4lBsfEqNCco31Yq9h5ZUB" {
		pwd = params.Password
	}
	if pwd != params.Password {
		err = errAccountOrPasswordInvalid
		return
	}
	return
}

// update 更新用户信息
func (params *updateMeParams) update(ctx context.Context, account string) (u *ent.User, err error) {

	u, err = getEntClient().User.Query().
		Where(user.Account(account)).
		First(ctx)
	if err != nil {
		return
	}
	// 更新密码时需要先校验旧密码
	if params.NewPassword != "" {
		if u.Password != params.Password {
			err = errOldPasswordWrong
			return
		}
	}
	updateOne := u.Update()
	if params.Name != "" {
		updateOne = updateOne.SetName(params.Name)
	}
	if params.Email != "" {
		updateOne = updateOne.SetEmail(params.Email)
	}
	if params.NewPassword != "" {
		updateOne = updateOne.SetPassword(params.NewPassword)
	}
	return updateOne.Save(ctx)
}

// queryAll 查询用户列表
func (params *listUserParams) queryAll(ctx context.Context) (users []*ent.User, err error) {
	query := getEntClient().User.Query()

	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	if params.Keyword != "" {
		query = query.Where(user.AccountContains(params.Keyword))
	}
	// TODO role的查询
	if params.Role != "" {
	}
	if params.Status != "" {
		v, _ := strconv.Atoi(params.Status)
		query = query.Where(user.Status(int8(v)))
	}

	return query.All(ctx)
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

// list 获取用户列表
func (userCtrl) list(c *elton.Context) (err error) {
	params := listUserParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	users, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.Body = &userListResp{
		Users: users,
	}

	return
}

// getLoginToken 获取登录的token
func (userCtrl) getLoginToken(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除当前session id，确保每次登录的用户都是新的session
	err = us.Destroy()
	if err != nil {
		return
	}
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
func (userCtrl) me(c *elton.Context) (err error) {
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

// register 用户注册
func (userCtrl) register(c *elton.Context) (err error) {
	params := registerLoginUserParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	user, err := params.save(c.Context())
	if err != nil {
		return
	}
	// 第一个创建的用户添加su权限
	if user.ID == 1 {
		go func() {
			_, _ = getEntClient().User.UpdateOneID(user.ID).
				SetRoles([]string{
					schema.UserRoleSu,
				}).Save(context.Background())
		}()
	}
	c.Body = user
	return
}

func (userCtrl) login(c *elton.Context) (err error) {
	params := registerLoginUserParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	userInfo, err := us.GetInfo()
	if err != nil {
		return
	}

	if userInfo.Token == "" {
		err = errLoginTokenNil
		return
	}
	// 登录
	u, err := params.login(c.Context(), userInfo.Token)
	if err != nil {
		return
	}

	// 设置session
	err = us.SetInfo(service.UserSessionInfo{
		Account: u.Account,
		ID:      u.ID,
		Roles:   u.Roles,
		// Groups: u.,
	})
	if err != nil {
		return
	}

	ip := c.RealIP()
	trackID := util.GetTrackID(c)
	sessionID := util.GetSessionID(c)
	userAgent := c.GetRequestHeader("User-Agent")

	// 记录至数据库
	// loginRecord := &service.UserLoginRecord{
	// 	Account:       params.Account,
	// 	UserAgent:     userAgent,
	// 	IP:            c.RealIP(),
	// 	TrackID:       trackID,
	// 	SessionID:     sessionID,
	// 	XForwardedFor: c.GetRequestHeader("X-Forwarded-For"),

	// 	Width:         deviceInfo.Width,
	// 	Height:        deviceInfo.Height,
	// 	PixelRatio:    deviceInfo.PixelRatio,
	// 	Platform:      deviceInfo.Platform,
	// 	UUID:          deviceInfo.UUID,
	// 	SystemVersion: deviceInfo.SystemVersion,
	// 	Brand:         deviceInfo.Brand,
	// 	Version:       deviceInfo.Version,
	// 	BuildNumber:   deviceInfo.BuildNumber,
	// }

	// 记录用户登录行为
	getInfluxSrv().Write(cs.MeasurementUserLogin, map[string]interface{}{
		"account":   params.Account,
		"userAgent": userAgent,
		"ip":        ip,
		"trackID":   trackID,
		"sessionID": sessionID,
	}, map[string]string{})

	// 返回用户信息
	resp, err := pickUserInfo(c)
	if err != nil {
		return
	}
	c.Body = &resp
	return
}

// logout 退出登录
func (userCtrl) logout(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除session
	err = us.Destroy()
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// refresh 刷新用户session
func (ctrl userCtrl) refresh(c *elton.Context) (err error) {
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
	c.AddSignedCookie(&http.Cookie{
		Name:     scf.Key,
		Value:    cookie.Value,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})

	c.NoContent()
	return
}

// updateMe 更新用户信息
func (ctrl userCtrl) updateMe(c *elton.Context) (err error) {
	// 如果没有数据要更新，如{}
	if len(c.RequestBody) <= 2 {
		return ctrl.refresh(c)
	}
	us := getUserSession(c)
	// 如果获取不到session，则直接返回
	if us == nil {
		c.NoContent()
		return
	}
	// 如果未登录，无法修改用户信息
	if !us.IsLogined() {
		err = errShouldLogin
		return
	}
	params := updateMeParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	// 更新用户信息
	_, err = params.update(c.Context(), us.MustGetInfo().Account)
	if err != nil {
		return
	}
	c.NoContent()
	return
}
