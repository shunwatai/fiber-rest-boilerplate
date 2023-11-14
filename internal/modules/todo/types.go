package todo

import (
	"encoding/json"
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

type Todo struct {
	MongoId   *string                `json:"_id,omitempty" bson:"_id,omitempty"` // https://stackoverflow.com/a/20739427
	Id        *int64                 `json:"id,omitempty" db:"id" bson:"id,omitempty"`
	Task      string                 `json:"task" db:"task" bson:"task,omitempty"`
	Done      bool                   `json:"done" db:"done" bson:"done,omitempty"`
	CreatedAt *helper.CustomDatetime `json:"createdAt" db:"created_at" bson:"created_at,omitempty"`
	UpdatedAt *helper.CustomDatetime `json:"updatedAt" db:"updated_at" bson:"updated_at,omitempty"`
	// CreatedAt *string `db:"created_at" json:"createdAt,omitempty"`
	// UpdatedAt *string `db:"updated_at" json:"updatedAt,omitempty"`
}

type Todos []*Todo

func (todos Todos) StructToMap() []map[string]interface{} {
	mapsResults := []map[string]interface{}{}
	for _, todo := range todos {
		tmp := map[string]interface{}{}
		result := map[string]interface{}{}
		data, _ := json.Marshal(todo)
		json.Unmarshal(data, &tmp)
		for k, v := range tmp {
			result[strcase.ToSnake(k)] = v
		}
		mapsResults = append(mapsResults, result)
	}

	return mapsResults
}

// func (todos Todos) rowsToStruct(rows *sqlx.Rows) []*Todo {
func (todos Todos) rowsToStruct(rows database.Rows) []*Todo {
	defer rows.Close()

	records := make([]*Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.StructScan(&todo)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		records = append(records, &todo)
	}

	return records
}

func (todos *Todos) printValue() {
	for _, v := range *todos {
		if v.Id != nil {
			fmt.Printf("existing --> id: %+v, v: %+v\n", *v.Id, *v)
		}
		fmt.Printf("new --> v: %+v\n", *v)
	}
}

// get the tags by key(json / db / bson) name from the struct
// ref: https://stackoverflow.com/a/40865028
func (todo Todo) getTags(key string) []string {
	cols := []string{}
	val := reflect.ValueOf(todo)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		switch jsonTag := t.Tag.Get(key); jsonTag {
		case "-":
		case "":
			// fmt.Println(fieldName)
		default:
			parts := strings.Split(jsonTag, ",")
			name := parts[0]
			if name == "" {
				name = fieldName
			}
			fmt.Println(name)
			cols = append(cols, name)
		}
	}
	return cols
}
