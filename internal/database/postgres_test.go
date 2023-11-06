package database

import (
	"golang-api-starter/internal/helper"
	"reflect"
	"strings"
	"testing"
)

var tableName = "todos"
var testDb = GetDatabase(tableName)

type sqlTests struct {
	name  string
	input map[string]interface{}
	want1 string
	want2 *helper.Pagination
	want3 map[string]interface{}
}

func TestConstructSelectStmtFromQuerystring(t *testing.T) {
	tests := []sqlTests{
		{
			name:  "get by ID",
			input: map[string]interface{}{"id": 2},
			want1: `SELECT * FROM todos WHERE id=:id ORDER BY id desc LIMIT 1 OFFSET 0`,
			// want2: &helper.Pagination{
			// 	Page: 1, Items: 0, Count: 1, OrderBy: map[string]string{"by": "desc", "key": "id"}, TotalPages: 1,
			// },
			want3: map[string]interface{}{"id": 2},
		},
		{
			name:  "get by IDs",
			input: map[string]interface{}{"id": []string{"2", "3"}},
			want1: `SELECT * FROM todos WHERE id IN (:id1,:id2) ORDER BY id desc LIMIT 2 OFFSET 0`,
			want3: map[string]interface{}{"id1": "2", "id2": "3"},
		},
		{
			name:  "get keyword by ILIKE",
			input: map[string]interface{}{"task": "show"},
			want1: `SELECT * FROM todos WHERE task ILIKE :task ORDER BY id desc LIMIT 1 OFFSET 0`,
			want3: map[string]interface{}{"task": "%show%"},
		},
		{
			name:  "get keywords by ~~ ANY(xx)",
			input: map[string]interface{}{"task": []string{"show", "stop"}, "page": "1", "items": "5"},
			want1: `SELECT * FROM todos WHERE lower(task) ~~ ANY(:task) ORDER BY id desc LIMIT 5 OFFSET 0`,
			want3: map[string]interface{}{"task": "{%show%,%stop%}"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got1, _, got3 := testDb.constructSelectStmtFromQuerystring(testCase.input)

			if eq := reflect.DeepEqual(strings.Fields(testCase.want1), strings.Fields(got1)); !eq {
				t.Errorf("got %q want %q", strings.Fields(got1), strings.Fields(testCase.want1))
			}

			// skip testing want2(pagination) because of the variation of the records in DB
			// if eq := reflect.DeepEqual(testCase.want2, got2); !eq {
			// 	t.Errorf("got %+v want %+v", got2, testCase.want2)
			// }

			if eq := reflect.DeepEqual(testCase.want3, got3); !eq {
				t.Errorf("got %+v want %+v", got3, testCase.want3)
			}
		})
	}
}
