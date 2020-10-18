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

package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
	"github.com/vicanso/forest/ent/schema"
	"github.com/vicanso/forest/service"
)

func TestListParams(t *testing.T) {
	assert := assert.New(t)
	params := listParams{
		Limit:  "10",
		Offset: "100",
		Order:  "-id,name",
		Fields: "id,updatedAt",
	}
	assert.Equal(10, params.GetLimit())
	assert.Equal(100, params.GetOffset())
	assert.Equal([]string{
		"id",
		"updated_at",
	}, params.GetFields())
	assert.Equal(2, len(params.GetOrders()))
}

func newContextAndUserSession() (*elton.Context, *service.UserSession) {
	s := session.Session{}
	_, _ = s.Fetch()
	c := elton.NewContext(nil, nil)
	c.Set(session.Key, &s)
	us := service.NewUserSession(c)
	return c, us
}

func TestIsLogin(t *testing.T) {
	assert := assert.New(t)
	c, us := newContextAndUserSession()
	assert.False(isLogin(c))
	err := us.SetInfo(service.UserSessionInfo{
		Account: "trexie",
	})
	assert.Nil(err)
	assert.True(isLogin(c))
}

func TestCheckLogin(t *testing.T) {
	assert := assert.New(t)
	c, us := newContextAndUserSession()
	err := checkLogin(c)
	assert.Equal(errShouldLogin, err)
	err = us.SetInfo(service.UserSessionInfo{
		Account: "trexie",
	})
	assert.Nil(err)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err = checkLogin(c)
	assert.Nil(err)
	assert.True(done)
}

func TestCheckAnonymous(t *testing.T) {
	assert := assert.New(t)
	c, us := newContextAndUserSession()
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err := checkAnonymous(c)
	assert.Nil(err)
	assert.True(done)
	err = us.SetInfo(service.UserSessionInfo{
		Account: "trexie",
	})
	assert.Nil(err)
	err = checkAnonymous(c)
	assert.Equal(errLoginAlready, err)
}

func TestNewCheckRolesMiddleware(t *testing.T) {
	assert := assert.New(t)
	fn := newCheckRolesMiddleware([]string{
		schema.UserRoleAdmin,
	})
	c, us := newContextAndUserSession()
	// 未登录
	err := fn(c)
	assert.Equal(errShouldLogin, err)

	// 已登录但无权限
	err = us.SetInfo(service.UserSessionInfo{
		Account: "trexie",
	})
	assert.Nil(err)
	err = fn(c)
	assert.Equal(errForbidden, err)

	// 已登录且权限允许
	err = us.SetInfo(service.UserSessionInfo{
		Account: "trexie",
		Roles: []string{
			schema.UserRoleAdmin,
		},
	})
	assert.Nil(err)
	done := false
	c.Next = func() error {
		done = true
		return nil
	}
	err = fn(c)
	assert.Nil(err)
	assert.True(done)
}

func TestGetIDFromParams(t *testing.T) {
	assert := assert.New(t)
	c := elton.NewContext(nil, nil)
	c.Params = new(elton.RouteParams)
	c.Params.Add("id", "1")
	id, err := getIDFromParams(c)
	assert.Nil(err)
	assert.Equal(1, id)
}
