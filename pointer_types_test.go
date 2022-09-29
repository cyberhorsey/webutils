package webutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPointerTypes(t *testing.T) {
	s := "abc"
	assert.Equal(t, String(s), &s)

	b := true
	assert.Equal(t, Bool(b), &b)

	var f32 float32 = -1.23

	assert.Equal(t, Float32(f32), &f32)

	var f64 float64 = -1.23

	assert.Equal(t, Float64(f64), &f64)

	var u8 uint8 = 1

	assert.Equal(t, UInt8(u8), &u8)

	var u16 uint16 = 1

	assert.Equal(t, UInt16(u16), &u16)

	var u32 uint32 = 1

	assert.Equal(t, UInt32(u32), &u32)

	var u64 uint64 = 1

	assert.Equal(t, UInt64(u64), &u64)

	var u uint = 1

	assert.Equal(t, UInt(u), &u)

	var i8 int8 = 1

	assert.Equal(t, Int8(i8), &i8)

	var i16 int16 = 1

	assert.Equal(t, Int16(i16), &i16)

	var i32 int32 = 1

	assert.Equal(t, Int32(i32), &i32)

	var i64 int64 = 1

	assert.Equal(t, Int64(i64), &i64)

	var i int = 1

	assert.Equal(t, Int(i), &i)

	now := time.Now()
	assert.Equal(t, Time(now), &now)
}
