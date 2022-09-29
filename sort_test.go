package webutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_SortMultiple(t *testing.T) {
	type sortable struct {
		Name      string
		CreatedAt time.Time
	}

	s1 := &sortable{
		Name:      "aaa",
		CreatedAt: time.Unix(1000, 0),
	}
	s2 := &sortable{
		Name:      "aaa",
		CreatedAt: time.Unix(1001, 0),
	}
	s3 := &sortable{
		Name:      "aab",
		CreatedAt: time.Unix(1000, 0),
	}

	sortByCreatedAsc := func(results []interface{}) func(i, j int) int {
		return func(i, j int) int {
			ir := results[i].(*sortable)
			jr := results[j].(*sortable)

			if ir.CreatedAt.Before(jr.CreatedAt) {
				return 1
			}

			if ir.CreatedAt.Equal(jr.CreatedAt) {
				return 0
			}

			return -1
		}
	}
	sortByCreatedDesc := func(results []interface{}) func(i, j int) int {
		return func(i, j int) int {
			ir := results[i].(*sortable)
			jr := results[j].(*sortable)

			if ir.CreatedAt.Before(jr.CreatedAt) {
				return -1
			}

			if ir.CreatedAt.Equal(jr.CreatedAt) {
				return 0
			}

			return 1
		}
	}
	sortByNameAsc := func(results []interface{}) func(i, j int) int {
		return func(i, j int) int {
			ir := results[i].(*sortable)
			jr := results[j].(*sortable)

			if ir.Name < jr.Name {
				return 1
			}

			if ir.Name == jr.Name {
				return 0
			}

			return -1
		}
	}
	sortByNameDesc := func(results []interface{}) func(i, j int) int {
		return func(i, j int) int {
			ir := results[i].(*sortable)
			jr := results[j].(*sortable)

			if ir.Name < jr.Name {
				return -1
			}

			if ir.Name == jr.Name {
				return 0
			}

			return 1
		}
	}

	tests := []struct {
		name   string
		values []interface{}
		sortBy []func(values []interface{}) func(i int, j int) int
		want   []interface{}
	}{
		{
			"empty",
			nil,
			nil,
			nil,
		},
		{
			"created_at,name",
			[]interface{}{s1, s2, s3},
			[]func(values []interface{}) func(i int, j int) int{
				sortByCreatedAsc, sortByNameAsc,
			},
			[]interface{}{s1, s3, s2},
		},
		{
			"created_at,-name",
			[]interface{}{s1, s2, s3},
			[]func(values []interface{}) func(i int, j int) int{
				sortByCreatedAsc, sortByNameDesc,
			},
			[]interface{}{s3, s1, s2},
		},
		{
			"name,created_at",
			[]interface{}{s1, s2, s3},
			[]func(values []interface{}) func(i int, j int) int{
				sortByNameAsc, sortByCreatedAsc,
			},
			[]interface{}{s1, s2, s3},
		},
		{
			"name,-created_at",
			[]interface{}{s1, s2, s3},
			[]func(values []interface{}) func(i int, j int) int{
				sortByNameAsc, sortByCreatedDesc,
			},
			[]interface{}{s2, s1, s3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortMultiple(tt.values, tt.sortBy)
			assert.Equal(t, tt.values, tt.want)
		})
	}
}
