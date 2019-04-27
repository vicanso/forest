package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
	concurrentLimiter "github.com/vicanso/cod-concurrent-limiter"
	"github.com/vicanso/hes"
)

func TestCreateConcurrentLimitLock(t *testing.T) {
	assert := assert.New(t)
	fn := createConcurrentLimitLock("test-create-concurrent-limit-", time.Second, true)
	c := cod.NewContext(nil, nil)
	key := "abcd"
	success, done, err := fn(key, c)
	assert.Nil(err)
	assert.True(success, "first should lock success")

	success, _, err = fn(key, c)
	assert.Nil(err)
	assert.False(success)
	done()

	success, _, err = fn(key, c)
	assert.Nil(err)
	assert.True(success, "after done should lock success")
}

func TestNewLimiter(t *testing.T) {
	fn := NewLimiter()
	c := cod.NewContext(nil, nil)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(t, err)
}

func TestNewConcurrentLimit(t *testing.T) {
	assert := assert.New(t)
	fn := NewConcurrentLimit([]string{
		"q:type",
	}, time.Second, "test-limit-")
	req := httptest.NewRequest("GET", "/users/me?type=1", nil)
	c := cod.NewContext(nil, req)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(err)
	if err != nil {
		t.Fatalf("concurrent limit fail, %v", err)
	}
	err = fn(c)
	he, ok := err.(*hes.Error)
	assert.True(ok)
	assert.Equal(he.Category, concurrentLimiter.ErrCategory)
}

func TestNewIPLimit(t *testing.T) {
	assert := assert.New(t)
	fn := NewIPLimit(1, time.Second, "test-ip-limit-")
	req := httptest.NewRequest("GET", "/users/me", nil)
	c := cod.NewContext(nil, req)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(err)

	err = fn(c)
	assert.Equal(err, errTooFrequently)
}
