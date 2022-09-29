package testutils

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := RandomInt()
		assert.True(t, r >= 1)
		assert.True(t, r <= 1_000_000)
	}
}

// nolint:gosec
func TestRandomIntn(t *testing.T) {
	for i := 0; i < 10; i++ {
		min, max := rand.Int(), rand.Int()
		if min > max {
			min, max = max, min
		}

		r := RandomIntn(min, max)
		assert.True(t, r >= min)
		assert.True(t, r <= max)
	}
}

// nolint: gosec
func TestRandomFloat64(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := RandomFloat64()
		assert.True(t, r >= 1)
		assert.True(t, r <= 1_000_000)
	}
}

// nolint: gosec
func TestRandomFloat64n(t *testing.T) {
	for i := 0; i < 10; i++ {
		min, max := rand.Int(), rand.Int()
		if min > max {
			min, max = max, min
		}

		r := RandomFloat64n(float64(min), float64(max))
		assert.True(t, r >= float64(min))
		assert.True(t, r <= float64(max))
	}
}

// nolint: gosec
func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		l := rand.Intn(100)
		s := RandomString(l)
		assert.Equal(t, len(s), l)
	}
}
