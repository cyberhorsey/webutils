package webutils

import "time"

func String(
	s string,
) *string {
	return &s
}

func Bool(
	b bool,
) *bool {
	return &b
}

func Float32(
	f float32,
) *float32 {
	return &f
}

func Float64(
	f float64,
) *float64 {
	return &f
}

func UInt8(
	u uint8,
) *uint8 {
	return &u
}

func UInt16(
	u uint16,
) *uint16 {
	return &u
}

func UInt32(
	u uint32,
) *uint32 {
	return &u
}

func UInt64(
	u uint64,
) *uint64 {
	return &u
}

func UInt(
	u uint,
) *uint {
	return &u
}

func Int8(
	i int8,
) *int8 {
	return &i
}
func Int16(
	i int16,
) *int16 {
	return &i
}

func Int32(
	i int32,
) *int32 {
	return &i
}

func Int64(
	i int64,
) *int64 {
	return &i
}

func Int(
	i int,
) *int {
	return &i
}

func Time(
	t time.Time,
) *time.Time {
	return &t
}
