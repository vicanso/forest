package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	assert := assert.New(t)
	assert.False(IsDevelopment())
	assert.False(IsProduction())
	assert.True(IsTest())
}
