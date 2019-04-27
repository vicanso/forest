package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStack(t *testing.T) {
	assert.True(t, strings.Contains(GetStack(0, 3)[0], "util.GetStack:"))
}
