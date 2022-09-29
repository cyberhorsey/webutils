package webutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	Username               string `json:"username" validate:"required_without=Email"`
	Email                  string `json:"email" validate:"required_without=Username"`
	FirstName              string `json:"first_name" validate:"required"`
	LastName               string `json:"last_name" validate:"required"`
	Type                   string `json:"type"`
	Place                  string `json:"place"`
	Language               string `json:"language" validate:"required_if=Type human Place earth"`
	Age                    int    `json:"age" validate:"gt=15"`
	Cats                   int    `json:"cats" validate:"lte=3"`
	FavouritePrimaryColour string `json:"favourite_primary_colour" validate:"oneof=red yellow blue"`
}

func Test_Validate(t *testing.T) {
	tests := []struct {
		name        string
		data        Data
		wantErrMsgs []string
	}{
		{
			"invalidFields",
			Data{
				Type:  "human",
				Place: "earth",
				Cats:  5,
			},
			[]string{
				"username is required unless email is provided",
				"email is required unless username is provided",
				"first_name is required",
				"last_name is required",
				"language is required if type is human, place is earth",
				"age must be greater than 15",
				"cats must be 3 or less",
				"favourite_primary_colour must be one of [red yellow blue]",
			},
		},
		{
			"noErrors",
			Data{
				Email:                  "bdole@example.com",
				FirstName:              "Bob",
				LastName:               "Dole",
				Type:                   "human",
				Language:               "English",
				Age:                    45,
				FavouritePrimaryColour: "red",
			},
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Validate(tt.data)
			gotErrMsgs := make([]string, len(errs))
			for i, err := range errs {
				gotErrMsgs[i] = err.Error()
			}
			assert.Equal(t, tt.wantErrMsgs, gotErrMsgs)
		})
	}
}

func Test_underscore(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"thisIsACamelString",
			"this_is_a_camel_string",
		},
		{
			"ThisIsAlsoACamelString",
			"this_is_also_a_camel_string",
		},
		{
			"CAPITAL",
			"capital",
		},
		{
			"Camel",
			"camel",
		},
		{
			"already_underscored",
			"already_underscored",
		},
		{
			"BIGCamel",
			"big_camel",
		},
		{
			"Spaced Out",
			"spaced_out",
		},
		{
			"spaced out",
			"spaced_out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := underscore(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
