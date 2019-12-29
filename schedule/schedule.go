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

package schedule

import (
	"github.com/robfig/cron/v3"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/service"

	"go.uber.org/zap"
)

func init() {
	c := cron.New()
	_, _ = c.AddFunc("@every 5m", redisCheck)
	_, _ = c.AddFunc("@every 1m", configRefresh)
	c.Start()
}

func redisCheck() {
	err := helper.RedisPing()
	if err != nil {
		log.Default().Error("redis check fail",
			zap.Error(err),
		)
		service.AlarmError("redis check fail")
	}
}

func configRefresh() {
	configSrv := new(service.ConfigurationSrv)
	err := configSrv.Refresh()
	if err != nil {
		log.Default().Error("config refresh fail",
			zap.Error(err),
		)
		service.AlarmError("config refresh fail")
	}
}
