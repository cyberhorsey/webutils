package testutils

import (
	"math/rand"
	"time"
	"unsafe"
)

// RandomInt generates a random int between 1 and 1,000,000
func RandomInt() int {
	return RandomIntn(1, 1_000_000)
}

// RandomIntn generates a random int between the specified min and max
// nolint: gosec
func RandomIntn(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomFloat64 generates a random f6 between 1 and 1,000,000
func RandomFloat64() float64 {
	return RandomFloat64n(1, 1_000_000)
}

// RandomFloat64n generates a random f64 between the specified min and max
// nolint: gosec
func RandomFloat64n(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// RandomString generates a random string with the specified length
func RandomString(length int) string {
	return randStringBytesMaskImprSrcUnsafe(length)
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
