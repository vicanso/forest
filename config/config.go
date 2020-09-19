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

// 应用中的所有配置获取，拉取配置信息使用default配置中的值为默认值，再根据GO_ENV配置的环境变量获取对应的环境配置，
// 需要注意，尽可能按单个key的形式来获取对应的配置，这样的方式可以保证针对单个key优先获取GO_ENV对应配置，
// 再获取默认配置，如果一次获取map的形式，如果当前配置对应的map的所有key不全，不会再获取default的配置

package config

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/vicanso/forest/validate"
	"github.com/vicanso/viperx"
)

var (
	box = packr.New("config", "../configs")
	env = os.Getenv("GO_ENV")

	defaultViperX *viperx.ViperX

	// 应用状态
	applicationStatus = ApplicationStatusStopped

	// 应用版本
	applicationVersion string
	// 应用构建时间
	applicationBuildedAt string
)

const (
	// Dev 开发模式下的环境变量
	Dev = "dev"
	// Test 测试环境下的环境变量
	Test = "test"
	// Production 生产环境下的环境变量
	Production = "production"
)

const (
	// ApplicationStatusStopped 应用停止
	ApplicationStatusStopped int32 = iota
	// ApplicationStatusRunning 应用运行中
	ApplicationStatusRunning
	// ApplicationStatusStopping 应用正在停止
	ApplicationStatusStopping
)

type (
	// BasicConfig 应用基本配置信息
	BasicConfig struct {
		// 监听地址
		Listen string `validate:"required,ascii"`
		// 最大处理请求数
		RequestLimit uint `validate:"required"`
		// 应用名称
		Name string `validate:"required,ascii"`
	}
	// SessionConfig session相关配置信息
	SessionConfig struct {
		// cookie的保存路径
		CookiePath string `validate:"required,ascii"`
		// cookie的key
		Key string `validate:"required,ascii"`
		// cookie的有效期
		TTL time.Duration `validate:"required"`
		// 用于加密cookie的key
		Keys []string `validate:"required"`
		// 用于跟踪用户的cookie
		TrackKey string `validate:"required,ascii"`
	}
	// RedisConfig redis配置
	RedisConfig struct {
		// 连接地址
		Addr string `validate:"required,hostname_port"`
		// 密码
		Password string
		// db(0,1,2等)
		DB int
		// 慢请求时长
		Slow time.Duration `validate:"required"`
		// 最大的正在处理请求量
		MaxProcessing uint32 `validate:"required"`
	}
	// PostgresConfig postgres配置
	PostgresConfig struct {
		URI string `validate:"required,uri"`
	}
	// MailConfig email的配置
	MailConfig struct {
		Host     string `validate:"required,hostname"`
		Port     int    `validate:"required,number"`
		User     string `validate:"required,email"`
		Password string `validate:"required,min=1,max=100"`
	}
	// Influxdb influxdb配置
	InfluxdbConfig struct {
		// 存储的bucket
		Bucket string `validate:"required,min=1,max=50"`
		// 配置的组织名称
		Org string `validate:"required,min=1,max=100"`
		// 连接地址
		URI string `validate:"required,url"`
		// 认证的token
		Token string `validate:"required,ascii"`
		// 批量提交大小
		BatchSize uint `validate:"required,min=1,max=5000"`
		// 间隔提交时长
		FlushInterval time.Duration `validate:"required"`
		// 是否禁用
		Disabled bool
	}
	// AlarmConfig alarm配置
	AlarmConfig struct {
		// 接收人列表
		Receivers []string `validate:"required"`
	}
	// LocationConfig 定位配置
	LocationConfig struct {
		Name    string        `validate:"required"`
		Timeout time.Duration `validate:"required"`
		BaseURL string        `validate:"required,url"`
	}

	// MinioConfig minio的配置信息
	MinioConfig struct {
		Endpoint        string `validate:"required,hostname_port"`
		AccessKeyID     string `validate:"required,min=3"`
		SecretAccessKey string `validate:"required,min=6"`
		SSL             bool
	}
)

func init() {

	configType := "yml"
	defaultViperX = viperx.New(configType)

	readers := make([]io.Reader, 0)
	for _, name := range []string{
		"default",
		GetENV(),
	} {
		data, err := box.Find(name + "." + configType)
		if err != nil {
			panic(err)
		}
		readers = append(readers, bytes.NewReader(data))
	}

	err := defaultViperX.ReadConfig(readers...)
	if err != nil {
		panic(err)
	}
}

func validatePanic(v interface{}) {
	err := validate.Do(v, nil)
	if err != nil {
		panic(err)
	}
}

// GetENV 获取当前运行环境
func GetENV() string {
	if env == "" {
		return Dev
	}
	return env
}

// SetApplicationStatus 设置应用运行状态
func SetApplicationStatus(status int32) {
	atomic.StoreInt32(&applicationStatus, status)
}

// GetApplicationStatus 获取应用运行状态
func GetApplicationStatus() int32 {
	return atomic.LoadInt32(&applicationStatus)
}

// ApplicationIsRunning 判断应用是否正在运行
func ApplicationIsRunning() bool {
	return atomic.LoadInt32(&applicationStatus) == ApplicationStatusRunning
}

// GetApplicationVersion 获取应用版本号
func GetApplicationVersion() string {
	return applicationVersion
}

// SetApplicationVersion 设置应用版本号
func SetApplicationVersion(v string) {
	applicationVersion = v
}

