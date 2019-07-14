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

package controller

import (
	"strconv"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/validate"
)

type (
	configurationCtrl      struct{}
	addConfigurationParams struct {
		Name     string `json:"name,omitempty" valid:"xConfigName"`
		Category string `json:"category,omitempty" valid:"xConfigCategory,optional"`
		Enabled  bool   `json:"enabled,omitempty" valid:"-"`
		Data     string `json:"data,omitempty" valid:"xConfigData"`
	}
	updateConfigurationParams struct {
		Enabled  bool   `json:"enabled,omitempty" valid:"-"`
		Category string `json:"category,omitempty" valid:"xConfigCategory,optional"`
		Data     string `json:"data,omitempty" valid:"xConfigData,optional"`
	}
)

func init() {
	// TODO 增加用户权限判断
	g := router.NewGroup("/configurations")
	ctrl := configurationCtrl{}

	g.POST(
		"/v1",
		newTracker("add-configuration"),
		ctrl.add,
	)
	g.PATCH(
		"/v1/:id",
		newTracker("update-configuration"),
		ctrl.update,
	)
}

// add configuration
func (ctrl configurationCtrl) add(c *cod.Context) (err error) {
	params := &addConfigurationParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	conf := &service.Configuration{
		Name:     params.Name,
		Category: params.Category,
		Enabled:  params.Enabled,
		Data:     params.Data,
		// TODO owner设置为当前登录用户
		Owner: "foo",
	}
	err = service.ConfigurationAdd(conf)
	if err != nil {
		return
	}
	c.Created(conf)
	return
}

func (ctrl configurationCtrl) update(c *cod.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := &updateConfigurationParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	err = service.ConfigurationUpdate(&service.Configuration{
		ID: uint(id),
	}, map[string]interface{}{
		"enabled":  params.Enabled,
		"data":     params.Data,
		"category": params.Category,
	})
	if err != nil {
		return
	}

	c.NoContent()
	return
}
