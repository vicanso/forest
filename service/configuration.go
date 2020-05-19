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

package service

import (
	"strings"
	"time"

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/util"
)

const (
	mockTimeKey              = "mockTime"
	sessionSignedKeyCateogry = "signedKey"
	blockIPCategory          = "blockIP"
	routerConfigCategory     = "routerConfig"

	defaultConfigurationLimit = 100
)

var (
	signedKeys = new(elton.AtomicSignedKeys)
)

type (
	// Configuration configuration of application
	Configuration struct {
		helper.Model

		// 配置名称，唯一
		Name string `json:"name" gorm:"type:varchar(30);not null;unique"`
		// 配置分类
		Category string `json:"category" gorm:"type:varchar(20)"`
		// 配置由谁创建
		Owner string `json:"owner" gorm:"type:varchar(20);not null"`
		// 配置状态
		Status int    `json:"status"`
		Data   string `json:"data"`
		// 启用开始时间
		BeginDate *time.Time `json:"beginDate"`
		// 启用结束时间
		EndDate *time.Time `json:"endDate"`
	}
	// ConfigurationQueryParmas configuration query params
	ConfigurationQueryParmas struct {
		Name     string
		Category string
		Limit    int
	}
	// ConfigurationSrv configuration service
	ConfigurationSrv struct {
	}
)

func init() {
	pgGetClient().AutoMigrate(&Configuration{})
	signedKeys.SetKeys(config.GetSignedKeys())
}

// IsValid check the config is valid
func (conf *Configuration) IsValid() bool {
	if conf.Status != cs.ConfigEnabled {
		return false
	}
	now := util.Now().Unix()
	// 如果开始时间大于当前时间，未开始启用
	if conf.BeginDate != nil && conf.BeginDate.Unix() > now {
		return false
	}
	// 如果结束时间少于当前时间，已结束
	if conf.EndDate != nil && conf.EndDate.Unix() < now {
		return false
	}
	return true
}

// createByID create a configuration by id
func (srv *ConfigurationSrv) createByID(id uint) *Configuration {
	c := &Configuration{}
	c.Model.ID = id
	return c
}

// Add add configuration
func (srv *ConfigurationSrv) Add(conf *Configuration) (err error) {
	err = pgCreate(conf)
	return
}

// UpdateByID update configuration by id
func (srv *ConfigurationSrv) UpdateByID(id uint, attrs ...interface{}) (err error) {
	err = pgGetClient().Model(srv.createByID(id)).Update(attrs...).Error
	return
}

// FindByID find configuration by id
func (srv *ConfigurationSrv) FindByID(id uint) (config *Configuration, err error) {
	config = new(Configuration)
	err = pgGetClient().First(config, "id = ?", id).Error
	return
}

// Available get available configs
func (srv *ConfigurationSrv) Available() (configs []*Configuration, err error) {
	result := make([]*Configuration, 0)
	configs = make([]*Configuration, 0)
	err = pgGetClient().Where("status = ?", cs.ConfigEnabled).Find(&result).Error
	if err != nil {
		return
	}
	for _, item := range result {
		if item.IsValid() {
			configs = append(configs, item)
		}
	}
	return
}

// Unavailable get unavailable configs
func (srv *ConfigurationSrv) Unavailable() (configs []*Configuration, err error) {
	result := make([]*Configuration, 0)
	configs = make([]*Configuration, 0)
	err = pgGetClient().Model(&Configuration{}).Find(&result).Error
	if err != nil {
		return
	}
	for _, item := range result {
		if !item.IsValid() {
			configs = append(configs, item)
		}
	}
	return
}

// Refresh refresh configurations
func (srv *ConfigurationSrv) Refresh() (err error) {
	configs, err := srv.Available()
	if err != nil {
		return
	}
	var mockTimeConfig *Configuration

	routerConfigs := make([]*Configuration, 0)
	var signedKeysConfig *Configuration
	blockIPList := make([]string, 0)

	for _, item := range configs {
		if item.Name == mockTimeKey {
			mockTimeConfig = item
			continue
		}

		// 路由配置
		if item.Category == routerConfigCategory {
			routerConfigs = append(routerConfigs, item)
			continue
		}

		// signed key配置
		if item.Category == sessionSignedKeyCateogry {
			signedKeysConfig = item
			continue
		}

		// 黑名单IP
		if item.Category == blockIPCategory {
			blockIPList = append(blockIPList, item.Data)
			continue
		}
	}

	// 如果未配置mock time，则设置为空
	if mockTimeConfig == nil {
		util.SetMockTime("")
	} else {
		util.SetMockTime(mockTimeConfig.Data)
	}

	// 如果数据库中未配置，则使用默认配置
	if signedKeysConfig == nil {
		signedKeys.SetKeys(config.GetSignedKeys())
	} else {
		keys := strings.Split(signedKeysConfig.Data, ",")
		signedKeys.SetKeys(keys)
	}

	// 更新router configs
	updateRouterConfigs(routerConfigs)

	ResetIPBlocker(blockIPList)
	return
}

// GetSignedKeys get signed keys
func GetSignedKeys() elton.SignedKeysGenerator {
	return signedKeys
}

// List list configurations
func (srv *ConfigurationSrv) List(params ConfigurationQueryParmas) (result []*Configuration, err error) {
	result = make([]*Configuration, 0)
	db := pgGetClient()

	if params.Limit <= 0 {
		db = db.Limit(defaultConfigurationLimit)
	} else {
		db = db.Limit(params.Limit)
	}

	if params.Name != "" {
		names := strings.Split(params.Name, ",")
		if len(names) > 1 {
			db = db.Where("name in (?)", names)
		} else {
			db = db.Where("name = (?)", names[0])
		}
	}

	if params.Category != "" {
		categories := strings.Split(params.Category, ",")
		if len(categories) > 1 {
			db = db.Where("category in (?)", categories)
		} else {
			db = db.Where("category = ?", categories[0])
		}
	}
	err = db.Find(&result).Error
	return
}

// DeleteByID delete configuration
func (srv *ConfigurationSrv) DeleteByID(id uint) (err error) {
	err = pgGetClient().Unscoped().Delete(srv.createByID(id)).Error
	return
}
