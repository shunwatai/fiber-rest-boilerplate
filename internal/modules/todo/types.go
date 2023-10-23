package todo

import (
	"encoding/json"
	"time"

	"github.com/iancoleman/strcase"
)

type TodoRequest struct {
	Task string `json:"task" db:"task"`
	Done bool   `json:"done" db:"done"`
}

type Todo struct {
	Id *int64 `json:"id" db:"id"`
	// *TodoRequest
	Task      string     `json:"task" db:"task"`
	Done      bool       `json:"done" db:"done"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
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
