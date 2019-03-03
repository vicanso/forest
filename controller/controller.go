package controller

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
)

var (
	logger     = log.Default()
	now        = util.NowString
	getTrackID = util.GetTrackID

	noQuery               = middleware.NoQuery
	waitFor               = middleware.WaitFor
	createConcurrentLimit = middleware.NewConcurrentLimit
	isLogin               = middleware.IsLogin
	isAnonymous           = middleware.IsAnonymous
	newIPLimit            = middleware.NewIPLimit

	getUserSession = service.NewUserSession

	userSession cod.Handler
)

func init() {
	userSession = middleware.NewSession()
}
