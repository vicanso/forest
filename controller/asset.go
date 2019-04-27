package controller

import (
	"bytes"
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/router"

	staticServe "github.com/vicanso/cod-static-serve"
)

type (
	// assetCtrl asset ctrl
	assetCtrl struct {
	}
	staticFile struct {
		box *packr.Box
	}
)

var (
	box = packr.New("asset", "../web/build")
)

func (sf *staticFile) Exists(file string) bool {
	return sf.box.Has(file)
}
func (sf *staticFile) Get(file string) ([]byte, error) {
	return sf.box.Find(file)
}
func (sf *staticFile) Stat(file string) os.FileInfo {
	return nil
}

func init() {
	g := router.NewGroup("")
	ctrl := assetCtrl{}
	g.GET("/", noQuery, ctrl.index)
	g.GET("/favicon.ico", noQuery, ctrl.favIcon)

	sf := &staticFile{
		box: box,
	}
	g.GET("/static/*file", staticServe.New(sf, staticServe.Config{
		Path: "/static",
		// 客户端缓存一年
		MaxAge: 365 * 24 * 3600,
		// 缓存服务器缓存一个小时
		SMaxAge:             60 * 60,
		DenyQueryString:     true,
		DisableLastModified: true,
	}))
}

func sendFile(c *cod.Context, file string) (err error) {
	buf, err := box.Find(file)
	if err != nil {
		return
	}
	// 根据文件后续设置类型
	c.SetContentTypeByExt(file)
	c.BodyBuffer = bytes.NewBuffer(buf)
	return
}

func (ctrl assetCtrl) index(c *cod.Context) (err error) {
	c.CacheMaxAge("10s")
	return sendFile(c, "index.html")
}

func (ctrl assetCtrl) favIcon(c *cod.Context) (err error) {
	c.SetHeader(cod.HeaderAcceptEncoding, "public, max-age=3600, s-maxage=600")
	return sendFile(c, "favicon.ico")
}
