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

package middleware

import (
	"bytes"
	"net/http"

	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

// New Error handler
func NewError() elton.Handler {
	return func(c *elton.Context) error {
		err := c.Next()
		if err == nil {
			return nil
		}
		he, ok := err.(*hes.Error)
		if !ok {
			// 如果不是以http error的形式返回的error则为非主动抛出错误
			logger.Warn("unexpected error",
				zap.String("route", c.Route),
				zap.String("uri", c.Request.RequestURI),
				zap.Error(err),
			)
			he = hes.NewWithError(err)
			he.StatusCode = http.StatusInternalServerError
			he.Exception = true
		}
		if he.StatusCode == 0 {
			he.StatusCode = http.StatusInternalServerError
		}
		if he.Extra == nil {
			he.Extra = make(map[string]interface{})
		}
		he.Extra["route"] = c.Route
		c.StatusCode = he.StatusCode
		c.BodyBuffer = bytes.NewBuffer(he.ToJSON())
		return nil
	}
}
