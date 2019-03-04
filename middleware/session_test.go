package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/service"
)

func TestNewSession(t *testing.T) {
	fn := NewSession()
	req := httptest.NewRequest("GET", "/users/me", nil)
	resp := httptest.NewRecorder()
	c := cod.NewContext(resp, req)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	if err != nil {
		t.Fatalf("get session fail, %v", err)
	}

	if IsAnonymous(c) != nil ||
		IsLogin(c) != errShouldLogin {
		t.Fatalf("the session status should be anonymous")
	}
	us := service.NewUserSession(c)
	us.SetAccount("tree.xie")

	if IsAnonymous(c) != errLoginAlready ||
		IsLogin(c) != nil {
		t.Fatalf("the session status should be login")
	}
}
