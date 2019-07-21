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
	"time"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/validate"
)

type (
	configurationCtrl      struct{}
	addConfigurationParams struct {
		Name      string    `json:"name,omitempty" valid:"xConfigName"`
		Category  string    `json:"category,omitempty" valid:"xConfigCategory,optional"`
		Enabled   bool      `json:"enabled,omitempty" valid:"-"`
		Data      string    `json:"data,omitempty" valid:"xConfigData"`
		BeginDate time.Time `json:"beginDate" valid:"-"`
		EndDate   time.Time `json:"endDate" valid:"-"`
	}
	updateConfigurationParams struct {
		Enabled   bool      `json:"enabled,omitempty" valid:"-"`
		Category  string    `json:"category,omitempty" valid:"xConfigCategory,optional"`
		Data      string    `json:"data,omitempty" valid:"xConfigData,optional"`
		BeginDate time.Time `json:"beginDate" valid:"-"`
		EndDate   time.Time `json:"endDate" valid:"-"`
	}
	listConfigurationParmas struct {
		Name     string `json:"name,omitempty" valid:"xConfigName,optional"`
		Category string `json:"category,omitempty" valid:"xConfigCategory,optional"`
	}
)

func init() {
	// TODO 增加用户权限判断
	g := router.NewGroup("/configurations", loadUserSession)
	ctrl := configurationCtrl{}

	g.GET(
		"/v1",
		shouldBeAdmin,
		ctrl.list,
	)

	g.POST(
		"/v1",
		newTracker(cs.ActionConfigurationAdd),
		shouldBeAdmin,
		ctrl.add,
	)
	g.PATCH(
		"/v1/:id",
		newTracker(cs.ActionConfigurationUpdate),
		shouldBeAdmin,
		ctrl.update,
	)
	g.DELETE(
		"/v1/:id",
		newTracker(cs.ActionConfigurationDelete),
		shouldBeAdmin,
		ctrl.delete,
	)
}

// list configuration
func (ctrl configurationCtrl) list(c *cod.Context) (err error) {
	params := &listConfigurationParmas{}
	err = validate.Do(params, c.Query())
	if err != nil {
		return
	}
	result, err := service.ConfigurationList(service.ConfigurationQueryParmas{
		Name:     params.Name,
		Category: params.Category,
	})
	if err != nil {
		return
	}
	c.Body = map[string]interface{}{
		"configs": result,
	}
	return
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
		Owner:     "foo",
		BeginDate: params.BeginDate,
		EndDate:   params.EndDate,
	}
	err = service.ConfigurationAdd(conf)
	if err != nil {
		return
	}
	c.Created(conf)
	return
}

// update configuration
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
		"enabled":   params.Enabled,
		"data":      params.Data,
		"category":  params.Category,
		"beginDate": params.BeginDate,
		"endDate":   params.EndDate,
	})
	if err != nil {
		return
	}

	c.NoContent()
	return
}

// delete configuration
func (ctrl configurationCtrl) delete(c *cod.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	err = service.ConfigurationDeleteByID(uint(id))
	if err != nil {
		return
	}
	c.NoContent()
	return
}
