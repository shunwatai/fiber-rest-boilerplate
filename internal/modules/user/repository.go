package user

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

func (r *Repository) Get(queries map[string]interface{}) ([]*User, *helper.Pagination) {
	fmt.Printf("user repo\n")
	rows, pagination := r.db.Select(queries)

	var records Users
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	return records, pagination
}

func (r *Repository) Create(users []*User) ([]*User, error) {
	for _, user := range users {
		fmt.Printf("user repo add: %+v\n", user)
	}
	rows, err := r.db.Save(Users(users))

	var records Users
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Update(users []*User) ([]*User, error) {
	fmt.Printf("user repo update\n")
	rows, err := r.db.Save(Users(users))

	var records Users
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Delete(ids *[]int64) ([]*User, error) {
	idsString, _ := helper.ConvertNumberSliceToString(*ids)
	rows, _ := r.db.Select(map[string]interface{}{"id": idsString})

	var records Users
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
