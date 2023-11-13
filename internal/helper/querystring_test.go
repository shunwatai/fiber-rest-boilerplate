package helper

import (
	"reflect"
	"testing"
)

type orderByTests struct {
	name  string
	input string
	want  map[string]string
}

func TestParseOrderBy(t *testing.T) {
	tests := []orderByTests{
		{name: "orderBy is empty string", input: "", want: map[string]string{"key": "id", "by": "desc"}},
		{name: "orderBy is name.asc", input: "name.asc", want: map[string]string{"key": "name", "by": "asc"}},
		{name: "orderBy is name", input: "name", want: map[string]string{"key": "name", "by": "desc"}},
		{name: "orderBy is id.asc", input: "id.asc", want: map[string]string{"key": "id", "by": "asc"}},
		{name: "orderBy is userId.asc", input: "userId.asc", want: map[string]string{"key": "user_id", "by": "asc"}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := parseOrderBy(testCase.input)

			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}
