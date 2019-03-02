package middleware

import (
	"github.com/vicanso/cod"
	session "github.com/vicanso/cod-session"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/hes"
)

var (
	errShouldLogin  = hes.New("should login first")
	errLoginAlready = hes.New("login already, please logout first")
)

// NewSession new session middleware
func NewSession() cod.Handler {
	client := service.GetRedisClient()
	if client == nil {
		panic("session store need redis client")
	}
	store := session.NewRedisStore(client, nil)
	store.Prefix = "ss-"
	scf := config.GetSessionConfig()
	return session.NewByCookie(session.CookieConfig{
		Store:   store,
		Signed:  true,
		Expired: scf.TTL,
		GenID: func() string {
			return util.GenUlid()
		},
		Name:     scf.Key,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})
}

func isLogin(c *cod.Context) bool {
	us := service.NewUserSession(c)
	if us == nil || us.GetAccount() == "" {
		return false
	}
	return true
}

// IsLogin check login status, if not login will return error
func IsLogin(c *cod.Context) (err error) {
	if !isLogin(c) {
		err = errShouldLogin
		return
	}
	return c.Next()
}

// IsAnonymous check login status, if login should return error
func IsAnonymous(c *cod.Context) (err error) {
	if isLogin(c) {
		err = errLoginAlready
		return
	}
	return c.Next()
}
