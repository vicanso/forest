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
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/hes"

	"go.uber.org/zap"
)

const (
	// UserAccount user account field
	UserAccount = "account"
	// UserLoginAt user login at
	UserLoginAt = "loginAt"
	// UserRoles user roles
	UserRoles = "roles"
	// UserGroups user groups
	UserGroups = "groups"
	// UserLoginToken user login token
	UserLoginToken = "loginToken"
)

var (
	errAccountOrPasswordInvalid = hes.New("account or password is invalid")
)

type (
	// UserSession user session struct
	UserSession struct {
		se *session.Session
	}
	// User user
	User struct {
		helper.Model

		Account  string         `json:"account" gorm:"type:varchar(20);not null;unique_index:idx_users_account"`
		Password string         `json:"-" gorm:"type:varchar(128);not null"`
		Roles    pq.StringArray `json:"roles" gorm:"type:text[]"`
		Groups   pq.StringArray `json:"groups" gorm:"type:text[]"`
		// 用户状态
		Status int    `json:"status"`
		Email  string `json:"email"`
		Mobile string `json:"mobile"`
	}
	// UserRole user role
	UserRole struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	// UserStatus user status
	UserStatus struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	// UserGroup user group
	UserGroup struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	// UserLoginRecord user login
	UserLoginRecord struct {
		helper.Model

		Account       string `json:"account" gorm:"type:varchar(20);not null;index:idx_user_logins_account"`
		UserAgent     string `json:"userAgent"`
		IP            string `json:"ip" gorm:"type:varchar(64);not null"`
		TrackID       string `json:"trackId" gorm:"type:varchar(64);not null"`
		SessionID     string `json:"sessionId" gorm:"type:varchar(64);not null"`
		XForwardedFor string `json:"xForwardedFor" gorm:"type:varchar(128)"`
		Country       string `json:"country" gorm:"type:varchar(64)"`
		Province      string `json:"province" gorm:"type:varchar(64)"`
		City          string `json:"city" gorm:"type:varchar(64)"`
		ISP           string `json:"isp" gorm:"type:varchar(64)"`
	}
	// UserTrackRecord user track record
	UserTrackRecord struct {
		helper.Model

		TrackID   string `json:"trackId" gorm:"type:varchar(64);not null;index:idx_user_track_id"`
		UserAgent string `json:"userAgent"`
		IP        string `json:"ip" gorm:"type:varchar(64);not null"`
		Country   string `json:"country" gorm:"type:varchar(64)"`
		Province  string `json:"province" gorm:"type:varchar(64)"`
		City      string `json:"city" gorm:"type:varchar(64)"`
		ISP       string `json:"isp" gorm:"type:varchar(64)"`
	}
	// UserSrv user service
	UserSrv struct {
	}
)

func init() {
	pgGetClient().AutoMigrate(&User{}).
		AutoMigrate(&UserLoginRecord{}).
		AutoMigrate(&UserTrackRecord{})
}

// ListRoles list all user roles
func (srv *UserSrv) ListRoles() []*UserRole {
	return []*UserRole{
		&UserRole{
			Name:  "普通用户",
			Value: cs.UserRoleNormal,
		},
		&UserRole{
			Name:  "管理员",
			Value: cs.UserRoleAdmin,
		},
		&UserRole{
			Name:  "超级用户",
			Value: cs.UserRoleSu,
		},
	}
}

// ListStatuses list all user status
func (srv *UserSrv) ListStatuses() []*UserStatus {
	return []*UserStatus{
		&UserStatus{
			Name:  "正常",
			Value: cs.AccountStatusEnabled,
		},
		&UserStatus{
			Name:  "禁用",
			Value: cs.AccountStatusForbidden,
		},
	}
}

// ListGroups list all user group
func (srv *UserSrv) ListGroups() []*UserGroup {
	return []*UserGroup{
		&UserGroup{
			Name:  "IT",
			Value: cs.UserGroupIT,
		},
		&UserGroup{
			Name:  "财务",
			Value: cs.UserGroupFinance,
		},
	}
}

// createByID create a user model by id
func (srv *UserSrv) createByID(id uint) *User {
	u := &User{}
	u.Model.ID = id
	return u
}

// createLoginRecordByID cerate login record by id
func (srv *UserSrv) createLoginRecordByID(id uint) *UserLoginRecord {
	ulr := &UserLoginRecord{}
	ulr.Model.ID = id
	return ulr
}

// Add add user
func (srv *UserSrv) Add(u *User) (err error) {
	if u.Status == 0 {
		u.Status = cs.AccountStatusEnabled
	}
	if len(u.Roles) == 0 {
		u.Roles = pq.StringArray([]string{
			cs.UserRoleNormal,
		})
	}
	err = pgCreate(u)
	// 首次创建账号，设置su权限
	if u.ID == 1 {
		_ = srv.UpdateByID(u.ID, User{
			Roles: []string{
				cs.UserRoleSu,
			},
		})
	}
	return
}

