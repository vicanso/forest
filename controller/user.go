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
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/tidwall/gjson"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/validate"
	"github.com/vicanso/hes"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

type userCtrl struct{}
type (
	userInfoResp struct {
		// 是否匿名
		// Example: true
		Anonymous bool `json:"anonymous"`
		// 账号
		// Example: vicanso
		Account string `json:"account"`
		// 角色
		// Example: ["su", "admin"]
		Roles []string `json:"roles"`
		// 分组
		// Example: ["it", "finance"]
		Groups []string `json:"groups"`
		// 系统时间
		// Example: 2019-10-26T10:11:25+08:00
		Date string `json:"date"`
		// 信息更新时间
		// Example: 2019-10-26T10:11:25+08:00
		UpdatedAt string `json:"updatedAt"`
		// IP地址
		// Example: 1.1.1.1
		IP string `json:"ip"`
		// rack id
		// Example: 01DPNPDXH4MQJHBF4QX1EFD6Y3
		TrackID string `json:"trackId"`
		// 登录时间
		// Example: 2019-10-26T10:11:25+08:00
		LoginAt string `json:"loginAt"`
	}
	loginTokenResp struct {
		// 登录Token
		// Example: IaHnYepm
		Token string `json:"token"`
	}
)

type (
	// 注册与登录参数
	registerLoginUserParams struct {
		// 账户
		// Example: vicanso
		Account string `json:"account" validate:"xUserAccount"`
		// 密码，密码为sha256后的加密串
		// Example: JgX9742WqzaNHVP+YiPy/RXP0eoX29k00hEF3BdghGU=
		Password string `json:"password" validate:"xUserPassword"`
	}

	listUserParams struct {
		listParams
		Keyword string `json:"keyword" validate:"omitempty,xKeyword"`
		Role    string `json:"role" validate:"omitempty,xUserRole"`
		Group   string `json:"group" validate:"omitempty,xUserGroup"`
		Status  string `json:"status" validate:"omitempty,xUserStatusString"`
	}

	updateUserParams struct {
		Roles  []string `json:"roles" validate:"omitempty,xUserRoles"`
		Groups []string `json:"groups" validate:"omitempty,xUserGroups"`
		Status int      `json:"status" validate:"omitempty,xUserStatus"`
	}
	updateMeParams struct {
		Email       string `json:"email" validate:"omitempty,xUserEmail"`
		Mobile      string `json:"mobile" validate:"omitempty,xUserMobile"`
		Password    string `json:"password" validate:"omitempty,xUserPassword"`
		NewPassword string `json:"newPassword" validate:"omitempty,xUserPassword"`
	}
	listUserLoginRecordParams struct {
		listParams
		Begin   time.Time `json:"begin"`
		End     time.Time `json:"end"`
		Account string    `json:"account" validate:"omitempty,xUserAccount"`
	}
)

var (
	errLoginTokenNil = hes.New("login token is nil")
)

func init() {

	g := router.NewGroup("/users", loadUserSession)
	ctrl := userCtrl{}
	// 获取用户列表
	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

	// 更新用户信息
	g.PATCH(
		"/v1/{id}",
		shouldBeAdmin,
		newTracker(cs.ActionUserInfoUpdate),
		ctrl.updateByID,
	)

	// 获取用户信息
	g.GET("/v1/me", ctrl.me)
	g.GET("/v1/me/profile", shouldLogined, ctrl.profile)

	// 用户注册
	g.POST(
		"/v1/me",
		newTracker(cs.ActionRegister),
		captchaValidate,
		// 限制相同IP在60秒之内只能调用5次
		newIPLimit(5, 60*time.Second, cs.ActionLogin),
		shouldAnonymous,
		ctrl.register,
	)
	// 刷新user session的ttl
	g.PATCH(
		"/v1/me",
		newTracker(cs.ActionUserMeUpdate),
		ctrl.updateMe,
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
		middleware.WaitFor(time.Second),
		newTracker(cs.ActionLogin),
		captchaValidate,
		shouldAnonymous,
		loginLimit,
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		// 限制10分钟内，相同的账号只允许出错5次
		newErrorLimit(5, 10*time.Minute, func(c *elton.Context) string {
			return gjson.GetBytes(c.RequestBody, "account").String()
		}),
		ctrl.login,
	)
	// 用户退出登录
	g.DELETE(
		"/v1/me",
		newTracker(cs.ActionLogout),
		shouldLogined,
		ctrl.logout,
	)

	// 获取客户登录记录
	g.GET(
		"/v1/login-records",
		shouldBeAdmin,
		ctrl.listLoginRecord,
	)
}

