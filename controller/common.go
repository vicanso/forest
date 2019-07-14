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
	"bytes"

	"github.com/vicanso/forest/service"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/router"
)

type (
	commonCtrl struct{}
)

func init() {
	ctrl := commonCtrl{}
	g := router.NewGroup("")

	g.GET("/ping", ctrl.ping)

	g.GET("/ip-location", ctrl.location)

	g.GET("/routers", ctrl.routers)
}

func (ctrl commonCtrl) ping(c *cod.Context) error {
	c.BodyBuffer = bytes.NewBufferString("pong")
	return nil
}

func (ctrl commonCtrl) location(c *cod.Context) (err error) {
	info, err := service.GetLocationByIP(c.RealIP(), c.ID)
	if err != nil {
		return
	}
	c.Body = info
	return
}

func (ctrl commonCtrl) routers(c *cod.Context) (err error) {
	c.Body = c.Cod().Routers
	return
}
