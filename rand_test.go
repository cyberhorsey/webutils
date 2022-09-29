package webutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SecureRandomHex(t *testing.T) {
	for i := 0; i < 999; i++ {
		hex, _ := SecureRandomHex(i)
		assert.Equal(t, i*2, len(hex))
	}
}

func Test_SecureRandomString(t *testing.T) {
	for i := 0; i < 999; i++ {
		str, _ := SecureRandomString(i)
		assert.Equal(t, i, len(str))
	}
}
