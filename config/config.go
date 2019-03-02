package config

import (
	"bytes"
	"os"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
)

var (
	box = packr.New("config", "../configs")
	env = os.Getenv("GO_ENV")
)

type (
	SessionConfig struct {
		TTL        time.Duration
		Key        string
		SignKeys   []string
		CookiePath string
	}
)

const (
	// Dev development env
	Dev = "dev"
	// Test test env
	Test = "test"
	// Production production env
	Production = "production"
)

func init() {
	configType := "yml"
	data, err := box.Find("default.yml")
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
	envConfigFile := GetENV() + "." + configType
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

// GetListen get server listen address
func GetListen() string {
	return GetString("listen")
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
	v := viper.GetInt(key)
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
	v := viper.GetString(key)
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
	v := viper.GetDuration(key)
	if v.Nanoseconds() != 0 {
		return v
	}
	return defaultValue
}

// GetStringSlice viper get string slice
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetTrackKey get the track cookie key
func GetTrackKey() string {
	return GetString("track")
}

// GetSessionConfig get session config
func GetSessionConfig() *SessionConfig {
	return &SessionConfig{
		TTL:        GetDurationDefault("session.expires", 24*time.Hour),
		Key:        GetStringDefault("session.name", "forest"),
		CookiePath: GetStringDefault("session.cookie.path", "/"),
	}
}
