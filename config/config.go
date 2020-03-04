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

package config

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"github.com/vicanso/forest/validate"
)

var (
	box = packr.New("config", "../configs")
	env = os.Getenv("GO_ENV")
)

type (
	// RedisOptions redis options
	RedisOptions struct {
		Addr     string
		Password string
		DB       int
		// 慢请求时长
		Slow time.Duration
		// 最大的正在处理请求量
		MaxProcessing int
	}
	// SessionConfig session's config
	SessionConfig struct {
		TTL        time.Duration `valid:"-"`
		Key        string        `valid:"ascii"`
		CookiePath string        `valid:"ascii"`
	}
	// MailConfig mail's config
	MailConfig struct {
		Host     string `valid:"host"`
		Port     int    `valid:"port"`
		User     string `valid:"email"`
		Password string `valid:"runelength(1|100)"`
	}

	// Influxdb config
	InfluxdbConfig struct {
		Bucket    string `valid:"runelength(1|50)"`
		Org       string `valid:"runelength(1|100)"`
		URI       string `valid:"url"`
		Token     string `valid:"ascii"`
		BatchSize int    `valid:"range(1|5000)"`
	}
)

const (
	// Dev development env
	Dev = "dev"
	// Test test env
	Test = "test"
	// Production production env
	Production = "production"

	defaultListen   = ":7001"
	defaultTrackKey = "jt"
)

func init() {
	configType := "yml"
	configExt := "." + configType
	data, err := box.Find("default" + configExt)
	if err != nil {
		panic(err)
	}
	viper.SetConfigType(configType)
	v := viper.New()
	v.SetConfigType(configType)
	// 读取默认配置中的所有配置
	err = v.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		viper.SetDefault(k, v)
	}

	// 根据当前运行环境配置读取
	envConfigFile := GetENV() + configExt
	data, err = box.Find(envConfigFile)
	if err != nil {
		panic(err)
	}
	// 读取当前运行环境对应的配置
	err = viper.ReadConfig(bytes.NewReader(data))
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

// GetENV get go env
func GetENV() string {
	if env == "" {
		return Dev
	}
	return env
}

// GetInt viper get int
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetIntDefault get int with default value
func GetIntDefault(key string, defaultValue int) int {
	v := GetInt(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetString viper get string
func GetString(key string) string {
	return viper.GetString(key)
}

// GetStringDefault get string with default value
func GetStringDefault(key, defaultValue string) string {
	v := GetString(key)
	if v != "" {
		return v
	}
	return defaultValue
}

// GetDuration viper get duration
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

// GetDurationDefault get duration with default value
func GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	v := GetDuration(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetStringSlice viper get string slice
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetListen get server listen address
func GetListen() string {
	return GetStringDefault("listen", defaultListen)
}

// GetTrackKey get the track cookie key
func GetTrackKey() string {
	return GetStringDefault("track", defaultTrackKey)
}

// GetRedisConfig get redis config
func GetRedisConfig() (options RedisOptions, err error) {
	value := GetString("redis")
	if value == "" {
		err = errors.New("redis connect uri can't be nil")
		return
	}
	info, err := url.Parse(value)
	if err != nil {
		return
	}
	// 设置默认值
	options.Slow = 300 * time.Millisecond
	options.MaxProcessing = 1000

	options.Addr = info.Host
	pass, ok := info.User.Password()
	// 密码从env中读取
	if ok {
		v := os.Getenv(pass)
		// 如果不为空，则表示从env获取成功
		if v != "" {
			pass = v
		}
	}
	options.Password = pass
	query := info.Query()
	db := query.Get("db")
	if db != "" {
		options.DB, _ = strconv.Atoi(db)
	}
	slow := query.Get("slow")
	if slow != "" {
		d, _ := time.ParseDuration(slow)
		if d != 0 {
			options.Slow = d
		}
	}
	maxProcessing := query.Get("maxProcessing")
	if maxProcessing != "" {
		v, _ := strconv.Atoi(maxProcessing)
		if v != 0 {
			options.MaxProcessing = v
		}
	}
	return
}

// GetPostgresConnectString get postgres connect string
func GetPostgresConnectString() string {
	getPostgresConfig := func(key string) string {
		return viper.GetString("postgres." + key)
	}
	keys := []string{
		"host",
		"port",
		"user",
		"dbname",
		"password",
		"sslmode",
	}
	arr := []string{}
	for _, key := range keys {
		value := getPostgresConfig(key)
		// 密码与用户名支持env中获取
		if key == "password" || key == "user" {
			v := os.Getenv(value)
			if v != "" {
				value = v
			}
		}
		if value != "" {
			arr = append(arr, key+"="+value)
		}
	}
	return strings.Join(arr, " ")
}

// GetSessionConfig get sesion config
func GetSessionConfig() SessionConfig {
	sessConfig := SessionConfig{
		TTL:        viper.GetDuration("session.ttl"),
		Key:        viper.GetString("session.key"),
		CookiePath: viper.GetString("session.path"),
	}
	// 如果session设置过短，则使用默认为24小时
	if sessConfig.TTL < time.Second {
		sessConfig.TTL = 24 * time.Hour
	}
	validatePanic(&sessConfig)
	return sessConfig
}

// GetSignedKeys get signed keys
func GetSignedKeys() []string {
	return viper.GetStringSlice("keys")
}

// GetRouterConcurrentLimit get router concurrent limit
func GetRouterConcurrentLimit() map[string]uint32 {
	limit := make(map[string]uint32)
	data := viper.GetStringMap("routerLimit")
	for key, value := range data {
		v, _ := value.(int)
		if v != 0 {
			arr := strings.Split(key, " ")
			limit[strings.ToUpper(arr[0])+" "+arr[1]] = uint32(v)
		}
	}
	return limit
}

// GetMailConfig get mail config
func GetMailConfig() MailConfig {
	prefix := "mail."
	pass := viper.GetString(prefix + "password")
	if os.Getenv(pass) != "" {
		pass = os.Getenv(pass)
	}
	mailConfig := MailConfig{
		Host:     viper.GetString(prefix + "host"),
		Port:     viper.GetInt(prefix + "port"),
		User:     viper.GetString(prefix + "user"),
		Password: pass,
	}
	validatePanic(&mailConfig)
	return mailConfig
}

// GetInfluxdbConfig get influxdb config
func GetInfluxdbConfig() InfluxdbConfig {
	prefix := "influxdb."
	token := viper.GetString(prefix + "token")
	if os.Getenv(token) != "" {
		token = os.Getenv(token)
	}
	influxdbConfig := InfluxdbConfig{
		URI:       viper.GetString(prefix + "uri"),
		Bucket:    viper.GetString(prefix + "bucket"),
		Org:       viper.GetString(prefix + "org"),
		Token:     token,
		BatchSize: viper.GetInt(prefix + "batchSize"),
	}
	validatePanic(&influxdbConfig)
	return influxdbConfig
}
