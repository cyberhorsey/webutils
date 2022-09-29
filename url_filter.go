package webutils

import (
	"net/url"
)

func FilterQuerystring(
	querystring url.Values,
	fields []string,
) map[string][]string {
	result := make(map[string][]string)

	for _, field := range fields {
		if values, ok := querystring[field]; ok {
			result[field] = values
		}
	}

	return result
}
