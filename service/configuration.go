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

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/util"
)

const (
	mockTimeKey              = "mockTime"
	sessionSignedKeyCateogry = "signedKey"
	routerConfigCategory     = "router-config"
)

var (
	signedKeys = new(cod.AtomicSignedKeys)
)

type (
	// Configuration configuration of application
	Configuration struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt time.Time  `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		// 配置名称，唯一
		Name string `json:"name,omitempty" gorm:"not null;unique"`
		// 配置分类
		Category string `json:"category,omitempty"`
		// 配置由谁创建
		Owner string `json:"owner,omitempty" gorm:"not null;"`
		// 是否启用
		Enabled bool   `json:"enabled,omitempty"`
		Data    string `json:"data,omitempty"`
		// 启用开始时间
		BeginDate time.Time `json:"beginDate,omitempty"`
		// 启用结束时间
		EndDate time.Time `json:"endDate,omitempty"`
	}
	// ConfigurationQueryParmas configuration query params
	ConfigurationQueryParmas struct {
		Name     string
		Category string
	}
)

func init() {
	pgGetClient().AutoMigrate(&Configuration{})
	signedKeys.SetKeys(config.GetSignedKeys())
}

// ConfigurationAdd add configuration
func ConfigurationAdd(conf *Configuration) (err error) {
	err = pgCreate(conf)
	return
}

// ConfigurationUpdate update configuration
func ConfigurationUpdate(conf *Configuration, attrs ...interface{}) (err error) {
	pgGetClient().Model(conf).Update(attrs)
	return
}

// ConfigurationRefresh refresh configurations
func ConfigurationRefresh() (err error) {
	configs := make([]*Configuration, 0)
	err = pgGetClient().Where("enabled = ?", true).Find(&configs).Error
	if err != nil {
		return
	}
	var mockTimeConfig *Configuration
	now := util.Now().Unix()
	routerConfigs := make([]*Configuration, 0)
	var signedKeysConfig *Configuration
	for _, item := range configs {
		// 如果开始时间大于当前时间，未开始启用
		if item.BeginDate.UTC().Unix() > now {
			continue
		}
		// 如果结束时间少于当前时间，已结束
		if item.EndDate.UTC().Unix() < now {
			continue
		}
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
	return
}

// GetSignedKeys get signed keys
func GetSignedKeys() cod.SignedKeysGenerator {
	return signedKeys
}

// ConfigurationList list configurations
func ConfigurationList(params ConfigurationQueryParmas) (result []*Configuration, err error) {
	result = make([]*Configuration, 0)
	db := pgGetClient()
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

// ConfigurationDeleteByID delete configuration
func ConfigurationDeleteByID(id uint) (err error) {
	err = pgGetClient().Unscoped().Delete(&Configuration{
		ID: id,
	}).Error
	return
}
