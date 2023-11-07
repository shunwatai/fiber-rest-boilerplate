// +build mariadb

package database

import (
	"golang-api-starter/internal/helper"
	"reflect"
	"strings"
	"testing"
)

type mariadbTests struct {
	name  string
	input map[string]interface{}
	want1 string
	want2 map[string]interface{}
}

func TestMariadbSelectStmtFromQuerystring(t *testing.T) {
	var tableName = "todos"
	var testDb = GetDatabase(tableName)
	tests := []mariadbTests{
		{
			name:  "get by ID",
			input: map[string]interface{}{"id": 2},
			want1: `SELECT * FROM todos WHERE id=:id ORDER BY id desc LIMIT 1 OFFSET 0`,
			want2: map[string]interface{}{"id": 2},
		},
		{
			name:  "get by IDs",
			input: map[string]interface{}{"id": []string{"2", "3"}},
			want1: `SELECT * FROM todos WHERE id IN (:id1,:id2) ORDER BY id desc LIMIT 2 OFFSET 0`,
			want2: map[string]interface{}{"id1": "2", "id2": "3"},
		},
		{
			name:  "get keyword by ILIKE",
			input: map[string]interface{}{"task": "show"},
			want1: `SELECT * FROM todos WHERE task like :task ORDER BY id desc LIMIT 1 OFFSET 0`,
			want2: map[string]interface{}{"task": "%show%"},
		},
		{
			name:  "get keywords by ~~ ANY(xx)",
			input: map[string]interface{}{"task": []string{"show", "stop"}, "page": "1", "items": "5"},
			want1: `SELECT * FROM todos WHERE (lower(task) like :task1 or lower(task) like :task2)  ORDER BY id desc LIMIT 5 OFFSET 0`,
			want2: map[string]interface{}{"task1": "%show%", "task2": "%stop%"},
		},
		{
			name:  "get records by keyword that matches in given ids",
			input: map[string]interface{}{"task": "wan", "id": []string{"13", "15"}, "page": "1", "items": "5"},
			want1: `SELECT * FROM todos WHERE task like :task AND id IN (:id1,:id2) ORDER BY id desc LIMIT 5 OFFSET 0`,
			want2: map[string]interface{}{"task": "%wan%", "id1": "13", "id2": "15"},
		},
		{
			name:  "get records by date range",
			input: map[string]interface{}{"withDateFilter": true, "created_at": "2023-01-01.2023-12-31", "page": "1", "items": "5"},
			want1: `SELECT * FROM todos WHERE created_at >= :created_atFrom AND created_at <= :created_atTo ORDER BY id desc LIMIT 5 OFFSET 0`,
			want2: map[string]interface{}{"created_atFrom": "2023-01-01", "created_atTo": "2023-12-31"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got1, _, got2 := testDb.constructSelectStmtFromQuerystring(testCase.input)

			if eq := reflect.DeepEqual(strings.Fields(testCase.want1), strings.Fields(got1)); !eq {
				t.Errorf("got %q \nwant %q", strings.Fields(got1), strings.Fields(testCase.want1))
			}

			if eq := reflect.DeepEqual(testCase.want2, got2); !eq {
				t.Errorf("got %+v \nwant %+v", got2, testCase.want2)
			}
		})
	}
}
