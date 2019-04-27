package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestNoQuery(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest("GET", "/users/me?a=1", nil)
	c := cod.NewContext(nil, req)
	c.Next = func() error {
		return nil
	}
	assert.Equal(NoQuery(c), errQueryNotAllow)

	req = httptest.NewRequest("GET", "/users/me", nil)
	c.Request = req
	assert.Nil(NoQuery(c))
}

func TestWaitFor(t *testing.T) {
	assert := assert.New(t)
	start := time.Now()
	d := 10 * time.Millisecond
	fn := WaitFor(d)
	c := cod.NewContext(nil, nil)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(err)
	use := time.Since(start)
	assert.True(use >= d)
}
