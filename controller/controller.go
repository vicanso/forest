package controller

import (
	"regexp"

	"github.com/vicanso/cod"
	tracker "github.com/vicanso/cod-tracker"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"go.uber.org/zap"
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

var maskReg = regexp.MustCompile(`password`)

func createUserTracker(category string) cod.Handler {
	return tracker.New(tracker.Config{
		Mask: maskReg,
		OnTrack: func(info *tracker.Info, c *cod.Context) {
			us := getUserSession(c)
			fields := make([]zap.Field, 0, 5)
			fields = append(fields, zap.String("cid", c.ID))
			if us != nil && us.GetAccount() != "" {
				fields = append(fields, zap.String("account", us.GetAccount()))
			}
			fields = append(fields, zap.String("category", category))
			fields = append(fields, zap.Int("result", info.Result))

			if info.Form != nil {
				fields = append(fields, zap.Any("form", info.Form))
			}

			if info.Params != nil {
				fields = append(fields, zap.Any("params", info.Params))
			}

			if info.Query != nil {
				fields = append(fields, zap.Any("query", info.Query))
			}
			if info.Err != nil {
				fields = append(fields, zap.Any("err", info.Err))
			}
			logger.Info(cs.UserTrackerTag, fields...)
		},
	})
}
