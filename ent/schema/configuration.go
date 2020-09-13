package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

const (
	// 状态启用
	ConfigurationStatusEnabled = iota + 1
	// 状态禁用
	ConfigurationStatusDisabled
)

// Configuration holds the schema definition for the Configuration entity.
type Configuration struct {
	ent.Schema
}

// Mixin 配置信息的mixin
func (Configuration) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields 配置信息的相关字段
func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			Immutable().
			Comment("配置名称"),
		field.String("category").
			NotEmpty().
			Comment("配置分类"),
		field.String("owner").
			NotEmpty().
			Comment("创建者"),
		field.Int8("status").
			Range(ConfigurationStatusEnabled, ConfigurationStatusDisabled).
			Default(ConfigurationStatusEnabled).
			Comment("配置状态，默认为启用状态"),
		field.String("data").
			NotEmpty().
			Comment("配置信息"),
		field.Time("started_at").
			StructTag(`json:"startedAt,omitempty"`).
			Comment("配置启用时间"),
		field.Time("ended_at").
			StructTag(`json:"endedAt,omitempty"`).
			Comment("配置停用时间"),
	}
}

// Edges of the Configuration.
func (Configuration) Edges() []ent.Edge {
	return nil
}