// Login user login
func (srv *UserSrv) Login(account, password, token string) (u *User, err error) {
	u = &User{}
	err = pgGetClient().Where("account = ?", account).First(u).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = errAccountOrPasswordInvalid
		}
		return
	}
	pwd := util.Sha256(u.Password + token)
	// 用于自动化测试使用
	if util.IsDevelopment() && password == "fEqNCco3Yq9h5ZUglD3CZJT4lBsfEqNCco31Yq9h5ZUB" {
		pwd = password
	}
	if pwd != password {
		err = errAccountOrPasswordInvalid
		return
	}
	return
}

// UpdateByID update user by id
func (srv *UserSrv) UpdateByID(id uint, value interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Updates(value).Error
	return
}

// UpdateByAccount update user by account
func (srv *UserSrv) UpdateByAccount(account string, value interface{}) (err error) {
	err = pgGetClient().Model(&User{}).Where("account = ?", account).Updates(value).Error
	return
}

// FindOneByAccount find one by account
func (srv *UserSrv) FindOneByAccount(account string) (user *User, err error) {
	user = &User{}
	err = pgGetClient().Where("account = ?", account).First(user).Error
	return
}

// UpdateLoginRecordByID update login record by id
func (srv *UserSrv) UpdateLoginRecordByID(id uint, value interface{}) (err error) {
	err = pgGetClient().Model(srv.createLoginRecordByID(id)).Updates(value).Error
	return
}

// AddLoginRecord add user login record
func (srv *UserSrv) AddLoginRecord(r *UserLoginRecord, c *elton.Context) (err error) {
	err = pgCreate(r)
	if r.ID != 0 {
		id := r.ID
		ip := r.IP
		go func() {
			lo, err := GetLocationByIP(ip, c)
			if err != nil {
				logger.Error("get location by ip fail",
					zap.String("ip", ip),
					zap.Error(err),
				)
				return
			}
			_ = srv.UpdateLoginRecordByID(id, map[string]string{
				"country":  lo.Country,
				"province": lo.Province,
				"city":     lo.City,
				"isp":      lo.ISP,
			})
		}()
	}
	return
}

// AddTrackRecord add track record
func (srv *UserSrv) AddTrackRecord(r *UserTrackRecord, c *elton.Context) (err error) {
	err = pgCreate(r)
	if r.ID != 0 {
		id := r.ID
		ip := r.IP
		go func() {
			lo, err := GetLocationByIP(ip, c)
			if err != nil {
				logger.Error("get location by ip fail",
					zap.String("ip", ip),
					zap.Error(err),
				)
				return
			}
			_ = srv.UpdateLoginRecordByID(id, map[string]string{
				"country":  lo.Country,
				"province": lo.Province,
				"city":     lo.City,
				"isp":      lo.ISP,
			})
		}()
	}
	return
}

// List list users
func (srv *UserSrv) List(params helper.PGQueryParams, args ...interface{}) (result []*User, err error) {
	result = make([]*User, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// Count count the users
func (srv *UserSrv) Count(args ...interface{}) (count int, err error) {
	return pgCount(&User{}, args...)
}

// ListLoginRecord list login record
func (srv *UserSrv) ListLoginRecord(params helper.PGQueryParams, args ...interface{}) (result []*UserLoginRecord, err error) {
	result = make([]*UserLoginRecord, 0)
	err = pgQuery(params, args...).Find(&result).Error
	return
}

// CountLoginRecord count login record
func (srv *UserSrv) CountLoginRecord(args ...interface{}) (count int, err error) {
	return pgCount(&UserLoginRecord{}, args...)
}

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

// SetRoles set user roles
func (u *UserSession) SetRoles(roles []string) error {
	return u.se.Set(UserRoles, roles)
}

// GetRoles get user roles
func (u *UserSession) GetRoles() []string {
	result, ok := u.se.Get(UserRoles).([]interface{})
	if !ok {
		return nil
	}
	roles := []string{}
	for _, item := range result {
		role, _ := item.(string)
		if role != "" {
			roles = append(roles, role)
		}
	}
	return roles
}

// SetGroups set user groups
func (u *UserSession) SetGroups(groups []string) error {
	return u.se.Set(UserGroups, groups)
}

// GetGroups get user groups
func (u *UserSession) GetGroups() []string {
	result, ok := u.se.Get(UserGroups).([]interface{})
	if !ok {
		return nil
	}
	groups := []string{}
	for _, item := range result {
		group, _ := item.(string)
		if group != "" {
			groups = append(groups, group)
		}
	}
	return groups
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
func NewUserSession(c *elton.Context) *UserSession {
	v, ok := c.Get(session.Key)
	if !ok {
		return nil
	}
	data, ok := c.Get(cs.UserSession)
	if ok {
		us, ok := data.(*UserSession)
		if ok {
			return us
		}
	}
	us := &UserSession{
		se: v.(*session.Session),
	}
	c.Set(cs.UserSession, us)

	return us
}
