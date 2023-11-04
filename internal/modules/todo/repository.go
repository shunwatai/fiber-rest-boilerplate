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
	rows, pagination := r.db.Select(queries)

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	return records,pagination
}

func (r *Repository) Create(todos []*Todo) []*Todo {
	fmt.Printf("todo repo add\n")
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

func (r *Repository) Delete(ids *[]int64) ([]*Todo, error) {
	idsString, _ := helper.ConvertNumberSliceToString(*ids)
	rows, _ := r.db.Select(map[string]interface{}{"id": idsString})

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	err := r.db.Delete(ids)
	if err != nil {
		return nil, err
	}

	return records, nil
}
