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

package helper

import (
	"context"
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/log"
	"go.uber.org/zap"
)

var (
	influxdbClient   *influxdb.Client
	defaultInfluxSrv *InfluxSrv

	initDefaultInfluxSrv sync.Once
)

type (
	InfluxSrv struct {
		sync.Mutex
		BatchSize int
		Bucket    string
		Org       string
		metrics   []influxdb.Metric
	}
)

func init() {
	influxbConfig := config.GetInfluxdbConfig()
	c, err := influxdb.New(influxbConfig.URI, influxbConfig.Token)
	if err != nil {
		panic(err)
	}
	influxdbClient = c
}

// GetInfluxSrv get default influx service
func GetInfluxSrv() *InfluxSrv {
	initDefaultInfluxSrv.Do(func() {
		influxbConfig := config.GetInfluxdbConfig()
		defaultInfluxSrv = &InfluxSrv{
			BatchSize: influxbConfig.BatchSize,
			Bucket:    influxbConfig.Bucket,
			Org:       influxbConfig.Org,
		}
	})
	return defaultInfluxSrv
}

// Write write metric to influxdb
func (srv *InfluxSrv) Write(measurement string, fields map[string]interface{}, tags map[string]string) {
	srv.Lock()
	defer srv.Unlock()
	if len(srv.metrics) == 0 {
		srv.metrics = make([]influxdb.Metric, 0, srv.BatchSize)
	}
	metric := influxdb.NewRowMetric(fields, measurement, tags, time.Now())
	srv.metrics = append(srv.metrics, metric)
	if len(srv.metrics) > srv.BatchSize {
		metrics := srv.metrics
		go func() {
			srv.writeMetrics(metrics)
		}()
		srv.metrics = nil
	}
}

// Flush flush metric list
func (srv *InfluxSrv) Flush() {
	srv.Lock()
	defer srv.Unlock()
	metrics := srv.metrics
	if len(metrics) == 0 {
		return
	}
	go func() {
		srv.writeMetrics(metrics)
	}()
	srv.metrics = nil
}

// writeMetrics write metric list to influxdb
func (srv *InfluxSrv) writeMetrics(metrics []influxdb.Metric) {
	// 如果未初始化client，直接返回
	if influxdbClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := influxdbClient.Write(
		ctx,
		srv.Bucket,
		srv.Org,
		metrics...,
	)
	if err != nil {
		log.Default().Error("influxdb write fail",
			zap.Error(err),
		)
	}
}
