package webutils

import (
	"net/url"
	"strings"
)

const (
	FieldsQuerystringField = "fields"
)

func FieldsQuerystring(
	querystring url.Values,
	requestFields map[string]struct{},
) []string {
	responseFields := []string{}
	fieldExists := make(map[string]struct{})

	if values, ok := querystring[FieldsQuerystringField]; ok && len(values) > 0 {
		for _, field := range strings.Split(
			strings.Join(values, ","), // ["a","b","c,d"] -> "a,b,c,d"
			",",                       // a,b,c,d -> ["a","b","c","d"]
		) {
			if _, ok := fieldExists[field]; !ok {
				if _, ok := requestFields[field]; ok {
					responseFields = append(responseFields, field)

					fieldExists[field] = struct{}{}
				}
			}
		}
	}

	return responseFields
}
