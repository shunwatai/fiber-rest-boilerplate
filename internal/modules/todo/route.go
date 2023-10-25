package todo

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/database"
)

// var db = &database.Postgre{
// 	Host:      "localhost",
// 	Port:      "3306",
// 	User:      "user",
// 	Pass:      "maria",
// 	TableName: "todo",
// }
var db = &database.Sqlite{
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

var tableName = "todos"
var repo = NewRepository(db)
var srvc = NewService(repo)
var ctrl = NewController(srvc)

func GetRoutes(router fiber.Router) {
	r := router.Group("/todo")

	r.Get("/", GetAll)
	r.Post("/", Create)
	r.Patch("/", Update)
	r.Delete("/", Delete)
}

func GetAll(c *fiber.Ctx) error {
	return ctrl.Get(c)
}

func Create(c *fiber.Ctx) error {
	return ctrl.Create(c)
}

func Update(c *fiber.Ctx) error {
	return ctrl.Update(c)
}

func Delete(c *fiber.Ctx) error {
	return ctrl.Delete(c)
}
