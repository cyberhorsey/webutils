package webutils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

// RAND strings
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// SecureRandomString generates a secure random string of a given length, returning a string of
// that same length built with characters from the charset above.
func SecureRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.Wrap(err, "rand.Read")
	}

	for i := range bytes {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}

	return string(bytes), nil
}

// SecureRandomHex returns a secure random hex of a given length times two.
// IE: if you pass in length of 5, you will get a hex encoded string of length 10 back.
// There is no way to generate odd numbered strings.
func SecureRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.Wrap(err, "rand.Read")
	}

	return hex.EncodeToString(bytes), nil
}
