package todo

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/modules/user"
	"strconv"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}

func cascadeFields(todos Todos) {
	// cascade user
	userIds := []string{}
	for _, todo := range todos {
		if todo.UserId == nil {
			continue
		}
		userId := strconv.Itoa(int(todo.UserId.(int64)))
		userIds = append(userIds, userId)
	}

	if len(userIds) > 0 {
		users, _ := user.Srvc.Get(map[string]interface{}{"id": userIds})
		userMap := user.Srvc.GetIdMap(users)

		for _, todo := range todos {
			if todo.UserId == nil {
				continue
			}
			user := userMap[strconv.Itoa(int(todo.UserId.(int64)))]
			todo.User = user
		}
	}

	// // cascade documents
	// var pagination helper.Pagination
	// docsByExpenseId := document.Service.GetByColumns(map[string]interface{}{"expense_id": e.Id}, &pagination)
	// documents := make([]*models.Document, len(docsByExpenseId))
	//
	// for i, d := range docsByExpenseId {
	// 	// fmt.Printf("docsByExpenseId %+v: %+v\n", *e.Id, d)
	// 	documents[i] = d
	// }
	// e.Documents = documents
}

func (r *Repository) Get(queries map[string]interface{}) ([]*Todo, *helper.Pagination) {
	fmt.Printf("todo repo\n")
	rows, pagination := r.db.Select(queries)

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	cascadeFields(records)

	return records, pagination
}

func (r *Repository) Create(todos []*Todo) ([]*Todo, error) {
	for _, todo := range todos {
		fmt.Printf("todo repo add: %+v\n", todo)
	}
	rows, err := r.db.Save(Todos(todos))

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Update(todos []*Todo) ([]*Todo, error) {
	fmt.Printf("todo repo update\n")
	rows, err := r.db.Save(Todos(todos))

	var records Todos
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Delete(ids []string) ([]*Todo, error) {
	rows, _ := r.db.Select(map[string]interface{}{"id": ids})

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
