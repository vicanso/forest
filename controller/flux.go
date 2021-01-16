// Copyright 2021 tree xie
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

// flux查询influxdb相关数据

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/forest/validate"
)

type (
	fluxCtrl struct{}

	// fluxListParams flux查询参数
	fluxListParams struct {
		Measurement string    `json:"measurement,omitempty"`
		Begin       time.Time `json:"begin,omitempty" validate:"required"`
		End         time.Time `json:"end,omitempty" validate:"required"`
		Account     string    `json:"account,omitempty" validate:"omitempty,xUserAccount"`
		Limit       string    `json:"limit,omitempty" validate:"required,xLargerLimit"`
	}
)

func init() {
	sessionConfig = config.GetSessionConfig()
	g := router.NewGroup("/fluxes", loadUserSession)

	ctrl := fluxCtrl{}
	// 查询用户tracker
	g.GET(
		"/v1/trackers",
		shouldBeAdmin,
		ctrl.listTracker,
	)
	// 查询http出错
	g.GET(
		"/v1/http-errors",
		shouldBeAdmin,
		ctrl.listHTTPError,
	)
}

func (params *fluxListParams) Query() string {
	start := util.FormatTime(params.Begin)
	stop := util.FormatTime(params.End)
	query := fmt.Sprintf(`
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> sort(columns:["_time"], desc: true)
		|> limit(n:%s)
		|> pivot(
			rowKey:["_time"],
			columnKey: ["_field"],
			valueColumn: "_value"
		)
		`,
		start,
		stop,
		params.Measurement,
		params.Limit,
	)
	if params.Account != "" {
		query += fmt.Sprintf(`|> filter(fn: (r) => r.account == "%s")`, params.Account)
	}
	return query
}

func (params *fluxListParams) Do(ctx context.Context) (items []map[string]interface{}, err error) {
	items, err = getInfluxSrv().Query(ctx, params.Query())
	if err != nil {
		return
	}
	// 清除不需要字段
	for _, item := range items {
		delete(item, "_measurement")
		delete(item, "_start")
		delete(item, "_stop")
		delete(item, "table")
	}
	return
}

// listHTTPError list http error
func (ctrl fluxCtrl) listHTTPError(c *elton.Context) (err error) {
	params := fluxListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.Measurement = cs.MeasurementHTTPError
	result, err := params.Do(c.Context())
	c.Body = map[string]interface{}{
		"httpErrors": result,
	}
	return
}

// listTracker list user tracker
func (ctrl fluxCtrl) listTracker(c *elton.Context) (err error) {
	params := fluxListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.Measurement = cs.MeasurementUserTracker

	result, err := params.Do(c.Context())
	c.Body = map[string]interface{}{
		"trackers": result,
	}
	return
}
