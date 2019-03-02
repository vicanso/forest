package util

import (
	"github.com/vicanso/forest/config"
)

// IsDevelopment 判断是否开发环境
func IsDevelopment() bool {
	return config.GetENV() == config.Dev
}

// IsTest 判断是否测试环境
func IsTest() bool {
	return config.GetENV() == config.Test
}

// IsProduction 判断是否生产环境
func IsProduction() bool {
	return config.GetENV() == config.Production
}
