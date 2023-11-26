package todo

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}

func (r *Repository) Get(queries map[string]interface{}) ([]*Todo, *helper.Pagination) {
	fmt.Printf("todo repo\n")
	queries["exactMatch"] = map[string]bool{
		"id":   true,
		"_id":  true,
		"done": true, // bool match needs exact match, parram can be 0(false) & 1(true)
	}

	cfg.LoadEnvVariables()
	if cfg.DbConf.Driver == "mongodb" {
		queries["columns"] = Todo{}.getTags("bson")
	} else {
		queries["columns"] = Todo{}.getTags("db") // TODO: use this to replace GetColumns()
	}
	rows, pagination := r.db.Select(queries)

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	return records, pagination
}

func (r *Repository) Create(todos []*Todo) []*Todo {
	for _, todo := range todos {
		fmt.Printf("todo repo add: %+v\n", todo)
	}
	rows := r.db.Save(Todos(todos))

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records
}

func (r *Repository) Update(todos []*Todo) []*Todo {
	fmt.Printf("todo repo update\n")
	rows := r.db.Save(Todos(todos))

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records
}

func (r *Repository) Delete(ids []string) error {
	err := r.db.Delete(ids)
	if err != nil {
		return err
	}

	return nil
}
