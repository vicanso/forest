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

package service

import (
	"encoding/json"

	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
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
		Token     string   `json:"token,omitempty"`
		Account   string   `json:"account,omitempty"`
		ID        int      `json:"id,omitempty"`
		Roles     []string `json:"roles,omitempty"`
		Groups    []string `json:"groups,omitempty"`
		UpdatedAt string   `json:"updatedAt,omitempty"`
		LoginedAt string   `json:"loginedAt,omitempty"`
	}
	// UserSession 用户session
	UserSession struct {
		unmarshalDone bool
		se            *session.Session
		info          UserSessionInfo
	}
)

// GetUesrInfo 获取用户信息
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

// MustGetInfo 获取用户信息，如果信息获取失败则触发panic
func (us *UserSession) MustGetInfo() (info UserSessionInfo) {
	info, err := us.GetInfo()
	if err != nil {
		panic(err)
	}
	return info
}

// IsLogined 判断用户是否已登录
func (us *UserSession) IsLogined() bool {
	info, err := us.GetInfo()
	if err != nil {
		return false
	}
	return info.Account != ""
}

// SetInfo 设置用户信息
func (us *UserSession) SetInfo(info UserSessionInfo) (err error) {
	if info.Account != "" && info.LoginedAt == "" {
		info.LoginedAt = util.NowString()
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

// Refresh 刷新用户sesion ttl
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
	v, ok := c.Get(session.Key)
	if !ok {
		return nil
	}
	us := &UserSession{
		se: v.(*session.Session),
	}
	c.Set(cs.UserSession, us)

	return us
}