// toConditions get conditions of list user
func (params *listUserParams) toConditions() (conditions []interface{}) {
	queryList := make([]string, 0)
	args := make([]interface{}, 0)
	if params.Role != "" {
		queryList = append(queryList, "? = ANY(roles)")
		args = append(args, params.Role)
	}
	if params.Group != "" {
		queryList = append(queryList, "? = ANY(groups)")
		args = append(args, params.Group)
	}
	if params.Keyword != "" {
		queryList = append(queryList, "account ILIKE ?")
		args = append(args, "%"+params.Keyword+"%")
	}
	if params.Status != "" {
		queryList = append(queryList, "status = ?")
		args = append(args, params.Status)
	}
	conditions = make([]interface{}, 0)
	if len(queryList) != 0 {
		conditions = append(conditions, strings.Join(queryList, " AND "))
		conditions = append(conditions, args...)
	}
	return
}

// toConditions get conditions of list user login
func (params *listUserLoginRecordParams) toConditions() (conditions []interface{}) {
	queryList := make([]string, 0)
	args := make([]interface{}, 0)
	if params.Account != "" {
		queryList = append(queryList, "account = ?")
		args = append(args, params.Account)
	}
	if !params.Begin.IsZero() {
		queryList = append(queryList, "created_at >= ?")
		args = append(args, util.FormatTime(params.Begin))
	}
	if !params.End.IsZero() {
		queryList = append(queryList, "created_at <= ?")
		args = append(args, util.FormatTime(params.End))
	}
	conditions = make([]interface{}, 0)
	if len(queryList) != 0 {
		conditions = append(conditions, strings.Join(queryList, " AND "))
		conditions = append(conditions, args...)
	}
	return
}

// get user info from session
func pickUserInfo(c *elton.Context) (userInfo *userInfoResp) {
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
		userInfo.Groups = us.GetGroups()
		userInfo.Anonymous = false
	}
	return
}

// 用户信息
// swagger:response usersMeInfoResponse
// nolint
type usersMeInfoResponse struct {
	// in: body
	Body *userInfoResp
}

// swagger:route GET /users/v1/me users usersMe
// getUserInfo
//
// 获取用户信息，如果用户已登录，则返回用户相关信息
// responses:
// 	200: usersMeInfoResponse
func (ctrl userCtrl) me(c *elton.Context) (err error) {
	key := config.GetTrackKey()
	cookie, _ := c.Cookie(key)
	// ulid的长度为26
	if cookie == nil || len(cookie.Value) != 26 {
		uid := util.GenUlid()
		_ = c.AddCookie(&http.Cookie{
			Name:     key,
			Value:    uid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   365 * 24 * 3600,
		})
		trackRecord := &service.UserTrackRecord{
			UserAgent: c.GetRequestHeader("User-Agent"),
			IP:        c.RealIP(),
			TrackID:   util.GetTrackID(c),
		}
		_ = userSrv.AddTrackRecord(trackRecord, c)
	}
	c.Body = pickUserInfo(c)
	return
}

func (ctrl userCtrl) profile(c *elton.Context) (err error) {
	us := getUserSession(c)
	user, err := userSrv.FindOneByAccount(us.GetAccount())
	if err != nil {
		return
	}
	c.Body = user
	return
}

// 用户登录Token，用于客户登录密码加密
// swagger:response usersLoginTokenResponse
// nolint
type usersLoginTokenResponse struct {
	// in: body
	Body *loginTokenResp
}

// swagger:route GET /users/v1/me/login users usersLoginToken
// getLoginToken
//
// 获取用户登录Token
// responses:
// 	200: usersLoginTokenResponse
func (ctrl userCtrl) getLoginToken(c *elton.Context) (err error) {
	us := getUserSession(c)
	// 清除当前session id，确保每次登录的用户都是新的session
	us.ClearSessionID()
	token := util.RandomString(8)
	err = us.SetLoginToken(token)
	if err != nil {
		return
	}
	c.Body = &loginTokenResp{
		Token: token,
	}
	return
}

func omitUserInfo(u *service.User) {
	u.Password = ""
}

// 用户注册响应
// swagger:response usersRegisterResponse
// nolint
type usersRegisterResponse struct {
	// in: body
	Body *service.User
}

// swagger:parameters usersRegister usersMeLogin
// nolint
type usersRegisterParams struct {
	// in: body
	Payload *registerLoginUserParams
	// in: header
	Captcha string `json:"X-Captcha"`
}

