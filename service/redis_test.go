package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/forest/util"
)

func TestGetRedisClient(t *testing.T) {
	assert.NotNil(t, GetRedisClient())
}

func TestLock(t *testing.T) {
	assert := assert.New(t)
	key := util.RandomString(8)
	ttl := 10 * time.Millisecond
	success, err := Lock(key, ttl)
	assert.Nil(err)
	assert.True(success, "the first time should lock success")

	success, err = Lock(key, ttl)
	assert.Nil(err)
	assert.False(success, "the second time should lock fail")

	time.Sleep(2 * ttl)
	success, err = Lock(key, ttl)
	assert.Nil(err)
	assert.True(success, "after expired should lock success")

}

func TestLockWithDone(t *testing.T) {
	assert := assert.New(t)
	key := util.RandomString(8)
	ttl := 10 * time.Second
	success, done, err := LockWithDone(key, ttl)
	assert.Nil(err)
	assert.True(success, "the first time should lock success")

	success, _, err = LockWithDone(key, ttl)
	assert.Nil(err)
	assert.False(success, "the second time should lock fail")

	done()
	success, _, err = LockWithDone(key, ttl)
	assert.Nil(err)
	assert.True(success, "after done should lock success")
}

func TestIncWithTTL(t *testing.T) {
	assert := assert.New(t)
	key := util.RandomString(8)
	ttl := 10 * time.Millisecond
	count, err := IncWithTTL(key, ttl)
	assert.Nil(err)
	assert.Equal(count, int64(1))

	count, err = IncWithTTL(key, ttl)
	assert.Nil(err)
	assert.Equal(count, int64(2))

	time.Sleep(ttl)
	count, err = IncWithTTL(key, ttl)
	assert.Nil(err)
	assert.Equal(count, int64(1), "after expired should be reseted")
}
