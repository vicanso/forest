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

package cs

const (
	// CID context id
	CID = "cid"
	// UserSession user session
	UserSession = "userSession"

	// UserRoleNormal normal user
	UserRoleNormal = "normal"
	// UserRoleSu super user
	UserRoleSu = "su"
	// UserRoleAdmin admin user
	UserRoleAdmin = "admin"
)

const (
	// ConfigEnabled config enabled
	ConfigEnabled = iota + 1
	// ConfigDiabled config disabled
	ConfigDiabled
)

const (
	// MagicalCaptcha magical captcha(for test only)
	MagicalCaptcha = "0145"
)

const (
	// AccountStatusEnabled account enabled
	AccountStatusEnabled = iota + 1
	// AccountStatusForbidden account forbidden
	AccountStatusForbidden
)