// swagger:route POST /users/v1/me users usersRegister
// userRegister
//
// 用户注册，注册需要使用通用图形验证码，在成功时返回用户信息
// responses:
// 	201: usersRegisterResponse
func (ctrl userCtrl) register(c *elton.Context) (err error) {
	params := registerLoginUserParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	u := &service.User{
		Account:  params.Account,
		Password: params.Password,
	}
	err = userSrv.Add(u)
	if err != nil {
		return
	}
	omitUserInfo(u)
	c.Created(u)
	return
}

// 用户登录响应
// swagger:response usersLoginResponse
// nolint
type usersLoginResponse struct {
	// in: body
	Body *service.User
}

// swagger:route POST /users/v1/me/login users usersMeLogin
// userLogin
//
// 用户登录，需要使用通用图形验证码
// responses:
// 	200: usersLoginResponse
func (ctrl userCtrl) login(c *elton.Context) (err error) {
	params := registerLoginUserParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	token := us.GetLoginToken()
	if token == "" {
		err = errLoginTokenNil
		return
	}
	u, err := userSrv.Login(params.Account, params.Password, token)
	if err != nil {
		return
	}
	loginRecord := &service.UserLoginRecord{
		Account:       params.Account,
		UserAgent:     c.GetRequestHeader("User-Agent"),
		IP:            c.RealIP(),
		TrackID:       util.GetTrackID(c),
		SessionID:     util.GetSessionID(c),
		XForwardedFor: c.GetRequestHeader("X-Forwarded-For"),
	}
	_ = userSrv.AddLoginRecord(loginRecord, c)
	omitUserInfo(u)
	_ = us.SetAccount(u.Account)
	_ = us.SetRoles(u.Roles)
	_ = us.SetGroups(u.Groups)
	c.Body = u
	return
}

// logout user logout
func (ctrl userCtrl) logout(c *elton.Context) (err error) {
	us := getUserSession(c)
	if us != nil {
		err = us.Destroy()
	}
	c.NoContent()
	return
}

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

// refresh user refresh
func (ctrl userCtrl) updateMe(c *elton.Context) (err error) {
	// 如果没有数据要更新，如{}
	if len(c.RequestBody) <= 2 {
		return ctrl.refresh(c)
	}
	us := getUserSession(c)
	if us == nil {
		c.NoContent()
		return
	}
	params := updateMeParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	account := us.GetAccount()
	isUpdatedPassword := params.NewPassword != ""
	// 如果要更新密码，先校验旧密码是否一致
	if isUpdatedPassword {
		user, e := userSrv.FindOneByAccount(account)
		if e != nil {
			err = e
			return
		}
		// 如果密码不一致
		if user.Password != params.Password {
			err = hes.New("password is incorrect")
			return
		}
	}

	err = userSrv.UpdateByAccount(account, &service.User{
		Email:    params.Email,
		Mobile:   params.Mobile,
		Password: params.NewPassword,
	})
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// list user list
func (ctrl userCtrl) list(c *elton.Context) (err error) {
	params := listUserParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	args := params.toConditions()
	queryParams := params.listParams.toPGQueryParams()
	if queryParams.Offset == 0 {
		count, err = userSrv.Count(args...)
		if err != nil {
			return
		}
	}
	users, err := userSrv.List(queryParams, args...)
	if err != nil {
		return
	}
	c.Body = &struct {
		Count int             `json:"count"`
		Users []*service.User `json:"users"`
	}{
		count,
		users,
	}
	return
}

// update user update
func (ctrl userCtrl) updateByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := updateUserParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	// 只能su用户才可以添加su权限
	if util.ContainsString(params.Roles, cs.UserRoleSu) {
		roles := getUserSession(c).GetRoles()
		if !util.ContainsString(roles, cs.UserRoleSu) {
			err = hes.New("add su role is forbidden")
			return
		}
	}
	user := service.User{}
	if params.Status != 0 {
		user.Status = params.Status
	}
	if len(params.Roles) != 0 {
		user.Roles = pq.StringArray(params.Roles)
	}
	if len(params.Groups) != 0 {
		user.Groups = pq.StringArray(params.Groups)
	}
	err = userSrv.UpdateByID(uint(id), user)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// listLoginRecord list login record
func (ctrl userCtrl) listLoginRecord(c *elton.Context) (err error) {
	params := listUserLoginRecordParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	queryParams := params.listParams.toPGQueryParams()
	count := -1
	args := params.toConditions()
	if queryParams.Offset == 0 {
		count, err = userSrv.CountLoginRecord(args...)
		if err != nil {
			return
		}
	}
	result, err := userSrv.ListLoginRecord(queryParams, args...)
	if err != nil {
		return
	}

	c.Body = struct {
		Logins []*service.UserLoginRecord `json:"logins"`
		Count  int                        `json:"count"`
	}{
		result,
		count,
	}
	return
}
