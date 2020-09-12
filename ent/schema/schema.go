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

package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/mixin"
)

// TimeMixin 公共的时间schema
type TimeMixin struct {
	mixin.Schema
}

// Fields 公共时间schema的字段，包括创建于与更新于
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			StructTag(`json:"createdAt,omitempty"`).
			Default(time.Now).
			Comment("创建时间，添加记录时由程序自动生成"),
		field.Time("updated_at").
			StructTag(`json:"updatedAt,omitempty"`).
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间，更新记录时由程序自动生成"),
	}
}
