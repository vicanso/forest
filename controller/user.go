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
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/tidwall/gjson"
	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/ent"
	"github.com/vicanso/forest/ent/predicate"
	"github.com/vicanso/forest/ent/schema"
	"github.com/vicanso/forest/ent/user"
	"github.com/vicanso/forest/ent/userlogin"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/tracer"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/forest/validate"
	"github.com/vicanso/hes"
)

type userCtrl struct{}

// 响应相关定义
type (
	// 用户登录Token响应
	userLoginTokenResp struct {
		// 用户登录的Token
		Token string `json:"token,omitempty"`
	}
	// userInfoResp 用户信息响应
	userInfoResp struct {
		// 服务器当前时间，2021-03-06T15:10:12+08:00
		Date string `json:"date,omitempty"`
		service.UserSessionInfo
	}

	// userListResp 用户列表响应
	userListResp struct {
		// 用户列表
		Users []*ent.User `json:"users,omitempty"`

		// 用户记录总数，如果返回-1表示此次查询未返回总数
		Count int `json:"count,omitempty"`
	}
	// userRoleListResp 用户角色列表响应
	userRoleListResp struct {
		UserRoles []*schema.UserRoleInfo `json:"userRoles,omitempty"`
	}
	// userLoginListResp 用户登录列表响应
	userLoginListResp struct {
		UserLogins []*ent.UserLogin `json:"userLogins,omitempty"`
		Count      int              `json:"count,omitempty"`
	}
)

