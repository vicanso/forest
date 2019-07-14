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

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/util"
)

type (
	userCtrl struct{}
	// userInfoResp user info response
	userInfoResp struct {
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
	g := router.NewGroup("/users", loadUserSession)
	ctrl := userCtrl{}

	g.GET("/v1/me", ctrl.me)
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
