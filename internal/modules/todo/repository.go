package todo

import (
	"fmt"
	"golang-api-starter/internal/database"
	"log"
)

type Repository struct {
	db database.IDatabase
}

func NewRepository(db database.IDatabase) *Repository {
	return &Repository{db}
}

func handleSqlite() {
}

func (r *Repository) Get(queries map[string]interface{}) []*Todo {
	fmt.Printf("todo repo\n")
	todos := []*Todo{}
	rows := r.db.Select(queries)
	defer rows.Close()

	for rows.Next() {
		var todo Todo

		err := rows.StructScan(&todo)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}

		todos = append(todos, &todo)
	}
	fmt.Printf("todos??: %+v\n", todos)

	return todos
}

func (r *Repository) Create(todos []*Todo) []*Todo {
	fmt.Printf("todo repo add\n")
	rows := r.db.Save(Todos(todos))
	defer rows.Close()

	savedRecords := make([]*Todo, 0)
	for rows.Next() {
		var todo Todo
		err := rows.StructScan(&todo)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		savedRecords = append(savedRecords, &todo)
	}
	fmt.Printf("savedRecords??: %+v\n", savedRecords)

	return savedRecords
}
func (r *Repository) Update() {
	r.db.Update()
}
func (r *Repository) Delete() {
	r.db.Delete()
}
