package webutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNullPointerTypesNonNil(t *testing.T) {
	s := "abc"
	assert.Equal(t, NullString(s), &s)

	b := true
	assert.Equal(t, NullBool(b), &b)

	var f32 float32 = -1.23

	assert.Equal(t, NullFloat32(f32), &f32)

	var f64 float64 = -1.23

	assert.Equal(t, NullFloat64(f64), &f64)

	var u8 uint8 = 1

	assert.Equal(t, NullUInt8(u8), &u8)

	var u16 uint16 = 1

	assert.Equal(t, NullUInt16(u16), &u16)

	var u32 uint32 = 1

	assert.Equal(t, NullUInt32(u32), &u32)

	var u64 uint64 = 1

	assert.Equal(t, NullUInt64(u64), &u64)

	var u uint = 1

	assert.Equal(t, NullUInt(u), &u)

	var i8 int8 = 1

	assert.Equal(t, NullInt8(i8), &i8)

	var i16 int16 = 1

	assert.Equal(t, NullInt16(i16), &i16)

	var i32 int32 = 1

	assert.Equal(t, NullInt32(i32), &i32)

	var i64 int64 = 1

	assert.Equal(t, NullInt64(i64), &i64)

	var i int = 1

	assert.Equal(t, NullInt(i), &i)

	now := time.Now()
	assert.Equal(t, NullTime(now), &now)
}

func TestNullPointerTypesNil(t *testing.T) {
	s := ""

	var s2 *string

	assert.Equal(t, NullString(s), s2)

	b := false

	var b2 *bool

	assert.Equal(t, NullBool(b), b2)

	var f32 float32 = 0

	var f322 *float32

	assert.Equal(t, NullFloat32(f32), f322)

	var f64 float64 = 0

	var f642 *float64

	assert.Equal(t, NullFloat64(f64), f642)

	var u8 uint8 = 0

	var u82 *uint8

	assert.Equal(t, NullUInt8(u8), u82)

	var u16 uint16 = 0

	var u162 *uint16

	assert.Equal(t, NullUInt16(u16), u162)

	var u32 uint32 = 0

	var u322 *uint32

	assert.Equal(t, NullUInt32(u32), u322)

	var u64 uint64 = 0

	var u642 *uint64

	assert.Equal(t, NullUInt64(u64), u642)

	var u uint = 0

	var u2 *uint

	assert.Equal(t, NullUInt(u), u2)

	var i8 int8 = 0

	var i82 *int8

	assert.Equal(t, NullInt8(i8), i82)

	var i16 int16 = 0

	var i162 *int16

	assert.Equal(t, NullInt16(i16), i162)

	var i32 int32 = 0

	var i322 *int32

	assert.Equal(t, NullInt32(i32), i322)

	var i64 int64 = 0

	var i642 *int64

	assert.Equal(t, NullInt64(i64), i642)

	var i int = 0

	var i2 *int

	assert.Equal(t, NullInt(i), i2)

	now := time.Time{}

	var now2 *time.Time

	assert.Equal(t, NullTime(now), now2)
}