// 参数相关定义
type (
	userListParams struct {
		listParams

		// 关键字搜索
		// validate:"omitempty,xKeyword"
		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`

		// 用户角色筛选
		// validate:"omitempty,xUserRole"
		Role string `json:"role,omitempty" validate:"omitempty,xUserRole"`

		// 用户分组筛选
		// validate:"omitempty,xUserGroup"
		Group string `json:"group,omitempty" validate:"omitempty,xUserGroup"`

		// 用户状态分组
		// validate:"omitempty,xStatus"
		Status string `json:"status,omitempty" validate:"omitempty,xStatus"`
	}

	// userLoginListParams 用户登录查询
	userLoginListParams struct {
		listParams

		Begin   time.Time `json:"begin,omitempty"`
		End     time.Time `json:"end,omitempty"`
		Account string    `json:"account,omitempty" validate:"omitempty,xUserAccount"`
	}

	// userRegisterLoginParams 注册与登录参数
	userRegisterLoginParams struct {
		// 账户
		Account string `json:"account,omitempty" validate:"required,xUserAccount"`
		// 用户密码，如果登录则是sha256(token + 用户密码)
		Password string `json:"password,omitempty" validate:"required,xUserPassword"`
	}

	// userUpdateMeParams 用户信息更新参数
	userUpdateMeParams struct {
		Name        string `json:"name,omitempty" validate:"omitempty,xUserName"`
		Email       string `json:"email,omitempty" validate:"omitempty,xUserEmail"`
		Password    string `json:"password,omitempty" validate:"omitempty,xUserPassword"`
		NewPassword string `json:"newPassword,omitempty" validate:"omitempty,xUserPassword"`
	}
	// userUpdateParams 更新用户信息参数
	userUpdateParams struct {
		Roles  []string      `json:"roles,omitempty" validate:"omitempty"`
		Status schema.Status `json:"status,omitempty" validate:"omitempty,xStatus"`
	}
	// userActionAddParams 用户添加行为记录的参数
	userActionAddParams struct {
		Actions []struct {
			// Category 用户行为类型
			Category string `json:"category,omitempty" validate:"required,xUserActionCategory"`
			// Route 触发时所在路由
			Route string `json:"route,omitempty" validate:"required,xUserActionRoute"`
			// Path 触发时的完整路径
			Path string `json:"path,omitempty" validate:"required,xPath"`
			// Result 操作结果，0:成功 1:失败
			Result int `json:"result,omitempty"`
			// Time 记录的时间戳，单位秒
			Time int64 `json:"time,omitempty" validate:"required"`
			// Extra 其它额外信息
			Extra map[string]interface{} `json:"extra,omitempty"`
		} `json:"actions,omitempty" validate:"required,dive"`
	}
)

var (
	// session配置信息
	sessionConfig config.SessionConfig
)

const (
	errUserCategory = "user"
)

func init() {
	sessionConfig = config.GetSessionConfig()
	prefix := "/users"
	g := router.NewGroup(prefix, loadUserSession)
	noneSessionGroup := router.NewGroup(prefix)

	ctrl := userCtrl{}

	// 获取用户列表
	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

	// 获取用户信息
	g.GET(
		"/v1/{id}",
		shouldBeAdmin,
		ctrl.findByID,
	)

	// 更新用户信息
	g.PATCH(
		"/v1/{id}",
		newTrackerMiddleware(cs.ActionUserInfoUpdate),
		shouldBeAdmin,
		ctrl.updateByID,
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

	// 获取用户信息
	g.GET(
		"/v1/detail",
		shouldBeLogin,
		ctrl.detail,
	)

	// 用户注册
	g.POST(
		"/v1/me",
		// 注册无论成功失败都最少等待1秒
		middleware.WaitFor(time.Second),
		newTrackerMiddleware(cs.ActionRegister),
		captchaValidate,
		// 限制相同IP在60秒之内只能调用5次
		newIPLimit(5, 60*time.Second, cs.ActionRegister),
		shouldBeAnonymous,
		ctrl.register,
	)

	// 用户登录
	g.POST(
		"/v1/me/login",
		// 登录如果失败则最少等待1秒
		middleware.WaitFor(time.Second, true),
		newTrackerMiddleware(cs.ActionLogin),
		captchaValidate,
		shouldBeAnonymous,
		// 同一个账号限制3秒只能登录一次（无论成功还是失败）
		newConcurrentLimit([]string{
			"account",
		}, 3*time.Second, cs.ActionLogin),
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		// 限制10分钟内，相同的账号只允许出错5次
		newErrorLimit(5, 10*time.Minute, func(c *elton.Context) string {
			return gjson.GetBytes(c.RequestBody, "account").String()
		}),
		ctrl.login,
	)
	// 内部登录
	g.POST(
		"/inner/v1/me/login",
		newTrackerMiddleware(cs.ActionLogin),
		captchaValidate,
		isIntranet,
		ctrl.login,
	)

	// 刷新user session的ttl或更新客户信息
	g.PATCH(
		"/v1/me",
		newTrackerMiddleware(cs.ActionUserMeUpdate),
		shouldBeLogin,
		ctrl.updateMe,
	)

	// 用户退出登录
	g.DELETE(
		"/v1/me",
		newTrackerMiddleware(cs.ActionLogout),
		shouldBeLogin,
		ctrl.logout,
	)

	// 获取客户登录记录
	g.GET(
		"/v1/login-records",
		shouldBeAdmin,
		ctrl.listLoginRecord,
	)

	// 添加用户行为
	g.POST(
		"/v1/actions",
		ctrl.addUserAction,
	)

	// 获取用户角色分组
	noneSessionGroup.GET(
		"/v1/roles",
		noCacheIfRequestNoCache,
		ctrl.getRoleList,
	)
}

// validateBeforeSave 保存前校验
func (params *userRegisterLoginParams) validateBeforeSave(ctx context.Context) (err error) {
	// 判断该账户是否已注册
	exists, err := getEntClient().User.Query().
		Where(user.Account(params.Account)).
		Exist(ctx)
	if err != nil {
		return
	}
	if exists {
		err = hes.New("该账户已注册", errUserCategory)
		return
	}

	return
}

// save 创建用户
func (params *userRegisterLoginParams) save(ctx context.Context) (*ent.User, error) {
	err := params.validateBeforeSave(ctx)
	if err != nil {
		return nil, err
	}
	return getEntClient().User.Create().
		SetAccount(params.Account).
		SetPassword(params.Password).
		Save(ctx)
}

// login 登录
func (params *userRegisterLoginParams) login(ctx context.Context, token string) (u *ent.User, err error) {
	u, err = getEntClient().User.Query().
		Where(user.Account(params.Account)).
		First(ctx)
	errAccountOrPasswordInvalid := hes.New("账户或者密码错误", errUserCategory)
	if err != nil {
		// 如果登录时账号不存在
		if ent.IsNotFound(err) {
			err = errAccountOrPasswordInvalid
		}
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
	// 禁止非正常状态用户登录
	if u.Status != schema.StatusEnabled {
		err = hes.NewWithStatusCode("该账户不允许登录", http.StatusForbidden, errUserCategory)
		return
	}
	return
}

// update 更新用户信息
func (params *userUpdateMeParams) updateOneAccount(ctx context.Context, account string) (u *ent.User, err error) {

	u, err = getEntClient().User.Query().
		Where(user.Account(account)).
		First(ctx)
	if err != nil {
		return
	}
	// 更新密码时需要先校验旧密码
	if params.NewPassword != "" {
		if u.Password != params.Password {
			err = hes.New("旧密码错误，请重新输入", errUserCategory)
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

// updateByID 通过ID更新信息
func (params *userUpdateParams) updateByID(ctx context.Context, id int) (u *ent.User, err error) {
	updateOne := getEntClient().User.UpdateOneID(id)
	if len(params.Roles) != 0 {
		updateOne = updateOne.SetRoles(params.Roles)
	}
	if params.Status != 0 {
		updateOne = updateOne.SetStatus(params.Status)
	}
	return updateOne.Save(ctx)
}

// where 将查询条件中的参数转换为对应的where条件
func (params *userListParams) where(query *ent.UserQuery) *ent.UserQuery {
	if params.Keyword != "" {
		query = query.Where(user.AccountContains(params.Keyword))
	}
	if params.Role != "" {
		query = query.Where(predicate.User(func(s *sql.Selector) {
			s.Where(sqljson.ValueContains(user.FieldRoles, params.Role))
		}))

	}
	if params.Status != "" {
		v, _ := strconv.Atoi(params.Status)
		query = query.Where(user.Status(schema.Status(v)))
	}
	return query
}

// queryAll 查询用户列表
func (params *userListParams) queryAll(ctx context.Context) (users []*ent.User, err error) {
	query := getEntClient().User.Query()

	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)

	return query.All(ctx)
}

// count 计算总数
func (params *userListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().User.Query()

	query = params.where(query)

	return query.Count(ctx)
}

// where 登录记录的where筛选
func (params *userLoginListParams) where(query *ent.UserLoginQuery) *ent.UserLoginQuery {
	if params.Account != "" {
		query = query.Where(userlogin.AccountEQ(params.Account))
	}
	query = query.Where(userlogin.CreatedAtGTE(params.Begin))
	query = query.Where(userlogin.CreatedAtLTE(params.End))
	return query
}

// queryAll 查询所有的登录记录
func (params *userLoginListParams) queryAll(ctx context.Context) (userLogins []*ent.UserLogin, err error) {
	query := getEntClient().UserLogin.Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)
	return query.All(ctx)
}

// count 计算登录记录总数
func (params *userLoginListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().UserLogin.Query()
	query = params.where(query)
	return query.Count(ctx)
}

// pickUserInfo 获取用户信息
func pickUserInfo(c *elton.Context) (resp userInfoResp, err error) {
	us := getUserSession(c)
	userInfo, err := us.GetInfo()
	if err != nil {
		return
	}
	resp = userInfoResp{
		Date: now(),
	}
	resp.UserSessionInfo = userInfo
	return
}

// swagger:route GET /users/v1 users userList
// 用户查询
//
// 根据查询条件返回用户列表，限制只允许管理人查询
// Responses:
// 	200: apiUserListResponse

func (*userCtrl) list(c *elton.Context) (err error) {
	params := userListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	users, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.Body = &userListResp{
		Count: count,
		Users: users,
	}

	return
}

// findByID 通过ID查询用户信息
func (*userCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	data, err := getEntClient().User.Get(c.Context(), id)
	if err != nil {
		return
	}
	c.Body = data
	return
}

// updateByID 更新信息
func (ctrl *userCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := userUpdateParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	user, err := params.updateByID(c.Context(), id)
	if err != nil {
		return
	}
	c.Body = user
	return
}

// swagger:route GET /users/v1/me/login users userLoginToken
// 获取登录的token
//
// 在登录之前需要先调用获取token，此token用于登录时与客户密码sha256生成hash，
// 保证客户每次登录时的密码均不相同，避免接口重放登录。
// Responses:
// 	200: apiUserLoginTokenResponse
func (*userCtrl) getLoginToken(c *elton.Context) (err error) {
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
	c.Body = &userLoginTokenResp{
		Token: userInfo.Token,
	}
	return
}

// swagger:route GET /users/v1/me users userMe
// 获取客户信息
//
// 若用户登录状态，则返回客户的相关信息。
// 若用户未登录，仅返回服务器当前时间。
// Responses:
// 	200: apiUserInfoResponse
func (*userCtrl) me(c *elton.Context) (err error) {
	cookie, _ := c.Cookie(sessionConfig.TrackKey)
	if cookie == nil {
		uid := util.GenXID()
		c.AddCookie(&http.Cookie{
			Name:     sessionConfig.TrackKey,
			Value:    uid,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   365 * 24 * 3600,
		})

		ip := c.RealIP()
		fields := map[string]interface{}{
			cs.FieldUserAgent: c.GetRequestHeader("User-Agent"),
			cs.FieldTID:       uid,
			cs.FieldIP:        ip,
		}

		// 记录创建user track
		tracerInfo := tracer.GetTracerInfo()
		go func() {
			tracer.SetTracerInfo(tracerInfo)
			location, _ := service.GetLocationByIP(ip, nil)
			if location.IP != "" {
				fields[cs.FieldCountry] = location.Country
				fields[cs.FieldProvince] = location.Province
				fields[cs.FieldCity] = location.City
				fields[cs.FieldISP] = location.ISP
			}
			getInfluxSrv().Write(cs.MeasurementUserAddTrack, nil, fields)
		}()
	}
	resp, err := pickUserInfo(c)
	if err != nil {
		return
	}
	c.Body = &resp
	return
}

// detail 详细信息
func (*userCtrl) detail(c *elton.Context) (err error) {
	us := getUserSession(c)
	user, err := getEntClient().User.Get(c.Context(), us.MustGetInfo().ID)
	if err != nil {
		return
	}

	c.Body = user
	return
}

// swagger:route POST /users/v1/me users userRegister
// 用户注册
//
// 用户注册时提交的密码需要在客户端以sha256后提交，
// 在成功注册后返回用户信息。
// 需注意此时并非登录状态，需要客户自主登录。
// Responses:
// 	201: apiUserRegisterResponse
func (*userCtrl) register(c *elton.Context) (err error) {
	params := userRegisterLoginParams{}
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
			_, _ = user.Update().
				SetRoles([]string{
					schema.UserRoleSu,
				}).
				Save(context.Background())
		}()
	}
	c.Created(user)
	return
}

// swagger:route POST /users/v1/me/login users userLogin
// 用户登录
//
// 用户登录时需要先获取token，之后使用token与密码sha256后提交，
// 登录成功后返回用户信息。
// Responses:
// 	200: apiUserInfoResponse
func (*userCtrl) login(c *elton.Context) (err error) {
	params := userRegisterLoginParams{}
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
		err = hes.New("登录令牌不能为空", errUserCategory)
		return
	}
	// 登录
	u, err := params.login(c.Context(), userInfo.Token)
	if err != nil {
		return
	}
	account := u.Account

	// 设置session
	err = us.SetInfo(service.UserSessionInfo{
		Account: account,
		ID:      u.ID,
		Roles:   u.Roles,
		// Groups: u.,
	})
	if err != nil {
		return
	}

	ip := c.RealIP()
	tid := util.GetTrackID(c)
	sid := util.GetSessionID(c)
	userAgent := c.GetRequestHeader("User-Agent")

	xForwardedFor := c.GetRequestHeader("X-Forwarded-For")
	tracerInfo := tracer.GetTracerInfo()
	go func() {
		tracer.SetTracerInfo(tracerInfo)
		fields := map[string]interface{}{
			cs.FieldAccount:   account,
			cs.FieldUserAgent: userAgent,
			cs.FieldIP:        ip,
			cs.FieldTID:       tid,
			cs.FieldSID:       sid,
		}
		location, _ := service.GetLocationByIP(ip, nil)
		country := ""
		province := ""
		city := ""
		isp := ""
		if location.IP != "" {
			country = location.Country
			province = location.Province
			city = location.City
			isp = location.ISP
			fields[cs.FieldCountry] = country
			fields[cs.FieldProvince] = province
			fields[cs.FieldCity] = city
			fields[cs.FieldISP] = isp
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// 记录至数据库
		_, err := getEntClient().UserLogin.Create().
			SetAccount(account).
			SetUserAgent(userAgent).
			SetIP(ip).
			SetTrackID(tid).
			SetSessionID(sid).
			SetXForwardedFor(xForwardedFor).
			SetCountry(country).
			SetProvince(province).
			SetCity(city).
			SetIsp(isp).
			Save(ctx)
		if err != nil {
			log.Default().Error().
				Err(err).
				Msg("save user login fail")
		}
		// 记录用户登录行为
		getInfluxSrv().Write(cs.MeasurementUserLogin, nil, fields)
	}()

	// 返回用户信息
	resp, err := pickUserInfo(c)
	if err != nil {
		return
	}
	c.Body = &resp
	return
}

// swagger:route DELETE /users/v1/me users userLogout
// 用户退出登录
//
// 退出用户当前登录状态，成功时并无内容返回。
// Responses:
// 	204: apiNoContentResponse
func (*userCtrl) logout(c *elton.Context) (err error) {
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
func (*userCtrl) refresh(c *elton.Context) (err error) {
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
func (ctrl *userCtrl) updateMe(c *elton.Context) (err error) {
	// 如果没有数据要更新，如{}
	if len(c.RequestBody) <= 2 {
		return ctrl.refresh(c)
	}
	us := getUserSession(c)
	params := userUpdateMeParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}

	// 更新用户信息
	_, err = params.updateOneAccount(c.Context(), us.MustGetInfo().Account)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// getRoleList 获取用户角色列表
func (*userCtrl) getRoleList(c *elton.Context) (err error) {
	c.CacheMaxAge(time.Minute)
	c.Body = &userRoleListResp{
		UserRoles: schema.GetUserRoleList(),
	}
	return
}

// listLoginRecord list login record
func (ctrl userCtrl) listLoginRecord(c *elton.Context) (err error) {
	params := userLoginListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	if params.ShouldCount() {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	userLogins, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.Body = &userLoginListResp{
		Count:      count,
		UserLogins: userLogins,
	}
	return
}

// addUserAction add user action
func (ctrl userCtrl) addUserAction(c *elton.Context) (err error) {
	tid := util.GetTrackID(c)
	// 如果没有tid，则直接返回
	if tid == "" {
		c.NoContent()
		return
	}

	params := userActionAddParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	now := time.Now().Unix()
	us := getUserSession(c)
	account := ""
	if us.IsLogin() {
		account = us.MustGetInfo().Account
	}

	count := 0
	for _, item := range params.Actions {
		// 如果时间大于当前时间或者一天前，则忽略
		if item.Time > now || item.Time < now-24*3600 {
			continue
		}
		count++
		// 由于客户端的统计时间精度只到second
		// 随机生成nano second填充
		nsec := rand.Int() % int(time.Second)
		t := time.Unix(item.Time, int64(nsec))
		fields := map[string]interface{}{
			cs.FieldRoute: item.Route,
			cs.FieldPath:  item.Path,
			cs.FieldTID:   tid,
		}
		if account != "" {
			fields[cs.FieldAccount] = account
		}
		fields = util.MergeMapStringInterface(fields, item.Extra)
		getInfluxSrv().Write(cs.MeasurementUserAction, map[string]string{
			cs.TagCategory: item.Category,
			cs.TagResult:   strconv.Itoa(item.Result),
		}, fields, t)
	}
	c.Body = map[string]int{
		"count": count,
	}
	return
}
