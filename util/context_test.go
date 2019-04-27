package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/stretchr/testify/assert"
)

func TestGetTrackID(t *testing.T) {
	req := httptest.NewRequest("GET", "/users/me", nil)
	cookieValue := "abcd"
	req.AddCookie(&http.Cookie{
		Name:  config.GetTrackKey(),
		Value: cookieValue,
	})
	c := cod.NewContext(nil, req)
	assert.Equal(t, GetTrackID(c), cookieValue)
}
