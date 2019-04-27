package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	assert.Equal(t, len(RandomString(8)), 8)
}

func TestGenUlid(t *testing.T) {
	assert.Equal(t, len(GenUlid()), 26)
}

func TestSha256(t *testing.T) {
	assert.Equal(t, Sha256("abcd"), "iNQmb9TmM40TuEX88olXnSCciXgjuSF9o+Fhk28DFYk=")
}

func TestContainsString(t *testing.T) {
	arr := []string{
		"a",
		"b",
	}
	assert.True(t, ContainsString(arr, "b"))
	assert.False(t, ContainsString(arr, "c"))
}
