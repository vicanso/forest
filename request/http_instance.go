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

package request

import (
	"sync/atomic"

	"github.com/vicanso/forest/config"
	"github.com/vicanso/go-axios"
)

var locationIns = newLocation()
var locationService = "location"
var insList = map[string]*axios.Instance{
	locationService: locationIns,
}

// newLocation 初始化location的实例
func newLocation() *axios.Instance {
	locationConfig := config.GetLocationConfig()
	return NewHTTP(locationService, locationConfig.BaseURL, locationConfig.Timeout)
}

// GetLocation get location instance
func GetLocation() *axios.Instance {
	return locationIns
}

// GetHTTPStats get http instance stats
func GetHTTPStats() map[string]interface{} {
	data := make(map[string]interface{})
	for name, ins := range insList {
		data[name] = int(ins.GetConcurrency())
	}
	return data
}

// UpdateConcurrencyLimit update the concurrency limit for instance
func UpdateConcurrencyLimit(limits map[string]int) {
	for name, ins := range insList {
		v := limits[name]
		limit := uint32(v)
		if ins.Config.MaxConcurrency != limit {
			atomic.StoreUint32(&ins.Config.MaxConcurrency, limit)
		}
	}
}
