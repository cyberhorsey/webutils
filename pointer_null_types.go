package webutils

import "time"

func NullString(
	s string,
) *string {
	if s == "" {
		return nil
	}

	return &s
}

func NullBool(
	b bool,
) *bool {
	if !b {
		return nil
	}

	return &b
}

func NullFloat32(
	f float32,
) *float32 {
	if f == 0 {
		return nil
	}

	return &f
}

func NullFloat64(
	f float64,
) *float64 {
	if f == 0 {
		return nil
	}

	return &f
}

func NullUInt8(
	u uint8,
) *uint8 {
	if u == 0 {
		return nil
	}

	return &u
}

func NullUInt16(
	u uint16,
) *uint16 {
	if u == 0 {
		return nil
	}

	return &u
}

func NullUInt32(
	u uint32,
) *uint32 {
	if u == 0 {
		return nil
	}

	return &u
}

func NullUInt64(
	u uint64,
) *uint64 {
	if u == 0 {
		return nil
	}

	return &u
}

func NullUInt(
	u uint,
) *uint {
	if u == 0 {
		return nil
	}

	return &u
}

func NullInt8(
	i int8,
) *int8 {
	if i == 0 {
		return nil
	}

	return &i
}
func NullInt16(
	i int16,
) *int16 {
	if i == 0 {
		return nil
	}

	return &i
}

func NullInt32(
	i int32,
) *int32 {
	if i == 0 {
		return nil
	}

	return &i
}

func NullInt64(
	i int64,
) *int64 {
	if i == 0 {
		return nil
	}

	return &i
}

func NullInt(
	i int,
) *int {
	if i == 0 {
		return nil
	}

	return &i
}

func NullTime(
	t time.Time,
) *time.Time {
	if t.IsZero() {
		return nil
	}

	return &t
}
