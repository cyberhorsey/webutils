package webutils

import "encoding/base64"

// Base64Decode decodes a base64 string and returns the relevant de-padded byte slice
func Base64Decode(str string) []byte {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(str)))
	n, _ := base64.StdEncoding.Decode(base64Text, []byte(str))

	return base64Text[:n]
}
