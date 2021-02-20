// Copyright 2021 tree xie
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
	"github.com/vicanso/elton"
	"github.com/vicanso/forest/tracer"
	"github.com/vicanso/forest/util"
)

// NewTracer create a tracer middleware
func NewTracer() elton.Handler {
	return func(c *elton.Context) error {
		deviceID := c.GetRequestHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = util.GetTrackID(c)
		}
		// 设置tracer的信息
		tracer.SetTracerInfo(tracer.TracerInfo{
			TraceID:  c.ID,
			DeviceID: deviceID,
		})
		return c.Next()
	}
}
