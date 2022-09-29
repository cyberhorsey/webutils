package webutils

import "sort"

func SortMultiple(
	values []interface{},
	sortBy []func(values []interface{}) func(i, j int) int,
) {
	sort.Slice(
		values,
		func(i, j int) bool {
			for _, sorter := range sortBy {
				comparer := sorter(values)
				compared := comparer(i, j)

				if compared == 1 {
					return true
				}

				if compared == -1 {
					return false
				}
			}

			return false
		},
	)
}
