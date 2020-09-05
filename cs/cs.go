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

package cs

import "strconv"

const (
	// CID context id
	CID = "cid"
	// UserSession user session
	UserSession = "userSession"
)

const (
	// MagicalCaptcha magical captcha(for test only)
	MagicalCaptcha = "0145"
)

const (
	// 状态启用
	StatusEnabled = iota + 1
	// 状态禁用
	StatusDisabled
)

var (
	// 状态列表
	Statuses = []int{
		// 启用
		StatusEnabled,
		// 禁用
		StatusDisabled,
	}
	// 状态列表（字符串）
	StatusesString = []string{
		strconv.Itoa(StatusEnabled),
		strconv.Itoa(StatusDisabled),
	}
)

// 用户角色
const (
	// UserRoleNormal normal user
	UserRoleNormal = "normal"
	// UserRoleSu super user
	UserRoleSu = "su"
	// UserRoleAdmin admin user
	UserRoleAdmin = "admin"
)

var (
	UserRoles = []string{
		UserRoleNormal,
		UserRoleSu,
		UserRoleAdmin,
	}
)