// GetApplicationBuildedAt 获取应用构建时间
func GetApplicationBuildedAt() string {
	return applicationBuildedAt
}

// SetApplicationBuildedAt 设置应用构建时间
func SetApplicationBuildedAt(v string) {
	applicationBuildedAt = v
}

// GetBasicConfig 获取基本配置信息
func GetBasicConfig() BasicConfig {
	prefix := "basic."
	basicConfig := BasicConfig{
		Name:         defaultViperX.GetString(prefix + "name"),
		RequestLimit: defaultViperX.GetUint(prefix + "requestLimit"),
		Listen:       defaultViperX.GetStringFromENV(prefix + "listen"),
	}
	validatePanic(&basicConfig)
	return basicConfig
}

// GetSessionConfig 获取session的配置
func GetSessionConfig() SessionConfig {
	prefix := "session."
	sessConfig := SessionConfig{
		TTL:        defaultViperX.GetDuration(prefix + "ttl"),
		Key:        defaultViperX.GetString(prefix + "key"),
		CookiePath: defaultViperX.GetString(prefix + "path"),
		Keys:       defaultViperX.GetStringSlice(prefix + "keys"),
		TrackKey:   defaultViperX.GetString(prefix + "trackKey"),
	}
	validatePanic(&sessConfig)
	return sessConfig
}

// GetRedisConfig 获取redis的配置
func GetRedisConfig() RedisConfig {
	prefix := "redis."
	uri := defaultViperX.GetStringFromENV(prefix + "uri")
	uriInfo, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	// 获取设置的db
	db := 0
	query := uriInfo.Query()
	dbValue := query.Get("db")
	if dbValue != "" {
		db, err = strconv.Atoi(dbValue)
		if err != nil {
			panic(err)
		}
	}
	// 获取密码
	password, _ := uriInfo.User.Password()

	// 获取slow设置的时间间隔
	slowValue := query.Get("slow")
	slow := 100 * time.Millisecond
	if slowValue != "" {
		slow, err = time.ParseDuration(slowValue)
		if err != nil {
			panic(err)
		}
	}

	// 获取最大处理数的配置
	maxProcessing := 1000
	maxValue := query.Get("maxProcessing")
	if maxValue != "" {
		maxProcessing, err = strconv.Atoi(maxValue)
		if err != nil {
			panic(err)
		}
	}

	redisConfig := RedisConfig{
		Addr:          uriInfo.Host,
		Password:      password,
		DB:            db,
		Slow:          slow,
		MaxProcessing: uint32(maxProcessing),
	}

	validatePanic(&redisConfig)
	return redisConfig
}

// GetPostgresConfig 获取postgres配置
func GetPostgresConfig() PostgresConfig {
	prefix := "postgres."
	postgresConfig := PostgresConfig{
		URI: defaultViperX.GetString(prefix + "uri"),
	}
	validatePanic(&postgresConfig)
	return postgresConfig
}

// GetMailConfig 获取邮件配置
func GetMailConfig() MailConfig {
	prefix := "mail."
	mailConfig := MailConfig{
		Host:     defaultViperX.GetString(prefix + "host"),
		Port:     defaultViperX.GetInt(prefix + "port"),
		User:     defaultViperX.GetStringFromENV(prefix + "user"),
		Password: defaultViperX.GetStringFromENV(prefix + "password"),
	}
	validatePanic(&mailConfig)
	return mailConfig
}

// GetInfluxdbConfig 获取influxdb配置
func GetInfluxdbConfig() InfluxdbConfig {
	prefix := "influxdb."
	influxdbConfig := InfluxdbConfig{
		URI:           defaultViperX.GetStringFromENV(prefix + "uri"),
		Bucket:        defaultViperX.GetString(prefix + "bucket"),
		Org:           defaultViperX.GetString(prefix + "org"),
		Token:         defaultViperX.GetStringFromENV(prefix + "token"),
		BatchSize:     defaultViperX.GetUint(prefix + "batchSize"),
		FlushInterval: defaultViperX.GetDuration(prefix + "flushInterval"),
		Disabled:      defaultViperX.GetBool(prefix + "disabled"),
	}
	validatePanic(&influxdbConfig)
	return influxdbConfig
}

// GetAlarmConfig 获取告警配置
func GetAlarmConfig() AlarmConfig {
	prefix := "alarm."
	alarmConfig := AlarmConfig{
		Receivers: defaultViperX.GetStringSlice(prefix + "receivers"),
	}
	validatePanic(&alarmConfig)
	return alarmConfig
}

// GetLocationConfig 获取定位的配置
func GetLocationConfig() LocationConfig {
	prefix := "location."
	locationConfig := LocationConfig{
		Name:    defaultViperX.GetString(prefix + "name"),
		BaseURL: defaultViperX.GetString(prefix + "baseURL"),
		Timeout: defaultViperX.GetDuration(prefix + "timeout"),
	}
	validatePanic(&locationConfig)
	return locationConfig
}

// GetMinioConfig 获取minio的配置
func GetMinioConfig() MinioConfig {
	prefix := "minio."
	minioConfig := MinioConfig{
		Endpoint:        defaultViperX.GetString(prefix + "endpoint"),
		AccessKeyID:     defaultViperX.GetStringFromENV(prefix + "accessKeyID"),
		SecretAccessKey: defaultViperX.GetStringFromENV(prefix + "secretAccessKey"),
		SSL:             defaultViperX.GetBool(prefix + "ssl"),
	}
	validatePanic(&minioConfig)
	return minioConfig
}
