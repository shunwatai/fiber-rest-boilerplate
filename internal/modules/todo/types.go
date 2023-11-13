package todo

import (
	"encoding/json"
	"fmt"
	"golang-api-starter/internal/helper"
	"log"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
)

type Todo struct {
	Id *int64 `json:"id" db:"id"`
	Task      string     `json:"task" db:"task"`
	Done      bool       `json:"done" db:"done"`
	CreatedAt *helper.CustomDatetime `db:"created_at" json:"createdAt"`
	UpdatedAt *helper.CustomDatetime `db:"updated_at" json:"updatedAt"`
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

func (todos Todos) rowsToStruct(rows *sqlx.Rows) []*Todo {
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
