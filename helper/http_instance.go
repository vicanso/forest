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

package helper

import (
	"sync/atomic"

	"github.com/vicanso/forest/config"
	"github.com/vicanso/go-axios"
)

var locationIns = newLocationInstance()
var insList = map[string]*axios.Instance{
	"location": locationIns,
}

// newLocationInstance 初始化location的实例
func newLocationInstance() *axios.Instance {
	locationConfig := config.GetLocationConfig()
	return NewHTTPInstance(locationConfig.Name, locationConfig.BaseURL, locationConfig.Timeout)
}

// GetLocationInstance get location instance
func GetLocationInstance() *axios.Instance {
	return locationIns
}

// GetHTTPInstanceStats get http instance stats
func GetHTTPInstanceStats() map[string]interface{} {
	data := make(map[string]interface{})
	for name, ins := range insList {
		data[name] = int(ins.GetConcurrency())
	}
	return data
}

// UpdateInstanceConcurrencyLimit update the concurrency limit for instance
func UpdateInstanceConcurrencyLimit(limits map[string]int) {
	for name, ins := range insList {
		v := limits[name]
		limit := uint32(v)
		if ins.Config.MaxConcurrency != limit {
			atomic.StoreUint32(&ins.Config.MaxConcurrency, limit)
		}
	}
}
