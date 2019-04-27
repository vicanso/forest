package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/forest/util"
)

func TestConnectingCount(t *testing.T) {
	var count uint32 = 100
	SaveConnectingCount(count)
	assert.Equal(t, GetConnectingCount(), count)
}

func TestSyncMap(t *testing.T) {
	assert := assert.New(t)
	key := util.RandomString(8)
	value := 1
	Store(key, value)
	v, ok := Load(key)
	assert.True(ok)
	assert.Equal(v.(int), value)
	_, loaded := LoadOrStore(key, 2)
	assert.True(loaded)
}

func TestLruCache(t *testing.T) {
	assert := assert.New(t)
	key := util.RandomString(8)
	value := 1
	Add(key, value)
	v, found := Get(key)
	assert.True(found)
	assert.Equal(v.(int), value)
	Remove(key)
	_, found = Get(key)
	assert.False(found)
}
