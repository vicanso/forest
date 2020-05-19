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

package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/vicanso/forest/cs"
)

func init() {
	// 应用配置名称
	AddAlias("xConfigName", "min=2,max=20")
	AddAlias("xConfigCategory", "alphanum,min=2,max=20")
	AddAlias("xConfigData", "min=0,max=500")

	Add("xConfigStatus", func(fl validator.FieldLevel) bool {
		return isInInt(fl, []int{
			cs.ConfigEnabled,
			cs.ConfigDiabled,
		})
	})
}
