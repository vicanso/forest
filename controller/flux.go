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
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/influx"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/validate"
)

type fluxCtrl struct{}

// 参数相关定义
type (
	fluxListTrackerParams struct {
		Begin time.Time `json:"begin" validate:"required"`
		End   time.Time `json:"end" validate:"required"`
		Limit int       `json:"limit"`

		Account string `json:"account" validate:"omitempty,xUserAccount"`
		// 用户行为类型筛选
		Action string `json:"action" validate:"omitempty,xTag"`
		// 结果
		Result string `json:"rslt" validate:"omitempty,xTag"`
	}
	fluxListHTTPErrorParams struct {
		Begin time.Time `json:"begin" validate:"required"`
		End   time.Time `json:"end" validate:"required"`
		Limit int       `json:"limit"`

		Account   string `json:"account" validate:"omitempty,xUserAccount"`
		Category  string `json:"category" validate:"omitempty,xTag"`
		Exception string `json:"exception" validate:"omitempty,xBoolean"`
	}

	fluxListRequestParams struct {
		Begin time.Time `json:"begin" validate:"required"`
		End   time.Time `json:"end" validate:"required"`
		Limit int       `json:"limit"`

		Route     string `json:"route" validate:"omitempty,xTag"`
		Service   string `json:"service" validate:"omitempty,xTag"`
		Exception string `json:"exception" validate:"omitempty,xBoolean"`
		// 结果
		Result string `json:"rslt" validate:"omitempty,xTag"`
	}

	// flux tags/fields查询参数
	fluxListTagOrFieldParams struct {
		Measurement string `json:"measurement" validate:"required,xMeasurement"`
	}
	// fluxListTagValuesParams flux tag values查询参数
	fluxListTagValuesParams struct {
		Measurement string `json:"measurement" validate:"required,xMeasurement"`
		Tag         string `json:"tag" validate:"required,xTag"`
	}
)

// 响应相关定义

func init() {
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

	// 获取request相关调用统计
	g.GET(
		"/v1/requests",
		shouldBeAdmin,
		ctrl.listRequest,
	)

	// 获取tag
	// 不校验登录状态，无敏感信息
	g.GET(
		"/v1/tags/{measurement}",
		ctrl.listTag,
	)
	// 获取tag的取值列表
	// 不校验登录状态，无敏感信息
	g.GET(
		"/v1/tag-values/{measurement}/{tag}",
		ctrl.listTagValue,
	)
	// 获取field
	g.GET(
		"/v1/fields/{measurement}",
		ctrl.ListField,
	)
}

// listTag returns the tags of measurement
func (ctrl fluxCtrl) listTag(c *elton.Context) error {
	params := fluxListTagOrFieldParams{}
	err := validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return err
	}
	tags, err := getInfluxSrv().ListTag(c.Context(), params.Measurement)
	if err != nil {
		return err
	}
	c.CacheMaxAge(time.Minute)
	c.Body = map[string][]string{
		"tags": tags,
	}
	return nil
}

// ListField return the fields of measurement
func (ctrl fluxCtrl) ListField(c *elton.Context) error {
	params := fluxListTagOrFieldParams{}
	err := validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return err
	}
	fields, err := getInfluxSrv().ListField(c.Context(), params.Measurement, "-30d")
	if err != nil {
		return err
	}
	c.CacheMaxAge(time.Minute)
	c.Body = map[string][]string{
		"fields": fields,
	}
	return nil
}

// listValue get the values of tag
func (ctrl fluxCtrl) listTagValue(c *elton.Context) error {
	params := fluxListTagValuesParams{}
	err := validate.Do(&params, c.Params.ToMap())
	if err != nil {
		return err
	}
	values, err := getInfluxSrv().ListTagValue(c.Context(), params.Measurement, params.Tag)
	if err != nil {
		return err
	}
	c.CacheMaxAge(time.Minute)
	c.Body = map[string][]string{
		"values": values,
	}
	return nil
}

// listHTTPError list http error
func (ctrl fluxCtrl) listHTTPError(c *elton.Context) error {
	params := fluxListHTTPErrorParams{}
	err := validate.Query(&params, c.Query())
	if err != nil {
		return err
	}
	fields := map[string]any{
		cs.FieldAccount: params.Account,
	}
	if params.Exception != "" {
		exception := false
		if params.Exception == "true" {
			exception = true
		}
		fields[cs.FieldException] = exception
	}
	result, err := getInfluxSrv().Query(c.Context(), influx.QueryParams{
		Measurement: cs.MeasurementHTTPError,
		Begin:       params.Begin,
		End:         params.End,
		Limit:       params.Limit,
		Tags: map[string]string{
			cs.TagCategory: params.Category,
		},
		Fields: fields,
	})
	if err != nil {
		return err
	}
	c.Body = map[string]any{
		"httpErrors": result,
		"count":      len(result),
	}
	return nil
}

// listTracker list user tracker
func (ctrl fluxCtrl) listTracker(c *elton.Context) error {
	params := fluxListTrackerParams{}
	err := validate.Query(&params, c.Query())
	if err != nil {
		return nil
	}

	result, err := getInfluxSrv().Query(c.Context(), influx.QueryParams{
		Measurement: cs.MeasurementUserTracker,
		Begin:       params.Begin,
		End:         params.End,
		Limit:       params.Limit,
		Tags: map[string]string{
			cs.TagAction: params.Action,
			cs.TagResult: params.Result,
		},
		Fields: map[string]any{
			cs.FieldAccount: params.Account,
		},
	})
	if err != nil {
		return err
	}
	c.Body = map[string]any{
		"trackers": result,
		"count":    len(result),
	}
	return nil
}

// listRequest list request
func (ctrl fluxCtrl) listRequest(c *elton.Context) error {
	params := fluxListRequestParams{}
	err := validate.Query(&params, c.Query())
	if err != nil {
		return err
	}
	fields := map[string]any{}
	if params.Exception != "" {
		exception := false
		if params.Exception == "true" {
			exception = true
		}
		fields[cs.FieldException] = exception
	}
	result, err := getInfluxSrv().Query(c.Context(), influx.QueryParams{
		Measurement: cs.MeasurementHTTPRequest,
		Begin:       params.Begin,
		End:         params.End,
		Limit:       params.Limit,
		Tags: map[string]string{
			cs.TagRoute:   params.Route,
			cs.TagService: params.Service,
			cs.TagResult:  params.Result,
		},
		Fields: fields,
	})
	if err != nil {
		return err
	}
	c.Body = map[string]any{
		"requests": result,
		"count":    len(result),
	}
	return nil
}
