package config

// 系统相关配置，读取default.yml的相关配置为默认配置，
// 再根据当前配置的GO_ENV读取对应的配置文件

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

var env = os.Getenv("GO_ENV")
var configPath = os.Getenv("CONFIG")
var viperInitTest = os.Getenv("VIPER_INIT_TEST")

// init viper config
func viperInit(path string) error {
	configType := "yml"
	defaultPath := "."
	v := viper.New()
	// 从default中读取默认的配置
	v.SetConfigName("default")
	v.AddConfigPath(defaultPath)
	v.AddConfigPath(path)
	v.SetConfigType(configType)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		viper.SetDefault(k, v)
	}
	// 根据配置的env读取相应的配置信息
	if env != "" {
		viper.SetConfigName(env)
		viper.AddConfigPath(defaultPath)
		viper.AddConfigPath(path)
		viper.SetConfigType(configType)
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

// set default config for test
func setDefaultForTest() {
	viper.Set("locationByIP", "http://ip.taobao.com/service/getIpInfo.php")
	viper.Set("redis", "redis://127.0.0.1:6379")
	viper.Set("db.uri", "postgres://tree:mypwd@127.0.0.1:5432/forest?connect_timeout=5&sslmode=disable")
	viper.Set("influxdb", "http://127.0.0.1:8086/?db=forest")
	viper.Set("app", "forest")
}

func init() {
	// 如果设置为测试，则直接读测试的配置
	if viperInitTest != "" {
		setDefaultForTest()
		return
	}
	if configPath == "" {
		runPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configPath = runPath + "/configs"
	}
	err := viperInit(configPath)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

// GetENV get the go env
func GetENV() string {
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

// GetTrackKey get track key
func GetTrackKey() string {
	v := viper.GetString("track")
	if v == "" {
		return "jt"
	}
	return v
}

// GetSessionKeys get the encrypt keys for session
func GetSessionKeys() []string {
	v := viper.GetStringSlice("session.keys")
	if len(v) == 0 {
		return []string{
			"cuttlefish",
		}
	}
	return v
}

// GetSessionCookie get the session cookie's name
func GetSessionCookie() string {
	v := viper.GetString("session.cookie.name")
	if v == "" {
		return "sess"
	}
	return v
}

// GetSessionCookieMaxAge get the session cookie's max age
func GetSessionCookieMaxAge() int {
	v := viper.GetDuration("session.cookie.maxAge")
	seconds := v / time.Second
	if seconds == 0 {
		return 24 * 3600
	}
	return int(seconds)
}

// GetCookiePath get the cookie's path
func GetCookiePath() string {
	v := viper.GetString("session.cookie.path")
	if v == "" {
		return "/"
	}
	return v
}
