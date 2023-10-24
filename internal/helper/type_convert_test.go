package helper

import (
	"reflect"
	"testing"
)

type test struct {
	name  string
	input interface{}
	want  []string
}

func TestConvertNumberSliceToString(t *testing.T) {
	tests := []test{
		{name: "[]interface to []string", input: []interface{}{1, 2, 3}, want: []string{"1", "2", "3"}},
		{name: "[]int64 to []string", input: []int64{1, 2, 3}, want: []string{"1", "2", "3"}},
		{name: "[]float64 to []string", input: []float64{1, 2, 3}, want: []string{"1", "2", "3"}},
		{name: "[]int to []string", input: []int{1, 2, 3}, want: []string{"1", "2", "3"}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, _ := ConvertNumberSliceToString(testCase.input)

			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}
