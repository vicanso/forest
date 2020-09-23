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

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vicanso/forest/util"
)

var defaultLogger = newLoggerX()

// newLoggerX 初始化logger
func newLoggerX() *zap.Logger {

	if util.IsDevelopment() {
		c := zap.NewDevelopmentConfig()
		l, err := c.Build()
		if err != nil {
			panic(err)
		}
		return l
	}
	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 只针对panic 以上的日志增加stack trace
	l, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}
	return l
}

// Default 获取默认的logger
func Default() *zap.Logger {
	return defaultLogger
}
