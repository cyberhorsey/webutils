package webutils

import (
	"net/url"
	"strings"
)

const (
	SortByField = "sort_by"
)

func SortByQuerystring(
	querystring url.Values,
	fields []string,
) []string {
	results := []string{}

	if fieldValues, ok := querystring[SortByField]; ok {
		for _, fieldValue := range parseSortValues(fieldValues) {
			for _, field := range fields {
				if field == fieldValue {
					results = append(results, field)
					break
				}
			}
		}
	}

	return results
}

func parseSortValues(
	sortValues []string,
) []string {
	exists := make(map[string]struct{})
	results := []string{}

	for _, values := range sortValues {
		for _, field := range parseSortBy(values) {
			if len(field) > 0 && field[0] == '-' {
				if _, ok := exists[field[1:]]; ok {
					// already being sorted
					continue
				}

				exists[field[1:]] = struct{}{}
			} else {
				if _, ok := exists[field]; ok {
					// already being sorted
					continue
				}

				exists[field] = struct{}{}
			}

			results = append(results, field)
		}
	}

	return results
}

func parseSortBy(
	sortby string,
) []string {
	results := []string{}

	for _, field := range strings.Split(sortby, ",") {
		if field != "" {
			results = append(results, field)
		}
	}

	return results
}
