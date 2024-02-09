package document

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	//"golang-api-starter/internal/modules/user"
	"golang.org/x/exp/maps"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}

// cascadeFields for joining other module, see the example in internal/modules/document/repository.go
func cascadeFields(documents Documents) {
	// cascade user
}

func (r *Repository) Get(queries map[string]interface{}) ([]*Document, *helper.Pagination) {
	fmt.Printf("document repo\n")
	defaultExactMatch := map[string]bool{
		"id":   true,
		"_id":  true,
		//"done": true, // bool match needs exact match, param can be 0(false) & 1(true)
	}
	if queries["exactMatch"] != nil {
		maps.Copy(queries["exactMatch"].(map[string]bool), defaultExactMatch)
	}

	queries["columns"] = Document{}.getTags()
	rows, pagination := r.db.Select(queries)

	var records Documents
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	//cascadeFields(records)

	return records, pagination
}

func (r *Repository) Create(documents []*Document) ([]*Document, error) {
	for _, document := range documents {
		fmt.Printf("document repo add: %+v\n", document)
	}
	rows, err := r.db.Save(Documents(documents))

	var records Documents
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Update(documents []*Document) ([]*Document, error) {
	fmt.Printf("document repo update\n")
	rows, err := r.db.Save(Documents(documents))

	var records Documents
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Delete(ids []string) error {
	fmt.Printf("document repo delete\n")
	err := r.db.Delete(ids)
	if err != nil {
		return err
	}

	return nil
}
