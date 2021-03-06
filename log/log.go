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
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/vicanso/forest/tracer"
	"github.com/vicanso/forest/util"
)

type TracerHook struct{}

func (h TracerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.NoLevel {
		return
	}
	info := tracer.GetTracerInfo()
	// 如果无trace id，则表示获取失
	if info.TraceID == "" {
		return
	}
	e.Str("deviceID", info.DeviceID).
		Str("traceID", info.TraceID).
		Str("account", info.Account)
}

var defaultLogger = mustNewLogger("")

// 如果有配置指定日志级别，则以配置指定的输出
var logLevel = os.Getenv("LOG_LEVEL")

var enabledDebugLog = false

func init() {
	lv, _ := strconv.Atoi(logLevel)
	if logLevel != "" && lv <= 0 {
		enabledDebugLog = true
	}
}

type httpServerLogger struct{}

func (hsl *httpServerLogger) Write(p []byte) (int, error) {
	Default().Info().
		Str("category", "httpServerLogger").
		Msg(string(p))
	return len(p), nil
}

type redisLogger struct{}

func (rl *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	Default().Info().
		Str("category", "redisLogger").
		Msg(fmt.Sprintf(format, v...))
}

type entLogger struct{}

func (el *entLogger) Log(args ...interface{}) {
	Default().Info().
		Msg(fmt.Sprint(args...))
}

// DebugEnabled 是否启用了debug日志
func DebugEnabled() bool {
	return enabledDebugLog
}

// mustNewLogger 初始化logger
func mustNewLogger(outputPath string) *zerolog.Logger {
	// 如果要节约日志空间，可以配置
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"

	l := zerolog.New(os.Stdout)
	if util.IsDevelopment() {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	l = l.Hook(&TracerHook{}).
		With().Timestamp().Logger()

	if logLevel != "" {
		lv, _ := strconv.Atoi(logLevel)
		l = l.Level(zerolog.Level(lv))
	}

	return &l
}

// Default 获取默认的logger
func Default() *zerolog.Logger {
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

// MapStringString create a map[string]string log event
func MapStringString(data map[string]string) *zerolog.Event {
	if len(data) == 0 {
		return zerolog.Dict()
	}
	m := make(map[string]interface{})
	for k, v := range data {
		m[k] = v
	}
	return zerolog.Dict().Fields(m)
}

// URLValues create a url.Values log event
func URLValues(query url.Values) *zerolog.Event {
	if len(query) == 0 {
		return zerolog.Dict()
	}
	m := make(map[string]interface{})
	for k, values := range query {
		m[k] = values
	}
	return zerolog.Dict().Fields(m)
}

// Struct create a struct log event
func Struct(data interface{}) *zerolog.Event {
	if data == nil {
		return zerolog.Dict()
	}
	buf, _ := json.Marshal(data)
	// 忽略错误，如果不成功则直接返回
	if len(buf) == 0 {
		return zerolog.Dict()
	}
	m := make(map[string]interface{})
	// 出错忽略
	_ = json.Unmarshal(buf, &m)
	return zerolog.Dict().Fields(m)
}
