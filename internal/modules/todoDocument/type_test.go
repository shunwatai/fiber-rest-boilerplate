package todoDocument

import (
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/helper/utils"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetId(t *testing.T) {
	todoDocument := &TodoDocument{
		MongoId: utils.ToPtr("xxxx-xxxx-xxxx-xxxx"),
		Id:      utils.ToPtr(int64(2)),
	}

	tests := []struct {
		name     string
		dbDriver string
		input    *TodoDocument
		want     string
	}{
		{name: "test Id", dbDriver: "postgres-or-mariadb-or-sqlite", input: todoDocument, want: "2"},
		{name: "test MongoId", dbDriver: "mongodb", input: todoDocument, want: "xxxx-xxxx-xxxx-xxxx"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cfg.LoadEnvVariables()
			cfg.Vpr.Set("database.engine", testCase.dbDriver)
			if err := cfg.Vpr.Unmarshal(cfg); err != nil {
				t.Logf("failed loading conf, err: %+v\n", err.Error())
			}

			got := todoDocument.GetId()
			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	var id int64 = 2
	now := time.Now()
	customDatetime := &helper.CustomDatetime{&now, utils.ToPtr(time.RFC3339)}
	timeStr, _ := customDatetime.MarshalJSON()
	timeJson := strings.Replace(string(timeStr), "\"", "", -1)
	todoDocuments := TodoDocuments{
		&TodoDocument{
			MongoId:    utils.ToPtr("xxxx-xxxx-xxxx-xxxx"),
			Id:         &id,
			TodoId:     float64(1),
			DocumentId: float64(2),
			CreatedAt:  customDatetime,
			UpdatedAt:  customDatetime,
		},
		&TodoDocument{
			MongoId:    utils.ToPtr("xxxx-xxxx-xxxx-xxxx"),
			Id:         &id,
			TodoId:     "yxyxyxyxyxyxyxyx",
			DocumentId: "zxzxzxzxzxzxzxzx",
			CreatedAt:  customDatetime,
			UpdatedAt:  customDatetime,
		},
	}

	tests := []struct {
		name  string
		input TodoDocuments
		want  []map[string]interface{}
	}{
		{name: "test StructToMap", input: todoDocuments, want: []map[string]interface{}{
			{"_id": "xxxx-xxxx-xxxx-xxxx", "id": float64(2), "created_at": timeJson, "updated_at": timeJson, "todo_id": float64(1), "document_id": float64(2), "document": nil},
			{"_id": "xxxx-xxxx-xxxx-xxxx", "id": float64(2), "created_at": timeJson, "updated_at": timeJson, "todo_id": "yxyxyxyxyxyxyxyx", "document_id": "zxzxzxzxzxzxzxzx", "document": nil},
		}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := todoDocuments.StructToMap()
			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	todoDocuments := TodoDocuments{
		&TodoDocument{},
	}

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "test get db tags", input: "db", want: []string{"id", "todo_id", "document_id", "created_at", "updated_at"}},
		{name: "test get bson tags", input: "bson", want: []string{"_id", "id", "todo_id", "document_id", "created_at", "updated_at"}},
		{name: "test get json tags", input: "json", want: []string{"_id", "id", "todoId", "documentId", "document", "createdAt", "updatedAt"}},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := todoDocuments.GetTags(testCase.input)
			eq := reflect.DeepEqual(testCase.want, got)

			if !eq {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}
