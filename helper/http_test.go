package helper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/cod"
	"github.com/vicanso/hes"
	gock "gopkg.in/h2non/gock.v1"
)

func TestHTTPRequest(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	t.Run("normal", func(t *testing.T) {
		gock.New("http://aslant.site").
			Get("/").
			Reply(200).
			JSON(map[string]string{
				"name": "tree.xie",
			})
		req := httptest.NewRequest("GET", "/users/me", nil)
		ipAddr := "1.1.1.1"
		req.Header.Set(xForwardedForHeader, ipAddr)
		c := cod.NewContext(nil, req)
		cid := "abcd"
		c.ID = cid
		d := GetWithContext("http://aslant.site/", c)
		assert.Equal(d.GetValue(contextID).(string), cid)

		d.SetClient(http.DefaultClient)
		resp, body, err := d.Do()
		assert.Nil(err)

		assert.Equal(d.Request.Header.Get(xForwardedForHeader), ipAddr)
		assert.Equal(resp.StatusCode, 200)
		assert.Equal(strings.TrimSpace(string(body)), `{"name":"tree.xie"}`)
	})

	t.Run("post", func(t *testing.T) {
		gock.New("http://aslant.site").
			Post("/user/login").
			Reply(200).
			JSON(map[string]string{
				"name": "tree.xie",
			})
		d := PostWithContext("http://aslant.site/user/login", nil).
			Send(map[string]string{
				"a": "1",
			})
		d.SetClient(http.DefaultClient)
		resp, _, err := d.Do()
		assert.Nil(err)
		assert.Equal(resp.StatusCode, 200)
	})

	t.Run("put", func(t *testing.T) {
		gock.New("http://aslant.site").
			Put("/user").
			Reply(200).
			JSON(map[string]string{
				"name": "tree.xie",
			})
		d := PutWithContext("http://aslant.site/user", nil)
		d.SetClient(http.DefaultClient)
		resp, _, err := d.Do()
		assert.Nil(err)
		assert.Equal(resp.StatusCode, 200)
	})

	t.Run("patch", func(t *testing.T) {
		gock.New("http://aslant.site").
			Patch("/user").
			Reply(200).
			JSON(map[string]string{
				"name": "tree.xie",
			})
		d := PatchWithContext("http://aslant.site/user", nil)
		d.SetClient(http.DefaultClient)
		resp, _, err := d.Do()
		assert.Nil(err)
		assert.Equal(resp.StatusCode, 200)
	})

	t.Run("delete", func(t *testing.T) {
		gock.New("http://aslant.site").
			Delete("/user").
			Reply(200).
			JSON(map[string]string{
				"name": "tree.xie",
			})
		d := DeleteWithContext("http://aslant.site/user", nil)
		d.SetClient(http.DefaultClient)
		resp, _, err := d.Do()
		assert.Nil(err)
		assert.Equal(resp.StatusCode, 200)
	})

	t.Run("error", func(t *testing.T) {
		gock.New("http://aslant.site").
			Get("/").
			Reply(500).
			JSON(map[string]string{
				"message": "get data fail",
			})
		d := GetWithContext("http://aslant.site/", nil)
		d.SetClient(http.DefaultClient)
		_, _, err := d.Do()
		he, ok := err.(*hes.Error)
		assert.True(ok)
		assert.Equal(he.Category, errCategoryHTTPRequest)
		assert.Equal(he.Extra["uri"], "/")
		assert.Equal(he.Extra["method"], "GET")
		assert.Equal(he.Extra["host"], "aslant.site")
		assert.Equal(d.Response.StatusCode, 500)
	})

	t.Run("timeout", func(t *testing.T) {
		d := GetWithContext("https://aslant.site/", nil)
		d.Timeout(time.Millisecond)
		_, _, err := d.Do()
		he, ok := err.(*hes.Error)
		assert.True(ok)
		assert.Equal(he.StatusCode, http.StatusRequestTimeout)
	})
}
