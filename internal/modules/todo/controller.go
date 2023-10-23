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
	for _, todo := range todos {
		// t := time.Now().Format("2006-01-02 15:04:05")
		t := time.Now()
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

func (c *Controller) Update() {
	fmt.Printf("todo ctrl update\n")
	c.service.Update()
}

func (c *Controller) Delete() {
	c.service.Delete()
}
