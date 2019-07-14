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
	"time"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/util"
)

const (
	mockTimeKey          = "mockTime"
	routerConfigCategory = "router-config"
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

		Name     string `json:"name,omitempty" gorm:"not null;unique"`
		Category string `json:"category,omitempty"`
		Owner    string `json:"owner,omitempty" gorm:"not null;"`
		Enabled  bool   `json:"enabled,omitempty"`
		Data     string `json:"data,omitempty"`
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
	routerConfigs := make([]*Configuration, 0)
	for _, item := range configs {
		if item.Name == mockTimeKey {
			mockTimeConfig = item
			continue
		}
		if item.Category == routerConfigCategory {
			routerConfigs = append(routerConfigs, item)
		}
	}

	// 如果未配置mock time，则设置为空
	if mockTimeConfig == nil {
		util.SetMockTime("")
	} else {
		util.SetMockTime(mockTimeConfig.Data)
	}

	// 更新router configs
	updateRouterConfigs(routerConfigs)
	return
}

// GetSignedKeys get signed keys
func GetSignedKeys() cod.SignedKeysGenerator {
	return signedKeys
}
