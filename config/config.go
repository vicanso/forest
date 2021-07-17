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
	"embed"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vicanso/forest/validate"
	"github.com/vicanso/viperx"
)

//go:embed *.yml
var configFS embed.FS

var (
	env = os.Getenv("GO_ENV")

	defaultViperX = mustLoadConfig()
)

const (
	// Dev 开发模式下的环境变量
	Dev = "dev"
	// Test 测试环境下的环境变量
	Test = "test"
	// Production 生产环境下的环境变量
	Production = "production"
)

type (
	// BasicConfig 应用基本配置信息
	BasicConfig struct {
		// 监听地址
		Listen string `validate:"required,ascii" default:":7001"`
		// 最大处理请求数
		RequestLimit uint `validate:"required" default:"1000"`
		// 应用名称
		Name string `validate:"required,ascii"`
		// PID文件
		PidFile string `validate:"required"`
		// 应用前缀
		Prefixes []string `validate:"omitempty,dive,xPath"`
		// 超时（用于设置所有请求)
		Timeout time.Duration `default:"2m"`
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
		Addrs []string `validate:"required,dive,hostname_port"`
		// 用户名
		Username string
		// 密码
		Password string
		// 慢请求时长
		Slow time.Duration `validate:"required"`
		// 最大的正在处理请求量
		MaxProcessing uint32 `validate:"required" default:"100"`
		// 连接池大小
		PoolSize int `default:"100"`
		// key前缀
		Prefix string
		// sentinel模式下使用的master name
		Master string
	}
	// PostgresConfig postgres配置
	PostgresConfig struct {
		// 连接串
		URI string `validate:"required,uri"`
		// 最大连接数
		MaxOpenConns int `default:"100"`
		// 最大空闲连接数
		MaxIdleConns int `default:"10"`
		// 最大空闲时长
		MaxIdleTime time.Duration `default:"5m"`
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
		BatchSize uint `default:"100" validate:"required,min=1,max=5000"`
		// 间隔提交时长
		FlushInterval time.Duration `default:"30s" validate:"required"`
		// 是否启用gzip
		Gzip bool
		// 是否禁用
		Disabled bool
	}

	// LocationConfig 定位配置
	LocationConfig struct {
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
	// PyroscopeConfig pyroscope的配置信息
	PyroscopeConfig struct {
		Addr  string `validate:"omitempty,url"`
		Token string
	}
)

// mustLoadConfig 加载配置，出错是则抛出panic
func mustLoadConfig() *viperx.ViperX {
	configType := "yml"
	defaultViperX := viperx.New(configType)

	readers := make([]io.Reader, 0)
	for _, name := range []string{
		"default",
		GetENV(),
	} {
		data, err := configFS.ReadFile(name + "." + configType)
		if err != nil {
			panic(err)
		}
		readers = append(readers, bytes.NewReader(data))
	}

	err := defaultViperX.ReadConfig(readers...)
	if err != nil {
		panic(err)
	}
	return defaultViperX
}

// mustValidate 对数据校验，如果出错则panic，仅用于初始化时的配置检查
func mustValidate(v interface{}) {
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

// GetBasicConfig 获取基本配置信息
func GetBasicConfig() *BasicConfig {
	prefix := "basic."
	basicConfig := &BasicConfig{
		Name:         defaultViperX.GetString(prefix + "name"),
		RequestLimit: defaultViperX.GetUint(prefix + "requestLimit"),
		Listen:       defaultViperX.GetStringFromENV(prefix + "listen"),
		Prefixes:     defaultViperX.GetStringSlice(prefix + "prefixes"),
		Timeout:      defaultViperX.GetDuration(prefix + "timeout"),
	}
	pidFile := fmt.Sprintf("%s.pid", basicConfig.Name)
	pwd, _ := os.Getwd()
	if pwd != "" {
		pidFile = pwd + "/" + pidFile
	}
	basicConfig.PidFile = pidFile
	mustValidate(basicConfig)
	return basicConfig
}

// GetSessionConfig 获取session的配置
func GetSessionConfig() *SessionConfig {
	prefix := "session."
	sessConfig := &SessionConfig{
		TTL:        defaultViperX.GetDuration(prefix + "ttl"),
		Key:        defaultViperX.GetString(prefix + "key"),
		CookiePath: defaultViperX.GetString(prefix + "path"),
		Keys:       defaultViperX.GetStringSlice(prefix + "keys"),
		TrackKey:   defaultViperX.GetString(prefix + "trackKey"),
	}
	mustValidate(sessConfig)
	return sessConfig
}

// GetRedisConfig 获取redis的配置
func GetRedisConfig() *RedisConfig {
	prefix := "redis."
	uri := defaultViperX.GetStringFromENV(prefix + "uri")
	uriInfo, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	// 获取密码
	password, _ := uriInfo.User.Password()
	username := uriInfo.User.Username()

	query := uriInfo.Query()
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

	// 转换失败则为0
	poolSize, _ := strconv.Atoi(query.Get("poolSize"))

	redisConfig := &RedisConfig{
		Addrs:         strings.Split(uriInfo.Host, ","),
		Username:      username,
		Password:      password,
		Slow:          slow,
		MaxProcessing: uint32(maxProcessing),
		PoolSize:      poolSize,
		Master:        query.Get("master"),
	}
	keyPrefix := query.Get("prefix")
	if keyPrefix != "" {
		redisConfig.Prefix = keyPrefix + ":"
	}

	mustValidate(redisConfig)
	return redisConfig
}

// GetPostgresConfig 获取postgres配置
func GetPostgresConfig() *PostgresConfig {
	prefix := "postgres."
	uri := defaultViperX.GetStringFromENV(prefix + "uri")
	rawQuery := ""
	uriInfo, _ := url.Parse(uri)
	maxIdleConns := 0
	maxOpenConns := 0
	var maxIdleTime time.Duration
	if uriInfo != nil {
		query := uriInfo.Query()
		rawQuery = uriInfo.RawQuery
		maxIdleConns, _ = strconv.Atoi(query.Get("maxIdleConns"))
		maxOpenConns, _ = strconv.Atoi(query.Get("maxOpenConns"))
		maxIdleTime, _ = time.ParseDuration(query.Get("maxIdleTime"))
	}

	postgresConfig := &PostgresConfig{
		URI:          strings.ReplaceAll(uri, rawQuery, ""),
		MaxIdleConns: maxIdleConns,
		MaxOpenConns: maxOpenConns,
		MaxIdleTime:  maxIdleTime,
	}
	mustValidate(postgresConfig)
	return postgresConfig
}

// GetMailConfig 获取邮件配置
func GetMailConfig() *MailConfig {
	prefix := "mail."
	mailConfig := &MailConfig{
		Host:     defaultViperX.GetString(prefix + "host"),
		Port:     defaultViperX.GetInt(prefix + "port"),
		User:     defaultViperX.GetStringFromENV(prefix + "user"),
		Password: defaultViperX.GetStringFromENV(prefix + "password"),
	}
	mustValidate(mailConfig)
	return mailConfig
}

// GetInfluxdbConfig 获取influxdb配置
func GetInfluxdbConfig() *InfluxdbConfig {
	prefix := "influxdb."
	influxdbConfig := &InfluxdbConfig{
		URI:           defaultViperX.GetStringFromENV(prefix + "uri"),
		Bucket:        defaultViperX.GetString(prefix + "bucket"),
		Org:           defaultViperX.GetString(prefix + "org"),
		Token:         defaultViperX.GetStringFromENV(prefix + "token"),
		BatchSize:     defaultViperX.GetUint(prefix + "batchSize"),
		FlushInterval: defaultViperX.GetDuration(prefix + "flushInterval"),
		Gzip:          defaultViperX.GetBool(prefix + "gzip"),
		Disabled:      defaultViperX.GetBool(prefix + "disabled"),
	}
	mustValidate(influxdbConfig)
	return influxdbConfig
}

// GetLocationConfig 获取定位的配置
func GetLocationConfig() *LocationConfig {
	prefix := "location."
	locationConfig := &LocationConfig{
		BaseURL: defaultViperX.GetString(prefix + "baseURL"),
		Timeout: defaultViperX.GetDuration(prefix + "timeout"),
	}
	mustValidate(locationConfig)
	return locationConfig
}

// GetMinioConfig 获取minio的配置
func GetMinioConfig() *MinioConfig {
	prefix := "minio."
	minioConfig := &MinioConfig{
		Endpoint:        defaultViperX.GetString(prefix + "endpoint"),
		AccessKeyID:     defaultViperX.GetStringFromENV(prefix + "accessKeyID"),
		SecretAccessKey: defaultViperX.GetStringFromENV(prefix + "secretAccessKey"),
		SSL:             defaultViperX.GetBool(prefix + "ssl"),
	}
	mustValidate(minioConfig)
	return minioConfig
}

// GetPyroscopeConfig 获取pyroscope的配置信息
func GetPyroscopeConfig() *PyroscopeConfig {
	prefix := "pyroscope."
	pyroscopeConfig := &PyroscopeConfig{
		Addr:  defaultViperX.GetString(prefix + "addr"),
		Token: defaultViperX.GetString(prefix + "token"),
	}
	return pyroscopeConfig
}
