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
		Exception   string    `json:"exception,omitempty" validate:"omitempty,xBoolean"`
		// 用户行为类型筛选
		Action      string `json:"action,omitempty" validate:"omitempty,xTag"`
		Result      string `json:"result,omitempty" validate:"omitempty,xTag"`
		Category    string `json:"category,omitempty" validate:"omitempty,xTag"`
		ErrCategory string `json:"errCategory,omitempty" validate:"omitempty,xTag"`
		Route       string `json:"route,omitempty" validate:"omitempty,xTag"`
		Service     string `json:"service,omitempty" validate:"omitempty,xTag"`
	}
	// fluxListTagValuesParams flux tag values查询参数
	fluxListTagValuesParams struct {
		Measurement string `json:"measurement,omitempty" validate:"required,xMeasurement"`
		Tag         string `json:"tag,omitempty" validate:"required,xTag"`
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
	// 获取用户action
	g.GET(
		"/v1/actions",
		shouldBeAdmin,
		ctrl.listAction,
	)
	// 获取request相关调用统计
	g.GET(
		"/v1/requests",
		shouldBeAdmin,
		ctrl.listRequest,
	)
	// 获取tag的值
	g.GET(
		"/v1/tag-values/{measurement}/{tag}",
		shouldBeAdmin,
		ctrl.listTagValue,
	)
}

// Query get flux query string
func (params *fluxListParams) Query() string {
	start := util.FormatTime(params.Begin.UTC())
	stop := util.FormatTime(params.End.UTC())
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
	addStrQuery := func(key, value string) {
		query += fmt.Sprintf(`|> filter(fn: (r) => r.%s == "%s")`, key, value)
	}
	addQuery := func(key string, value interface{}) {
		query += fmt.Sprintf(`|> filter(fn: (r) => r.%s == %s)`, key, value)
	}
	// 用户行为类型
	if params.Action != "" {
		addStrQuery("action", params.Action)
	}
	// 结果
	if params.Result != "" {
		addStrQuery("result", params.Result)
	}
	if params.Category != "" {
		addStrQuery("category", params.Category)
	}
	// 账号
	if params.Account != "" {
		addStrQuery("account", params.Account)
	}
	// 异常
	if params.Exception != "" {
		value := "true"
		if params.Exception == "0" {
			value = "false"
		}
		addQuery("exception", value)
	}

	// service
	if params.Service != "" {
		addStrQuery("service", params.Service)
	}

	// route
	if params.Route != "" {
		addStrQuery("route", params.Route)
	}

	// 出错类型
	if params.ErrCategory != "" {
		addStrQuery("errCategory", params.ErrCategory)

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

// listValue get the values of tag
func (ctrl fluxCtrl) listTagValue(c *elton.Context) (err error) {
	params := fluxListTagValuesParams{}
	err = validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return
	}
	values, err := getInfluxSrv().ListTagValue(c.Context(), params.Measurement, params.Tag)
	if err != nil {
		return
	}
	c.Body = map[string][]string{
		"values": values,
	}
	return
}

func (ctrl fluxCtrl) list(c *elton.Context, measurement, responseKey string) (err error) {
	params := fluxListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.Measurement = measurement
	result, err := params.Do(c.Context())
	c.Body = map[string]interface{}{
		responseKey: result,
	}
	return
}

// listHTTPError list http error
func (ctrl fluxCtrl) listHTTPError(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementHTTPError, "httpErrors")
}

// listTracker list user tracker
func (ctrl fluxCtrl) listTracker(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementUserTracker, "trackers")
}

// listAction list user action
func (ctrl fluxCtrl) listAction(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementUserAction, "actions")
}

// listRequest list request
func (ctrl fluxCtrl) listRequest(c *elton.Context) (err error) {
	return ctrl.list(c, cs.MeasurementHTTPRequest, "requests")
}
