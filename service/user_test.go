package service

import (
	"testing"
	"time"

	"github.com/vicanso/cod"

	"github.com/vicanso/forest/util"

	session "github.com/vicanso/cod-session"
)

func TestUserSession(t *testing.T) {
	c := cod.NewContext(nil, nil)
	se := &session.Session{
		Store: session.NewRedisStore(GetRedisClient(), nil),
	}
	se.Fetch()
	c.Set(session.Key, se)
	us := NewUserSession(c)
	account := "tree.xie"
	err := us.SetAccount(account)
	if err != nil ||
		us.GetAccount() != account {
		t.Fatalf("get/set account fail, %v", err)
	}
	if us.GetUpdatedAt() == "" {
		t.Fatalf("get updated at fail")
	}

	loginAt := util.NowString()
	us.SetLoginAt(loginAt)
	if us.GetLoginAt() != loginAt {
		t.Fatalf("get/set login at fail")
	}

	token := util.RandomString(8)
	us.SetLoginToken(token)
	if us.GetLoginToken() != token {
		t.Fatalf("get/set login token fail")
	}
	us.se.ID = "abcd"
	us.ClearSessionID()
	if us.se.ID != "" {
		t.Fatalf("clear session id fail")
	}

	updatedAt := us.GetUpdatedAt()
	time.Sleep(time.Second)
	us.Refresh()
	if us.GetUpdatedAt() == updatedAt {
		t.Fatalf("session refresh fail")
	}
}
