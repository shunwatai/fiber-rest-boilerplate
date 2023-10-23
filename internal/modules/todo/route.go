package todo

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/iancoleman/strcase"
	"golang-api-starter/internal/database"
	"log"
	"net/url"
)

var tableName = "todos"

func GetRoutes(router fiber.Router) {
	// db := &database.Postgre{
	// 	Host:      "localhost",
	// 	Port:      "3306",
	// 	User:      "user",
	// 	Pass:      "maria",
	// 	TableName: "todo",
	// }
	db := &database.Sqlite{
		ConnectionInfo: &database.ConnectionInfo{
			Driver:   "sqlite",
			Host:     "localhost",
			Port:     "",
			User:     "user",
			Pass:     "user",
			Database: "fiber-starter",
		},
		TableName: tableName,
	}

	repo := NewRepository(db)
	srvc := NewService(repo)
	ctrl := NewController(srvc)

	r := router.Group("/todo")

	r.Get("/", func(c *fiber.Ctx) error {
		queries := c.Queries()

		params, err := url.ParseQuery(string(c.Request().URI().QueryString()))
		if err != nil {
			log.Printf("ParseQuery err: %+v\n", err.Error())
		}
		fmt.Printf("queries: %+v\n", queries)

		var paramsMap = make(map[string]interface{}, 0)

		for key, value := range params {
			// fmt.Printf("  %v = %v\n", key, value)
			fmt.Printf("  %v = %v\n", key, value)
			snakeCase := strcase.ToSnake(key)
			if len(value) == 1 {
				paramsMap[snakeCase] = value[0]
				continue
			}
			paramsMap[snakeCase] = value
		}

		// if paramsMap["page"] != nil && paramsMap["items"] != nil {
		// 	pagination.Page, _ = strconv.ParseInt(paramsMap["page"].(string), 10, 64)
		// 	pagination.Items, _ = strconv.ParseInt(paramsMap["items"].(string), 10, 64)
		// }
		//
		// if paramsMap["order_by"] != nil {
		// 	pagination.OrderBy = parseOrderBy(paramsMap["order_by"].(string))
		// }

		fmt.Printf("test: %+v\n", paramsMap)
		results := ctrl.Get(paramsMap)
		// return c.JSON(map[string]interface{}{"message": "todos"})
		return c.JSON(map[string]interface{}{"data": results})
	})

	r.Post("/", func(c *fiber.Ctx) error {
		todo := &Todo{}
		todos := []*Todo{}
		todoErr, todosErr := c.BodyParser(todo), c.BodyParser(&todos)
		if todosErr != nil {
			log.Printf("BodyParser err: %+v\n", todosErr.Error())
		}

		if todoErr == nil {
			todos = append(todos, todo)
		}
		fmt.Printf("save todos: %+v\n", todos)

		results := ctrl.Create(todos)

		if todoErr == nil {
			return c.JSON(map[string]interface{}{"data": results[0]})
		}
		return c.JSON(map[string]interface{}{"data": results})
	})

	r.Patch("/", func(c *fiber.Ctx) error {
		todo := &Todo{}
		todos := []*Todo{}
		todoErr, todosErr := c.BodyParser(todo), c.BodyParser(&todos)
		if todosErr != nil {
			log.Printf("BodyParser err: %+v\n", todosErr.Error())
		}

		if todoErr == nil {
			todos = append(todos, todo)
		}
		fmt.Printf("update todos: %+v\n", todos)

		results := ctrl.Update(todos)

		if todoErr == nil {
			return c.JSON(map[string]interface{}{"data": results[0]})
		}
		return c.JSON(map[string]interface{}{"data": results})
	})

	r.Delete("/", func(c *fiber.Ctx) error {
		ctrl.Delete()
		return nil
	})
}
