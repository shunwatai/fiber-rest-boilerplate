package todo

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/database"
)

var tableName = "todos"
var db = database.GetDatabase(tableName)
var Repo = NewRepository(db)
var Srvc = NewService(Repo)
var ctrl = NewController(Srvc)

func GetRoutes(router fiber.Router) {
	r := router.Group("/todo")
	r.Get("/", GetAll)
	r.Post("/", Create)
	r.Patch("/", Update)
	r.Delete("/", Delete)

	rById := router.Group("/todo/:id")
	rById.Get("/", GetById)
}

// TodoGetAll godoc
//
//	@Summary		List Todos
//	@Description	get Todos
//	@Tags			Todos
//	@Accept			json
//	@Produce		json
//	@Param			id		query	number	false	"id"							example(2)
//	@Param			task	query	string	false	"search by task"				example(go practice)
//	@Param			page	query	string	false	"page number for pagination"	example(1)
//	@Param			items	query	string	false	"items per page for pagination"	example(10)
//	@Security		ApiKeyAuth
//	@Router			/todos [get]
func GetAll(c *fiber.Ctx) error {
	return ctrl.Get(c)
}

func GetById(c *fiber.Ctx) error {
	return ctrl.GetById(c)
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
