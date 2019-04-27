package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeApplicationStatus(t *testing.T) {
	assert := assert.New(t)
	if IsApplicationRunning() {
		defer StartApplication()
	}
	StartApplication()
	assert.True(IsApplicationRunning())

	PauseApplication()
	assert.False(IsApplicationRunning())
}
