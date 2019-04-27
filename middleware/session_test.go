package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/service"
)

func TestNewSession(t *testing.T) {
	assert := assert.New(t)
	fn := NewSession()
	req := httptest.NewRequest("GET", "/users/me", nil)
	resp := httptest.NewRecorder()
	c := cod.NewContext(resp, req)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(err)

	assert.Nil(IsAnonymous(c))
	assert.Equal(IsLogin(c), errShouldLogin)
	us := service.NewUserSession(c)
	us.SetAccount("tree.xie")

	assert.Equal(IsAnonymous(c), errLoginAlready)
	assert.Nil(IsLogin(c))
}
