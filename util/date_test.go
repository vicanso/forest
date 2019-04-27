package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	assert := assert.New(t)
	utcTimeStringLen := 20
	now := time.Now()
	assert.True(Now().Unix()-now.Unix() <= 1)

	assert.True(len(NowString()) >= utcTimeStringLen)

	assert.True(UTCNow().Unix()-now.Unix() <= 1)

	assert.Equal(len(UTCNowString()), utcTimeStringLen)

	_, err := ParseTime(NowString())
	assert.Nil(err)

	assert.True(len(FormatTime(now)) >= utcTimeStringLen)
}
