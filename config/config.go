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

// 应用中的所有配置获取，拉取配置信息使用default配置中的值为默认值，再根据GO_ENV配置的环境变量获取对应的环境配置

package config

import (
	"bytes"
	"net/url"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"github.com/vicanso/forest/validate"
)

var (
	box = packr.New("config", "../configs")
	env = os.Getenv("GO_ENV")

	defaultViper = viper.New()

	// 应用状态
	applicationStatus = ApplicationStatusStopped

	// 应用名称
	appName   string
	version   string
	buildedAt string
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
		Listen       string `validate:"ascii,required"`
		RequestLimit uint   `validate:"required"`
		Name         string `validate:"ascii"`
	}
	// SessionConfig session相关配置信息
	SessionConfig struct {
		// cookie的保存路径
		CookiePath string `validate:"ascii,required"`
		// cookie的key
		Key string `validate:"ascii,required"`
		// cookie的有效期
		TTL time.Duration `validate:"required"`
		// 用于加密cookie的key
		Keys []string `validate:"required"`
		// 用于跟踪用户的cookie
		TrackKey string `validate:"ascii,required"`
	}
	// RedisConfig redis配置
	RedisConfig struct {
		// 连接地址
		Addr string `validate:"hostname_port,required"`
		// 密码
		Password string
		// db(0,1,2等)
		DB int
		// 慢请求时长
		Slow time.Duration `validate:"required"`
		// 最大的正在处理请求量
		MaxProcessing uint32 `validate:"required"`
	}
	// MailConfig email的配置
	MailConfig struct {
		Host     string `validate:"hostname,required"`
		Port     int    `validate:"number,required"`
		User     string `validate:"email,required"`
		Password string `validate:"min=1,max=100"`
	}
	// Influxdb influxdb配置
	InfluxdbConfig struct {
		Bucket        string        `validate:"min=1,max=50"`
		Org           string        `validate:"min=1,max=100"`
		URI           string        `validate:"url,required"`
		Token         string        `validate:"ascii,required"`
		BatchSize     uint          `validate:"min=1,max=5000"`
		FlushInterval time.Duration `validate:"required"`
		Disabled      bool
	}
)

func init() {
	configType := "yml"
	v := viper.New()
	defaultViper.SetConfigType(configType)
	v.SetConfigType(configType)

	configExt := "." + configType
	// 加载默认配置
	data, err := box.Find("default" + configExt)
	if err != nil {
		panic(err)
	}
	// 读取默认配置中的所有配置
	err = v.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		defaultViper.SetDefault(k, v)
	}

	// 根据当前运行环境配置读取
	envConfigFile := GetENV() + configExt
	data, err = box.Find(envConfigFile)
	if err != nil {
		panic(err)
	}
	// 读取当前运行环境对应的配置
	err = defaultViper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	appName = GetString("app")
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

// GetBool 获取bool配置
func GetBool(key string) bool {
	return defaultViper.GetBool(key)
}

// GetInt 获取int配置
func GetInt(key string) int {
	return defaultViper.GetInt(key)
}

// GetUint 获取uint配置
func GetUint(key string) uint {
	return defaultViper.GetUint(key)
}

// GetUint32 获取uint32配置
func GetUint32(key string) uint32 {
	return defaultViper.GetUint32(key)
}

// GetIntDefault 获取int配置，如果为0则返回默认值
func GetIntDefault(key string, defaultValue int) int {
	v := GetInt(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetUint32Default 获取uint32配置，如果为0则返回默认值
func GetUint32Default(key string, defaultValue uint32) uint32 {
	v := GetUint32(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetString 获取string配置
func GetString(key string) string {
	return defaultViper.GetString(key)
}

// GetStringFromENV 根据配置的值，以此为key从环境变量中获取配置的值，如果环境变量中未配置，则返回当前配置中的值
func GetStringFromENV(key string) string {
	value := GetString(key)
	v := os.Getenv(value)
	if v != "" {
		return v
	}
	return value
}

// GetStringDefault 获取string配置，如果未配置则返回默认值
func GetStringDefault(key, defaultValue string) string {
	v := GetString(key)
	if v != "" {
		return v
	}
	return defaultValue
}

// GetDuration 获取duration配置
func GetDuration(key string) time.Duration {
	return defaultViper.GetDuration(key)
}

// GetDurationDefault 获取duration配置，如果为0则返回默认值
func GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	v := GetDuration(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetStringSlice 获取string slice配置
func GetStringSlice(key string) []string {
	return defaultViper.GetStringSlice(key)
}

// GetStringMap 配置string map配置
func GetStringMap(key string) map[string]interface{} {
	return defaultViper.GetStringMap(key)
}

// SetApplicationStatus 设置应用运行状态
func SetApplicationStatus(status int32) {
	atomic.StoreInt32(&applicationStatus, status)
}

// ApplicationIsRunning 判断应用是否正在运行
func ApplicationIsRunning() bool {
	return atomic.LoadInt32(&applicationStatus) == ApplicationStatusRunning
}

// GetVersion 获取应用版本号
func GetVersion() string {
	return version
}

// SetVersion 设置应用版本号
func SetVersion(v string) {
	version = v
}

// GetBuildedAt 获取应用构建时间
func GetBuildedAt() string {
	return buildedAt
}

// SetBuildedAt 设置应用构建时间
func SetBuildedAt(v string) {
	buildedAt = v
}

func GetBasicConfig() BasicConfig {
	prefix := "basic."
	basicConfig := BasicConfig{
		Name:         GetString(prefix + "name"),
		RequestLimit: GetUint(prefix + "requestLimit"),
		Listen:       GetStringFromENV(prefix + "listen"),
	}
	validatePanic(&basicConfig)
	return basicConfig
}

// GetSessionConfig 获取session的配置
func GetSessionConfig() SessionConfig {
	prefix := "session."
	sessConfig := SessionConfig{
		TTL:        GetDuration(prefix + "ttl"),
		Key:        GetString(prefix + "key"),
		CookiePath: GetString(prefix + "path"),
		Keys:       GetStringSlice(prefix + "keys"),
		TrackKey:   GetString(prefix + "trackKey"),
	}
	validatePanic(&sessConfig)
	return sessConfig
}

// GetRedisConfig 获取redis的配置
func GetRedisConfig() RedisConfig {
	prefix := "redis."
	uri := GetStringFromENV(prefix + "uri")
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

// GetMailConfig 获取邮件配置
func GetMailConfig() MailConfig {
	prefix := "mail."
	mailConfig := MailConfig{
		Host:     GetString(prefix + "host"),
		Port:     GetInt(prefix + "port"),
		User:     GetStringFromENV(prefix + "user"),
		Password: GetStringFromENV(prefix + "password"),
	}
	validatePanic(&mailConfig)
	return mailConfig
}

// GetInfluxdbConfig 获取influxdb配置
func GetInfluxdbConfig() InfluxdbConfig {
	prefix := "influxdb."
	influxdbConfig := InfluxdbConfig{
		URI:           GetStringFromENV(prefix + "uri"),
		Bucket:        GetString(prefix + "bucket"),
		Org:           GetString(prefix + "org"),
		Token:         GetStringFromENV(prefix + "token"),
		BatchSize:     GetUint(prefix + "batchSize"),
		FlushInterval: GetDuration(prefix + "flushInterval"),
		Disabled:      GetBool(prefix + "disabled"),
	}
	validatePanic(&influxdbConfig)
	return influxdbConfig
}
