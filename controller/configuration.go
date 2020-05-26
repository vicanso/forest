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

	"github.com/vicanso/elton"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/router"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/validate"
)

type (
	configurationCtrl      struct{}
	addConfigurationParams struct {
		Name      string     `json:"name" validate:"xConfigName"`
		Category  string     `json:"category" validate:"xConfigCategory"`
		Status    int        `json:"status" validate:"xConfigStatus"`
		Data      string     `json:"data" validate:"xConfigData"`
		BeginDate *time.Time `json:"beginDate"`
		EndDate   *time.Time `json:"endDate"`
	}
	updateConfigurationParams struct {
		Status    int    `json:"status" validate:"omitempty,xConfigStatus"`
		Category  string `json:"category" validate:"omitempty,xConfigCategory"`
		Data      string `json:"data" validate:"omitempty,xConfigData"`
		BeginDate *time.Time
		EndDate   *time.Time
	}
	listConfigurationParmas struct {
		Name     string `json:"name" validate:"omitempty,xConfigName"`
		Category string `json:"category" validate:"omitempty,xConfigCategory"`
	}
)

func init() {
	g := router.NewGroup("/configurations", loadUserSession, shouldBeSu)
	ctrl := configurationCtrl{}

	g.GET(
		"/v1",
		ctrl.list,
	)

	g.POST(
		"/v1",
		newTracker(cs.ActionConfigurationAdd),
		ctrl.add,
	)
	g.GET(
		"/v1/{id}",
		ctrl.findByID,
	)
	g.PATCH(
		"/v1/{id}",
		newTracker(cs.ActionConfigurationUpdate),
		ctrl.updateByID,
	)
	g.DELETE(
		"/v1/{id}",
		newTracker(cs.ActionConfigurationDelete),
		ctrl.delete,
	)
}

// list configuration
func (ctrl configurationCtrl) list(c *elton.Context) (err error) {
	params := listConfigurationParmas{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	result, err := configSrv.List(service.ConfigurationQueryParmas{
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
func (ctrl configurationCtrl) add(c *elton.Context) (err error) {
	params := addConfigurationParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	us := getUserSession(c)
	conf := &service.Configuration{
		Name:      params.Name,
		Category:  params.Category,
		Status:    params.Status,
		Data:      params.Data,
		Owner:     us.GetAccount(),
		BeginDate: params.BeginDate,
		EndDate:   params.EndDate,
	}
	err = configSrv.Add(conf)
	if err != nil {
		return
	}
	c.Created(conf)
	return
}

// updateByID configuration
func (ctrl configurationCtrl) updateByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := updateConfigurationParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	err = configSrv.UpdateByID(uint(id), service.Configuration{
		Status:    params.Status,
		Data:      params.Data,
		Category:  params.Category,
		BeginDate: params.BeginDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		return
	}

	c.NoContent()
	return
}

// delete configuration
func (ctrl configurationCtrl) delete(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	err = configSrv.DeleteByID(uint(id))
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// findByID find configuration by id
func (ctrl configurationCtrl) findByID(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	data, err := configSrv.FindByID(uint(id))
	if err != nil {
		return
	}
	c.Body = data
	return
}
