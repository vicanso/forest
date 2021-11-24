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

package schema

import (
	"net/http"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type File struct {
	ent.Schema
}

// Mixin 文件表的mixin
func (File) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields 文件表的字段配置
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("bucket").
			NotEmpty().
			Immutable().
			Comment("文件所在bucket"),
		field.String("filename").
			NotEmpty().
			Immutable().
			Comment("文件名"),
		field.String("contentType").
			NotEmpty().
			Comment("文件数据类型"),
		field.Int64("size").
			NonNegative().
			Comment("文件长度"),
		field.JSON("metadata", &http.Header{}).
			Comment("metadata"),
		field.String("creator").
			NotEmpty().
			Comment("创建者"),
		field.Bytes("data").
			Comment("文件数据"),
	}
}

// Indexes 文件表索引
func (File) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("bucket", "filename").Unique(),
	}
}
