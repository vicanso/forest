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

package influx

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/vicanso/forest/cache"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/forest/validate"
	"github.com/vicanso/go-parallel"
)

type influxSrv struct {
	db *helper.InfluxDB
}

type (
	QueryParams struct {
		Measurement string    `validate:"required"`
		Begin       time.Time `validate:"required"`
		End         time.Time `validate:"required"`
		Limit       int       `default:"20" validate:"required"`
		Tags        map[string]string
		Fields      map[string]any
	}
)

type mutexMapSlice struct {
	mutex sync.Mutex
	data  []map[string]any
}

func newMutexMapSlice() *mutexMapSlice {
	return &mutexMapSlice{}
}

func (m *mutexMapSlice) Add(data ...map[string]any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.data == nil {
		m.data = make([]map[string]any, 0)
	}
	m.data = append(m.data, data...)
}

func (m *mutexMapSlice) List() []map[string]any {
	return m.data
}

// 缓存的数据
type (
	// 缓存的values
	fluxCacheValues struct {
		Values []string `json:"values"`
	}
)

var srv = mustNewSrv()
var ignoreFields = "_start _stop _field _measurement"

// 缓存flux的tags tag-values等
var fluxCache = cache.NewMultilevelCache(10, 5*time.Minute, "flux:")

// mustNewSrv 创建新的influx服务
func mustNewSrv() *influxSrv {
	return &influxSrv{
		db: helper.GetInfluxDB(),
	}
}

// New 获取默认的influxdb服务
func New() *influxSrv {
	return srv
}

// ListTagValue list value of tag
func (srv *influxSrv) ListTagValue(ctx context.Context, measurement, tag string) (values []string, err error) {
	if srv.db == nil {
		return
	}
	// 优先取缓存
	key := fmt.Sprintf("tagValues:%s:%s", measurement, tag)
	result := fluxCacheValues{}
	// 忽略获取失败
	_ = fluxCache.Get(ctx, key, &result)
	if len(result.Values) != 0 {
		values = result.Values
		return
	}
	query := fmt.Sprintf(`import "influxdata/influxdb/schema"
schema.measurementTagValues(
	bucket: "%s",
	measurement: "%s",
	tag: "%s"
)`, srv.db.GetBucket(), measurement, tag)
	items, err := srv.db.Query(ctx, query)
	if err != nil {
		return
	}
	for _, item := range items {
		v, ok := item["_value"]
		if !ok {
			continue
		}
		value, ok := v.(string)
		if !ok {
			continue
		}
		values = append(values, value)
	}
	sort.Strings(values)
	if len(values) != 0 {
		result.Values = values
		_ = fluxCache.Set(ctx, key, &result)
	}
	return
}

// ListTag returns the tag list of measurement
func (srv *influxSrv) ListTag(ctx context.Context, measurement string) (tags []string, err error) {
	if srv.db == nil {
		return
	}
	// 优先取缓存
	key := fmt.Sprintf("tags:%s", measurement)
	result := fluxCacheValues{}
	_ = fluxCache.Get(ctx, key, &result)
	if len(result.Values) != 0 {
		tags = result.Values
		return
	}
	query := fmt.Sprintf(`import "influxdata/influxdb/schema"
schema.measurementTagKeys(
	bucket: "%s",
	measurement: "%s"
)`, srv.db.GetBucket(), measurement)
	items, err := srv.db.Query(ctx, query)
	if err != nil {
		return
	}
	for _, item := range items {
		v, ok := item["_value"]
		if !ok {
			continue
		}
		tag, ok := v.(string)
		if !ok {
			continue
		}
		if strings.Contains(ignoreFields, tag) {
			continue
		}
		tags = append(tags, tag)
	}
	if len(tags) != 0 {
		result.Values = tags
		_ = fluxCache.Set(ctx, key, &result)
	}
	return
}

func addTagQuery(query, key string, value any) string {
	str, ok := value.(string)
	template := `|> filter(fn: (r) => r.%s == %v)
`
	if ok {
		value = strings.ReplaceAll(str, `"`, `\"`)
		template = `|> filter(fn: (r) => r.%s == "%s")
`
	}
	query += fmt.Sprintf(template, key, value)
	return query
}

func addFieldsQuery(query string, fields map[string]any) string {
	if len(fields) == 0 {
		return query
	}
	arr := make([]string, 0)
	for key, value := range fields {
		if value == nil {
			continue
		}
		str, ok := value.(string)
		if ok && str == "" {
			continue
		}
		template := `r._value == %v`
		if ok {
			value = strings.ReplaceAll(str, `"`, `\"`)
			template = `r._value == "%s"`
		}
		template = fmt.Sprintf(template, value)
		arr = append(arr, fmt.Sprintf(`(r._field == "%s" and %s)`, key, template))
	}
	if len(arr) == 0 {
		return query
	}

	query += fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(arr, " or "))
	return query
}

