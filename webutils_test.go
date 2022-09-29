package webutils

import "regexp"

func formatJSONString(s string) string {
	formatted := regexp.MustCompile(`[\n\t]`).ReplaceAllString(s, "")
	return formatted
}
