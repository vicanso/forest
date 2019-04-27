package service

import (
	"testing"
	"time"

	"github.com/vicanso/cod"

	"github.com/vicanso/forest/util"

	"github.com/stretchr/testify/assert"
	session "github.com/vicanso/cod-session"
)

func TestUserSession(t *testing.T) {
	assert := assert.New(t)
	c := cod.NewContext(nil, nil)
	se := &session.Session{
		Store: session.NewRedisStore(GetRedisClient(), nil),
	}
	se.Fetch()
	c.Set(session.Key, se)
	us := NewUserSession(c)
	account := "tree.xie"
	err := us.SetAccount(account)
	assert.Nil(err)
	assert.Equal(us.GetAccount(), account)

	assert.NotEmpty(us.GetUpdatedAt())

	loginAt := util.NowString()
	us.SetLoginAt(loginAt)
	assert.Equal(us.GetLoginAt(), loginAt)

	token := util.RandomString(8)
	us.SetLoginToken(token)
	assert.Equal(us.GetLoginToken(), token)

	us.se.ID = "abcd"
	us.ClearSessionID()
	assert.Empty(us.se.ID)

	updatedAt := us.GetUpdatedAt()
	time.Sleep(time.Second)
	us.Refresh()
	assert.NotEqual(us.GetUpdatedAt(), updatedAt)
}
