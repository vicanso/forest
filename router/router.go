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

package router

import (
	"github.com/vicanso/cod"
)

var (
	// groupList 路由组列表
	groupList = make([]*cod.Group, 0)
)

// NewGroup new router group
func NewGroup(path string, handlerList ...cod.Handler) *cod.Group {
	// 如果配置文件中有配置路由
	g := cod.NewGroup(path, handlerList...)
	groupList = append(groupList, g)
	return g
}

// GetGroups get groups
func GetGroups() []*cod.Group {
	return groupList
}
