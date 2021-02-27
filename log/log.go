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

// 可通过zap.RegisterSink添加更多的sink实现不同方式的日志传输

package log

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vicanso/forest/tracer"
	"github.com/vicanso/forest/util"
)

var defaultLogger = mustNewLogger("")

// 如果有配置指定日志级别，则以配置指定的输出
var logLevel = os.Getenv("LOG_LEVEL")

var enabledDebugLog = false

func init() {
	lv, _ := strconv.Atoi(logLevel)
	if lv < 0 {
		enabledDebugLog = true
	}
}

type httpServerLogger struct{}

func (hsl *httpServerLogger) Write(p []byte) (int, error) {
	Default().Info(string(p),
		zap.String("category", "httpServerLogger"),
	)
	return len(p), nil
}

type redisLogger struct{}

func (rl *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	Default().Info(fmt.Sprintf(format, v...),
		zap.String("category", "redisLogger"),
	)
}

type entLogger struct{}

func (el *entLogger) Log(args ...interface{}) {
	Default().Info(fmt.Sprint(args...))
}

type logger struct {
	zapLogger *zap.Logger
}

func (l *logger) getFields(fields []zap.Field) []zap.Field {
	info := tracer.GetTracerInfo()
	// 如果无trace id，则表示获取失败
	if info.TraceID == "" {
		return fields
	}
	return append([]zap.Field{
		zap.String("deviceID", info.DeviceID),
		zap.String("traceID", info.TraceID),
		zap.String("account", info.Account),
	}, fields...)
}
func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, l.getFields(fields)...)
}
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, l.getFields(fields)...)
}
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, l.getFields(fields)...)
}
func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, l.getFields(fields)...)
}
func (l *logger) DPanic(msg string, fields ...zap.Field) {
	l.zapLogger.DPanic(msg, l.getFields(fields)...)
}
func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.zapLogger.Panic(msg, l.getFields(fields)...)
}
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, l.getFields(fields)...)
}
func (l *logger) Sync() error {
	return l.zapLogger.Sync()
}

// DebugEnabled 是否启用了debug日志
func DebugEnabled() bool {
	return enabledDebugLog
}

// mustNewLogger 初始化logger
func mustNewLogger(outputPath string) *logger {

	var c zap.Config
	opts := make([]zap.Option, 0)
	if util.IsDevelopment() {
		c = zap.NewDevelopmentConfig()
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	} else {
		c = zap.NewProductionConfig()
		// 在一秒钟内, 如果某个级别的日志输出量超过了 Initial, 那么在超过之后, 每 Thereafter 条日志才会输出一条, 其余的日志都将被删除
		c.Sampling.Initial = 1000
		// 如果不希望任何日志丢失，则设置为nil
		// c.Sampling = nil

		c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		// 只针对panic 以上的日志增加stack trace
		opts = append(opts, zap.AddStacktrace(zap.DPanicLevel))
	}

	if logLevel != "" {
		lv, _ := strconv.Atoi(logLevel)
		c.Level = zap.NewAtomicLevelAt(zapcore.Level(lv))
	}

	if outputPath != "" {
		c.OutputPaths = []string{
			outputPath,
		}
		c.ErrorOutputPaths = []string{
			outputPath,
		}
	}

	l, err := c.Build(opts...)

	if err != nil {
		panic(err)
	}
	return &logger{
		zapLogger: l,
	}
}

// Default 获取默认的logger
func Default() *logger {
	return defaultLogger
}

// NewHTTPServerLogger create a http server logger
func NewHTTPServerLogger() *log.Logger {
	return log.New(&httpServerLogger{}, "", 0)
}

// NewRedisLogger create a redis logger
func NewRedisLogger() *redisLogger {
	return &redisLogger{}
}

// NewEntLogger create a ent logger
func NewEntLogger() *entLogger {
	return &entLogger{}
}
