// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"

	"go.uber.org/zap"

	tracker "github.com/vicanso/cod-tracker"
)

var (
	logger     = log.Default()
	now        = util.NowString
	getTrackID = util.GetTrackID

	getUserSession  = service.NewUserSession
	loadUserSession cod.Handler
)

func init() {
	loadUserSession = middleware.NewSession()
}

func newTracker(action string) cod.Handler {
	return tracker.New(tracker.Config{
		// TODO 添加当前登录用户
		OnTrack: func(info *tracker.Info, _ *cod.Context) {
			logger.Info("tracker",
				zap.String("action", action),
				zap.String("cid", info.CID),
				zap.Int("result", info.Result),
				zap.Any("query", info.Query),
				zap.Any("params", info.Params),
				zap.Any("form", info.Form),
				zap.Error(info.Err),
			)
		},
	})
}
