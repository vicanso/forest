package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetListen(t *testing.T) {
	assert.Equal(t, GetListen(), defaultListen)
}

func TestGetENV(t *testing.T) {
	assert.Equal(t, GetENV(), "test")
}

func TestGet(t *testing.T) {
	randomKey := "xx_xx_xx"
	assert := assert.New(t)
	assert.Equal(GetIntDefault("requestLimit", 0), 1024)

	assert.Equal(GetIntDefault(randomKey, 1), 1)

	assert.Equal(GetString("app"), "forest")
	assert.Equal(GetStringDefault(randomKey, "1"), "1")

	assert.Equal(GetDurationDefault(randomKey, time.Second), time.Second)

	assert.Equal(GetStringSlice("keys"), []string{
		"cuttlefish",
		"secret",
	})
}

func TestGetTrackKey(t *testing.T) {
	assert.Equal(t, GetTrackKey(), defaultTrackKey)
}

func TestGetSessionConfig(t *testing.T) {
	assert := assert.New(t)
	scf := GetSessionConfig()
	assert.Equal(scf.TTL, defaultSessionTTL)
	assert.Equal(scf.Key, defaultSessionKey)
	assert.Equal(scf.CookiePath, defaultCookiePath)
}
