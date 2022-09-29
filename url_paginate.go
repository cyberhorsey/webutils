package webutils

import (
	"net/url"
	"strconv"
)

const (
	PaginationPageField  = "page"
	PaginationLimitField = "limit"
	PaginationLimitMax   = 100
)

type PaginatedResults struct {
	Page  int
	Limit int
}

func PaginateQuerystring(
	querystring url.Values,
) PaginatedResults {
	result := PaginatedResults{}

	filters := FilterQuerystring(
		querystring,
		[]string{
			PaginationPageField,
			PaginationLimitField,
		},
	)

	if values := filters[PaginationPageField]; len(values) > 0 {
		page, _ := strconv.ParseUint(values[0], 10, 64)
		result.Page = int(page)
	}

	if result.Page < 1 {
		result.Page = 1
	}

	if values := filters[PaginationLimitField]; len(values) > 0 {
		limit, _ := strconv.ParseUint(values[0], 10, 64)
		result.Limit = int(limit)
	}

	if result.Limit < 1 ||
		result.Limit > PaginationLimitMax {
		result.Limit = PaginationLimitMax
	}

	return result
}
