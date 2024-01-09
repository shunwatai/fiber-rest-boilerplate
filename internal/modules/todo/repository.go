package todo

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/modules/user"
	"golang.org/x/exp/maps"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}

func cascadeFields(todos Todos) {
	// cascade user
	var (
		userIds []string
		userId  string
	)
	// get all userIds
	for _, todo := range todos {
		if todo.UserId == nil {
			continue
		}

		userId = todo.GetUserId()
		userIds = append(userIds, userId)
	}

	if len(userIds) > 0 {
		users := []*user.User{}

		// get users by userIds
		condition := helper.GetIdMap(userIds)
		users, _ = user.Srvc.Get(condition)
		// get the map[userId]user
		userMap := user.Srvc.GetIdMap(users)

		for _, todo := range todos {
			if todo.UserId == nil {
				continue
			}
			user := &user.User{}
			// take out the user by userId in map and assign
			user = userMap[todo.GetUserId()]
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
	defaultExactMatch := map[string]bool{
		"id":   true,
		"_id":  true,
		"done": true, // bool match needs exact match, parram can be 0(false) & 1(true)
	}
	if queries["exactMatch"] != nil {
		maps.Copy(queries["exactMatch"].(map[string]bool), defaultExactMatch)
	}

	queries["columns"] = Todo{}.getTags()
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

func (r *Repository) Delete(ids []string) error {
	fmt.Printf("todo repo delete\n")
	err := r.db.Delete(ids)
	if err != nil {
		return err
	}

	return nil
}
