package webutils

func MappedFields(fields []string, fieldsToMappedFields map[string]string) []string {
	responseColumns := []string{}
	columnExists := make(map[string]struct{})

	for _, field := range fields {
		if column, ok := fieldsToMappedFields[field]; ok {
			if _, ok := columnExists[column]; !ok {
				responseColumns = append(responseColumns, column)

				columnExists[column] = struct{}{}
			}
		}
	}

	if len(responseColumns) == 0 {
		return []string{}
	}

	return responseColumns
}

// IncludesAnyReturnField returns true if returnFields is empty or includes at least one of fields
func IncludesAnyReturnField(returnFields []string, fields ...string) bool {
	if len(returnFields) == 0 {
		return true
	}

	for _, field := range fields {
		if ContainsString(returnFields, field) {
			return true
		}
	}

	return false
}

// ReturnFieldsWith appends a field to fields slice if not already appended
func ReturnFieldsWith(fields []string, field string) []string {
	if !IncludesAnyReturnField(fields, field) {
		fields = append(fields, field)
	}

	return fields
}
