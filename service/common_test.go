package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocationByIP(t *testing.T) {
	assert := assert.New(t)
	ip := "1.1.1.1"
	l, err := GetLocationByIP(ip, nil)
	assert.Nil(err)
	assert.Equal(l.IP, ip)
}
