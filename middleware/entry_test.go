package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
)

func TestNewEntry(t *testing.T) {
	assert := assert.New(t)
	fn := NewEntry()
	req := httptest.NewRequest("GET", "/users/me", nil)
	res := httptest.NewRecorder()
	c := cod.NewContext(res, req)
	c.Next = func() error {
		return nil
	}
	err := fn(c)
	assert.Nil(err)
	assert.Equal(c.GetHeader(xResponseID), c.ID)
}
