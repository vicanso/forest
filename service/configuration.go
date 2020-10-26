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

package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/ent"
	"github.com/vicanso/forest/ent/configuration"
	"github.com/vicanso/forest/ent/schema"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/util"
)

type (
	// ConfigurationSrv 配置的相关函数
	ConfigurationSrv struct{}
)

var (
	sessionSignedKeys = new(elton.RWMutexSignedKeys)

	// sessionInterceptorConfig session拦截的配置
	sessionInterceptorConfig = new(sync.Map)
)

const (
	sessionInterceptorKey = "sessionInterceptor"
)

func init() {
	sessionConfig := config.GetSessionConfig()
	// session中用于cookie的signed keys
	sessionSignedKeys.SetKeys(sessionConfig.Keys)
}

// GetSignedKeys 获取用于cookie加密的key列表
func GetSignedKeys() elton.SignedKeysGenerator {
	return sessionSignedKeys
}

// GetSessionInterceptorMessage 获取session拦截的配置信息
func GetSessionInterceptorMessage() (string, bool) {
	value, ok := sessionInterceptorConfig.Load(sessionInterceptorKey)
	if !ok {
		return "", false
	}
	str, ok := value.(string)
	if !ok {
		return "", false
	}
	return str, true
}

// available 获取可用的配置
func (*ConfigurationSrv) available() ([]*ent.Configuration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := util.Now()
	return helper.EntGetClient().Configuration.Query().
		Where(configuration.Status(schema.StatusEnabled)).
		Where(configuration.StartedAtLT(now)).
		Where(configuration.EndedAtGT(now)).
		Order(ent.Desc(configuration.FieldUpdatedAt)).
		All(ctx)
}

// Refresh 刷新配置
func (srv *ConfigurationSrv) Refresh() (err error) {
	configs, err := srv.available()
	if err != nil {
		return
	}
	var mockTimeConfig *ent.Configuration
	routerConcurrencyConfigs := make([]string, 0)
	routerConfigs := make([]string, 0)
	var signedKeys []string
	blockIPList := make([]string, 0)
	sessionInterceptorValue := ""
	for _, item := range configs {
		switch item.Category {
		case schema.ConfigurationCategoryMockTime:
			// 由于排序是按更新时间，因此取最新的记录
			if mockTimeConfig == nil {
				mockTimeConfig = item
			}
		case schema.ConfigurationCategoryBlockIP:
			blockIPList = append(blockIPList, item.Data)
		case schema.ConfigurationCategorySignedKey:
			if len(signedKeys) == 0 {
				signedKeys = strings.Split(item.Data, ",")
			}
		case schema.ConfigurationCategoryRouterConcurrency:
			routerConcurrencyConfigs = append(routerConcurrencyConfigs, item.Data)
		case schema.ConfigurationCategoryRouter:
			routerConfigs = append(routerConfigs, item.Data)
		case schema.ConfigurationCategorySessionInterceptor:
			sessionInterceptorValue = item.Data
		}
	}

	// 设置session interceptor的拦截信息
	if sessionInterceptorValue == "" {
		sessionInterceptorConfig.Delete(sessionInterceptorKey)
	} else {
		sessionInterceptorConfig.Store(sessionInterceptorKey, sessionInterceptorValue)
	}

	// 如果未配置mock time，则设置为空
	if mockTimeConfig == nil {
		util.SetMockTime("")
	} else {
		util.SetMockTime(mockTimeConfig.Data)
	}

	// 如果数据库中未配置，则使用默认配置
	if len(signedKeys) == 0 {
		sessionConfig := config.GetSessionConfig()
		sessionSignedKeys.SetKeys(sessionConfig.Keys)
	} else {
		sessionSignedKeys.SetKeys(signedKeys)
	}

	// 更新router configs
	updateRouterConfigs(routerConfigs)

	// 重置IP拦截列表
	ResetIPBlocker(blockIPList)

	// 重置路由并发限制
	ResetRouterConcurrency(routerConcurrencyConfigs)
	return
}
