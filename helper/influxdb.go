// Copyright 2020 tree xie
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

package helper

import (
	"time"

	influxdb "github.com/influxdata/influxdb-client-go"
	influxdbAPI "github.com/influxdata/influxdb-client-go/api"
	"github.com/vicanso/forest/config"
	"go.uber.org/zap"
)

var (
	defaultInfluxSrv *InfluxSrv
)

type (
	InfluxSrv struct {
		client influxdb.Client
		writer influxdbAPI.WriteAPI
	}
)

func init() {
	influxdbConfig := config.GetInfluxdbConfig()
	if influxdbConfig.Disabled {
		defaultInfluxSrv = new(InfluxSrv)
		return
	}
	opts := influxdb.DefaultOptions()
	// 设置批量提交的大小
	opts.SetBatchSize(influxdbConfig.BatchSize)
	// 如果定时提交间隔大于1秒，则设定定时提交间隔
	if influxdbConfig.FlushInterval > time.Millisecond {
		v := influxdbConfig.FlushInterval / time.Millisecond
		opts.SetFlushInterval(uint(v))
	}
	logger.Info("new influxdb client",
		zap.String("uri", influxdbConfig.URI),
		zap.String("org", influxdbConfig.Org),
		zap.String("bucket", influxdbConfig.Bucket),
		zap.Uint("batchSize", influxdbConfig.BatchSize),
		zap.String("token", influxdbConfig.Token[:5]+"..."),
		zap.Duration("interval", influxdbConfig.FlushInterval),
	)
	c := influxdb.NewClientWithOptions(influxdbConfig.URI, influxdbConfig.Token, opts)
	writer := c.WriteAPI(influxdbConfig.Org, influxdbConfig.Bucket)
	newInfluxdbErrorLogger(writer)
	defaultInfluxSrv = &InfluxSrv{
		client: c,
		writer: writer,
	}
}

// newInfluxdbErrorLogger 创建读取出错日志处理，需要注意此功能需要启用新的goroutine
func newInfluxdbErrorLogger(writer influxdbAPI.WriteAPI) {
	go func() {
		for err := range writer.Errors() {
			logger.Error("influxdb write fail",
				zap.Error(err),
			)
		}
	}()
}

// GetInfluxSrv 获取默认的influxdb服务
func GetInfluxSrv() *InfluxSrv {
	return defaultInfluxSrv
}

// Write 写入数据
func (srv *InfluxSrv) Write(measurement string, fields map[string]interface{}, tags map[string]string, ts ...time.Time) {
	if srv.writer == nil {
		return
	}
	var now time.Time
	if len(ts) != 0 {
		now = ts[0]
	} else {
		now = time.Now()
	}

	srv.writer.WritePoint(influxdb.NewPoint(measurement, tags, fields, now))
}

// Close 关闭当前client
func (srv *InfluxSrv) Close() {
	if srv.client == nil {
		return
	}
	srv.client.Close()
}
