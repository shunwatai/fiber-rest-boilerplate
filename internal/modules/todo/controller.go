package todo

import (
	"fmt"
	"time"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) Controller {
	return Controller{s}
}

func (c *Controller) Get(queries map[string]interface{}) []*Todo {
	fmt.Printf("todo ctrl\n")
	results := c.service.Get(queries)

	return results
}

func (c *Controller) Create(todos []*Todo) []*Todo {
	fmt.Printf("todo ctrl create\n")
	t := time.Now()
	for _, todo := range todos {
		// t := time.Now().Format("2006-01-02 15:04:05")
		if todo.CreatedAt == nil {
			todo.CreatedAt = &t
		}
		if todo.UpdatedAt == nil {
			todo.UpdatedAt = &t
		}
	}
	// return []*Todo{}
	return c.service.Create(todos)
}

func (c *Controller) Update(todos []*Todo) []*Todo {
	fmt.Printf("todo ctrl update\n")
	t := time.Now()
	for _, todo := range todos {
		if todo.Id == nil {
			todo.CreatedAt = &t
		}
		todo.UpdatedAt = &t
	}

	return c.service.Update(todos)
}

func (c *Controller) Delete() {
	c.service.Delete()
}
