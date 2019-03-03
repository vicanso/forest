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
	// userSession user session struct
	userSession struct {
		se *session.Session
	}
	// User user
	User struct{}
)

// GetAccount get the account
func (u *userSession) GetAccount() string {
	if u.se == nil {
		return ""
	}
	return u.se.GetString(UserAccount)
}

// SetAccount set the account
func (u *userSession) SetAccount(account string) error {
	return u.se.Set(UserAccount, account)
}

// GetUpdatedAt get updated at
func (u *userSession) GetUpdatedAt() string {
	return u.se.GetUpdatedAt()
}

// SetLoginAt set the login at
func (u *userSession) SetLoginAt(date string) error {
	return u.se.Set(UserLoginAt, date)
}

// GetLoginAt get login at
func (u *userSession) GetLoginAt() string {
	return u.se.GetString(UserLoginAt)
}

// SetLoginToken get user login token
func (u *userSession) SetLoginToken(token string) error {
	return u.se.Set(UserLoginToken, token)
}

// GetLoginToken get user login token
func (u *userSession) GetLoginToken() string {
	return u.se.GetString(UserLoginToken)
}

// Destroy destroy user session
func (u *userSession) Destroy() error {
	return u.se.Destroy()
}

// Refresh refresh user session
func (u *userSession) Refresh() error {
	return u.se.Refresh()
}

// ClearSessionID clear session id
func (u *userSession) ClearSessionID() {
	u.se.ID = ""
}

// NewUserSession create a user session
func NewUserSession(c *cod.Context) *userSession {
	v := c.Get(session.Key)
	if v == nil {
		return nil
	}
	return &userSession{
		se: v.(*session.Session),
	}
}
