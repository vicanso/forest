package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/forest/helper"
)

func TestMultiCache(t *testing.T) {
	type TestData struct {
		Name string `json:"name,omitempty"`
	}

	assert := assert.New(t)

	mc := NewMultiCacheWithPrefix(1, time.Minute, "test-cache:")

	data := TestData{
		Name: "nickname",
	}

	key := "test"
	// 首次无数据
	err := mc.Get(key, &TestData{})
	assert.True(helper.RedisIsNilError(err))

	// 设置数据后，查询成功（从lru获取)
	err = mc.Set(key, &data)
	assert.Nil(err)
	result := TestData{}
	err = mc.Get(key, &result)
	assert.Nil(err)
	assert.Equal(data.Name, result.Name)

	// 添加新的数据，lru的数据被更新
	err = mc.Set("a", &TestData{})
	assert.Nil(err)
	result = TestData{}
	// 从redis中获取数据
	err = mc.Get(key, &result)
	assert.Nil(err)
	assert.Equal(data.Name, result.Name)
}
