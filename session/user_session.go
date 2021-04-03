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

package session

import (
	"encoding/json"

	"github.com/vicanso/elton"
	se "github.com/vicanso/elton-session"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/util"
)

const (
	// UserSessionInfoKey user session info
	UserSessionInfoKey = "user-session-info"
)

type (
	// UserSessionInfo 用户session中的信息
	UserSessionInfo struct {
		// 登录时使用的Token，此字段不返回
		Token string `json:"token,omitempty"`
		// 用户账号
		Account string `json:"account,omitempty"`
		// 用户ID
		ID int `json:"id,omitempty"`
		// 用户角色列表
		Roles []string `json:"roles,omitempty"`
		// 用户分组列表
		Groups []string `json:"groups,omitempty"`
		// Session信息更新时间
		UpdatedAt string `json:"updatedAt,omitempty"`
		// Session信息创建时间
		LoginAt string `json:"loginAt,omitempty"`
	}
	// UserSession 用户session
	UserSession struct {
		unmarshalDone bool
		se            *se.Session
		info          UserSessionInfo
	}
)

// GetUserInfo 获取用户信息
func (us *UserSession) GetInfo() (info UserSessionInfo, err error) {
	if us.unmarshalDone {
		info = us.info
		return
	}
	data := us.se.GetString(UserSessionInfoKey)
	if data == "" {
		data = "{}"
	}
	info = UserSessionInfo{}
	err = json.Unmarshal([]byte(data), &info)
	if err != nil {
		return
	}
	us.info = info
	us.unmarshalDone = true
	return
}

// MustGetInfo 获取用户信息，如果信息获取失败则触发panic，
// 如果前置中间件已保证是登录状态，可以使用此函数，否则禁止使用
func (us *UserSession) MustGetInfo() (info UserSessionInfo) {
	info, err := us.GetInfo()
	if err != nil {
		panic(err)
	}
	return info
}

// IsLogin 判断用户是否已登录
func (us *UserSession) IsLogin() bool {
	info, err := us.GetInfo()
	if err != nil {
		return false
	}
	return info.Account != ""
}

// SetInfo 设置用户信息
func (us *UserSession) SetInfo(info UserSessionInfo) (err error) {
	// 登录时设置登录时间
	if info.Account != "" && info.LoginAt == "" {
		info.LoginAt = util.NowString()
	}
	info.UpdatedAt = util.NowString()
	us.info = info
	us.unmarshalDone = true
	buf, err := json.Marshal(&info)
	if err != nil {
		return
	}
	err = us.se.Set(UserSessionInfoKey, string(buf))
	if err != nil {
		return
	}
	return
}

// Destroy 清除用户session
func (us *UserSession) Destroy() error {
	return us.se.Destroy()
}

// Refresh 刷新用户session ttl
func (us *UserSession) Refresh() error {
	return us.se.Refresh()
}

// NewUserSession 创建新的用户session对象
func NewUserSession(c *elton.Context) *UserSession {
	if data, ok := c.Get(cs.UserSession); ok {
		us, ok := data.(*UserSession)
		if ok {
			return us
		}
	}
	v, ok := c.Get(se.Key)
	if !ok {
		return nil
	}
	us := &UserSession{
		se: v.(*se.Session),
	}
	c.Set(cs.UserSession, us)

	return us
}