func (srv *influxSrv) pivotQuery(ctx context.Context, measurement string, params map[string]any) ([]map[string]any, error) {
	t, ok := params["_time"].(time.Time)
	if !ok {
		return nil, errors.New("time can not be nil")
	}

	start := t.Format(time.RFC3339Nano)
	stop := t.Add(time.Nanosecond).Format(time.RFC3339Nano)
	ingores := "_measurement _start _stop _time result table"

	filter := ""
	for k, v := range params {
		if v == nil || strings.Contains(ingores, k) {
			continue
		}
		// 由于此处是pivot之后处理，因此都使用tag query
		filter = addTagQuery(filter, k, v)
	}
	query := fmt.Sprintf(`|> range(start: %s, stop: %s)
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	%s`, start, stop, measurement, filter)
	items, err := srv.QueryRaw(ctx, query)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (srv *influxSrv) Query(ctx context.Context, params QueryParams) ([]map[string]any, error) {
	err := validate.Struct(&params)
	if err != nil {
		return nil, err
	}

	start := util.FormatTime(params.Begin.UTC())
	stop := util.FormatTime(params.End.UTC())
	query := fmt.Sprintf(`|> range(start: %s, stop: %s)
|> filter(fn: (r) => r["_measurement"] == "%s")
`,
		start,
		stop,
		params.Measurement,
	)
	for k, v := range params.Tags {
		if v == "" {
			continue
		}
		query = addTagQuery(query, k, v)
	}
	fields := make(map[string]any)
	// 过滤空值
	for k, v := range params.Fields {
		if v == nil {
			continue
		}
		str, ok := v.(string)
		if ok && str == "" {
			continue
		}
		fields[k] = v
	}
	query = addFieldsQuery(query, fields)

	// 筛选完成之后执行pivot
	query += fmt.Sprintf(`|> sort(columns:["_time"], desc: true)
	|> limit(n:%d)
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	`, params.Limit)

	items, err := srv.QueryRaw(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	result := newMutexMapSlice()
	err = parallel.Parallel(func(index int) error {
		// 首次筛选的结果均符合tag，因此pivot的时候将fields也增加匹配
		query := lo.Assign[string, any](items[index], fields)
		tmpItems, err := srv.pivotQuery(ctx, params.Measurement, query)
		if err != nil {
			return err
		}
		result.Add(tmpItems...)
		return nil
	}, len(items), 5)
	if err != nil {
		return nil, err
	}
	items = result.List()
	// 清除不需要字段
	for _, item := range items {
		delete(item, "_measurement")
		delete(item, "_start")
		delete(item, "_stop")
		delete(item, "table")
	}
	sort.Slice(items, func(i, j int) bool {
		t1, _ := items[i]["_time"].(time.Time)
		t2, _ := items[j]["_time"].(time.Time)
		return t1.Before(t2)
	})

	return items, nil
}

func (srv *influxSrv) QueryRaw(ctx context.Context, query string) (items []map[string]any, err error) {
	if srv.db == nil {
		return
	}
	query = fmt.Sprintf(`from(bucket: "%s")
`, srv.db.GetBucket()) + query
	return srv.db.Query(ctx, query)
}

// ListField return the fields of measurement
func (srv *influxSrv) ListField(ctx context.Context, measurement, duration string) (fields []string, err error) {
	if srv.db == nil {
		return
	}
	// 优先取缓存
	key := fmt.Sprintf("fields:%s:%s", measurement, duration)
	result := fluxCacheValues{}
	_ = fluxCache.Get(ctx, key, &result)
	if len(result.Values) != 0 {
		fields = result.Values
		return
	}

	query := fmt.Sprintf(`import "influxdata/influxdb/schema"
schema.measurementFieldKeys(
	bucket: "%s",
	measurement: "%s",
	start: %s
)`, srv.db.GetBucket(), measurement, duration)
	items, err := srv.db.Query(ctx, query)
	if err != nil {
		return
	}
	for _, item := range items {
		v, ok := item["_value"]
		if !ok {
			continue
		}
		field, ok := v.(string)
		if !ok {
			continue
		}
		if strings.Contains(ignoreFields, field) {
			continue
		}
		fields = append(fields, field)
	}
	if len(fields) != 0 {
		result.Values = fields
		_ = fluxCache.Set(ctx, key, &result)
	}
	return
}

// Write 写入数据
func (srv *influxSrv) Write(measurement string, tags map[string]string, fields map[string]any, ts ...time.Time) {
	if srv.db == nil {
		return
	}
	srv.db.Write(measurement, tags, fields, ts...)
}
