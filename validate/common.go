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

package validate

func init() {
	// 字符串，max表示字符串长度
	AddAlias("xLimit", "number,min=1,max=2")
	AddAlias("xOffset", "number,min=0,max=5")
	AddAlias("xOrder", "ascii,min=0,max=100")
	AddAlias("xFields", "ascii,min=0,max=100")
	AddAlias("xKeyword", "min=1,max=10")
	// 状态：启用、禁用
	AddAlias("xStatus", "numeric,min=1,max=2")
	// path校验
	AddAlias("xPath", "startswith=/")
}
