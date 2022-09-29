package webutils

import (
	"bytes"
	"encoding/json"
)

func ResponseFields(
	v interface{},
	fields []string,
) (
	interface{},
	error,
) {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var response interface{}

	decoder := json.NewDecoder(bytes.NewBuffer(bs))
	decoder.UseNumber()

	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	if filter, ok := response.(map[string]interface{}); ok {
		filtered := make(map[string]interface{})

		for _, field := range fields {
			if value, ok := filter[field]; ok {
				filtered[field] = value
			}
		}

		return filtered, nil
	}
	// t,f,n,[,",-,0123456789
	return response, nil
}
