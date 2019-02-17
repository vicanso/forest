package router

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
)

const (
	apiPrefixKey = "apiPrefix"
)

var (
	// groupList 路由组列表
	groupList = make([]*cod.Group, 0)
)

// NewGroup new router group
func NewGroup(path string, handlerList ...cod.Handler) *cod.Group {
	// 如果配置文件中有配置路由
	path = config.GetString(apiPrefixKey) + path
	g := cod.NewGroup(path, handlerList...)
	groupList = append(groupList, g)
	return g
}

// GetGroups get groups
func GetGroups() []*cod.Group {
	return groupList
}
