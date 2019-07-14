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

package service

import (
	"github.com/vicanso/cod"
	session "github.com/vicanso/cod-session"
)

const (
	// UserAccount user account field
	UserAccount = "account"
	// UserLoginAt user login at
	UserLoginAt = "loginAt"
	// UserRoles user roles
	UserRoles = "roles"
	// UserLoginToken user login token
	UserLoginToken = "loginToken"
)

type (
	// UserSession user session struct
	UserSession struct {
		se *session.Session
	}
	// User user
	User struct{}
)

// GetAccount get the account
func (u *UserSession) GetAccount() string {
	if u.se == nil {
		return ""
	}
	return u.se.GetString(UserAccount)
}

// SetAccount set the account
func (u *UserSession) SetAccount(account string) error {
	return u.se.Set(UserAccount, account)
}

// GetUpdatedAt get updated at
func (u *UserSession) GetUpdatedAt() string {
	return u.se.GetUpdatedAt()
}

// SetLoginAt set the login at
func (u *UserSession) SetLoginAt(date string) error {
	return u.se.Set(UserLoginAt, date)
}

// GetLoginAt get login at
func (u *UserSession) GetLoginAt() string {
	return u.se.GetString(UserLoginAt)
}

// SetLoginToken get user login token
func (u *UserSession) SetLoginToken(token string) error {
	return u.se.Set(UserLoginToken, token)
}

// GetLoginToken get user login token
func (u *UserSession) GetLoginToken() string {
	return u.se.GetString(UserLoginToken)
}

// Destroy destroy user session
func (u *UserSession) Destroy() error {
	return u.se.Destroy()
}

// Refresh refresh user session
func (u *UserSession) Refresh() error {
	return u.se.Refresh()
}

// ClearSessionID clear session id
func (u *UserSession) ClearSessionID() {
	u.se.ID = ""
}

// NewUserSession create a user session
func NewUserSession(c *cod.Context) *UserSession {
	v := c.Get(session.Key)
	if v == nil {
		return nil
	}
	return &UserSession{
		se: v.(*session.Session),
	}
}
