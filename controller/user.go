package controller

import (
	"net/http"
	"time"

	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/validate"
	"github.com/vicanso/hes"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/util"
)

type (
	userCtrl struct{}
	// UserInfoResp user info response
	UserInfoResp struct {
		Anonymous bool   `json:"anonymous,omitempty"`
		Account   string `json:"account,omitempty"`
		Date      string `json:"date,omitempty"`
		UpdatedAt string `json:"updatedAt,omitempty"`
		IP        string `json:"ip,omitempty"`
		TrackID   string `json:"trackId,omitempty"`
		LoginAt   string `json:"loginAt,omitempty"`
	}

	// UserLoginParams user login params
	UserLoginParams struct {
		Account  string `valid:"ascii,runelength(4|10)"`
		Password string `valid:"runelength(6|64)"`
	}
)

var (
	errLoginTokenNil = hes.New("login token is nil")
)

func init() {
	g := router.NewGroup("/users")
	ctrl := userCtrl{}

	g.GET(
		"/v1/me",
		userSession,
		ctrl.me,
	)

	// 用户登录
	g.GET(
		"/v1/me/login",
		userSession,
		isAnonymous,
		ctrl.getLoginToken,
	)
	// 限制3秒只能登录一次（无论成功还是失败）
	loginLimit := createConcurrentLimit([]string{
		"account",
	}, 3*time.Second, cs.ActionLogin)
	g.POST(
		"/v1/me/login",
		userSession,
		isAnonymous,
		loginLimit,
		// 限制相同IP在60秒之内只能调用10次
		newIPLimit(10, 60*time.Second, cs.ActionLogin),
		ctrl.login,
	)

	// 退出登录
	g.DELETE(
		"/v1/me/logout",
		userSession,
		ctrl.logout,
	)
}

// get user info from session
func (ctrl userCtrl) pickUserInfo(c *cod.Context) (userInfo *UserInfoResp) {
	us := getUserSession(c)
	userInfo = &UserInfoResp{
		Anonymous: true,
		Date:      now(),
		IP:        c.RealIP(),
		TrackID:   getTrackID(c),
	}
	if us == nil {
		return
	}
	account := us.GetAccount()
	if account != "" {
		userInfo.Account = account
		userInfo.Anonymous = false
		userInfo.UpdatedAt = us.GetUpdatedAt()
		userInfo.LoginAt = us.GetLoginAt()
	}
	return
}

// me get info of user
func (ctrl userCtrl) me(c *cod.Context) (err error) {
	key := config.GetTrackKey()
	// 如果没有track cookie，则生成
	if cookie, _ := c.SignedCookie(key); cookie == nil {
		c.AddSignedCookie(&http.Cookie{
			Name:     key,
			Value:    util.GenUlid(),
			Path:     "/",
			HttpOnly: true,
			// 有效期一年
			MaxAge: 365 * 24 * 3600,
		})
	}
	c.Body = ctrl.pickUserInfo(c)
	return
}

// login user login
func (ctrl userCtrl) login(c *cod.Context) (err error) {
	us := getUserSession(c)
	token := us.GetLoginToken()
	if token == "" {
		err = errLoginTokenNil
		return
	}
	params := &UserLoginParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	// TODO 从数据库读取客户密码与token sha1再校验
	us.SetAccount(params.Account)
	us.SetLoginAt(now())
	c.Body = ctrl.pickUserInfo(c)
	return
}

// getLoginToken get login token
func (ctrl userCtrl) getLoginToken(c *cod.Context) (err error) {
	us := getUserSession(c)
	us.ClearSessionID()
	token := util.RandomString(8)
	err = us.SetLoginToken(token)
	if err != nil {
		return
	}
	c.Body = &struct {
		Token string `json:"token"`
	}{
		token,
	}
	return
}

// logout logout
func (ctrl userCtrl) logout(c *cod.Context) (err error) {
	us := getUserSession(c)
	if us != nil {
		err = us.Destroy()
	}
	c.NoContent()
	return
}
