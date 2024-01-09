package todoDocument

import (
	"fmt"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	"golang-api-starter/internal/modules/document"

	"golang.org/x/exp/maps"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}

// cascadeFields for joining other module, see the example in internal/modules/todo/repository.go
func cascadeFields(todoDocuments TodoDocuments) {
	// cascade documents
	var (
		documentIds []string
		documentId  string
	)
	// get all userIds
	for _, todoDocument := range todoDocuments {
		if todoDocument.DocumentId == nil {
			continue
		}

		documentId = todoDocument.GetDocumentId()
		documentIds = append(documentIds, documentId)
	}

	if len(documentIds) > 0 {
		documents := []*document.Document{}

		// get users by userIds
		condition := helper.GetIdMap(documentIds)
		documents, _ = document.Srvc.Get(condition)
		// get the map[userId]user
		documentMap := document.Srvc.GetIdMap(documents)

		for _, todoDocument := range todoDocuments {
			if todoDocument.DocumentId == nil {
				continue
			}
			document := &document.Document{}
			// take out the user by userId in map and assign
			document = documentMap[todoDocument.GetDocumentId()]
			todoDocument.Document = document
		}
	}
}

func (r *Repository) Get(queries map[string]interface{}) ([]*TodoDocument, *helper.Pagination) {
	fmt.Printf("todoDocument repo\n")
	defaultExactMatch := map[string]bool{
		"id":   true,
		"_id":  true,
		//"done": true, // bool match needs exact match, param can be 0(false) & 1(true)
	}
	if queries["exactMatch"] != nil {
		maps.Copy(queries["exactMatch"].(map[string]bool), defaultExactMatch)
	}

	queries["columns"] = TodoDocument{}.getTags()
	rows, pagination := r.db.Select(queries)

	var records TodoDocuments
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	// records.printValue()

	cascadeFields(records)

	return records, pagination
}

func (r *Repository) Create(todoDocuments []*TodoDocument) ([]*TodoDocument, error) {
	for _, todoDocument := range todoDocuments {
		fmt.Printf("todoDocument repo add: %+v\n", todoDocument)
	}
	rows, err := r.db.Save(TodoDocuments(todoDocuments))

	var records TodoDocuments
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Update(todoDocuments []*TodoDocument) ([]*TodoDocument, error) {
	fmt.Printf("todoDocument repo update\n")
	rows, err := r.db.Save(TodoDocuments(todoDocuments))

	var records TodoDocuments
	if rows != nil {
		records = records.rowsToStruct(rows)
	}
	records.printValue()

	return records, err
}

func (r *Repository) Delete(ids []string) error {
	fmt.Printf("todoDocument repo delete\n")
	err := r.db.Delete(ids)
	if err != nil {
		return err
	}

	return nil
}
