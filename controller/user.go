package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vicanso/forest/config"

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
)

func init() {
	g := router.NewGroup("/users")
	ctrl := userCtrl{}

	// user login
	// 限制3秒只能登录一次（无论成功还是失败）
	loginLimit := createConcurrentLimit([]string{
		"account",
	}, 3*time.Second, cs.ActionLogin)

	g.GET("/v1/me", ctrl.me)
	g.POST("/v1/me/login", loginLimit, ctrl.login)
}

// get user info from session
func (ctrl userCtrl) pickUserInfo(c *cod.Context) (userInfo *UserInfoResp) {
	userInfo = &UserInfoResp{
		Account: "tree.xie",
	}
	// us := getUserSession(c)
	// userInfo = &UserInfoResp{
	// 	Anonymous: true,
	// 	Date:      now(),
	// 	IP:        c.RealIP(),
	// 	TrackID:   getTrackID(c),
	// }
	// if us == nil {
	// 	return
	// }
	// account := us.GetAccount()
	// if account != "" {
	// 	userInfo.Account = account
	// 	userInfo.Anonymous = false
	// 	userInfo.UpdatedAt = us.GetUpdatedAt()
	// 	userInfo.LoginAt = us.GetLoginAt()
	// }
	return
}

// me get info of user
func (ctrl userCtrl) me(c *cod.Context) (err error) {
	key := config.GetTrackKey()
	// 如果没有track cookie，则生成
	if cookie, _ := c.Cookie(key); cookie == nil {
		c.AddCookie(&http.Cookie{
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
	fmt.Println(c.RequestBody)
	c.Body = &struct {
		Name string `json:"name"`
	}{
		"tree.xie",
	}
	return
}
