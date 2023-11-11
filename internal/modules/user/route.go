package user

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/database"
)

var tableName = "users"
var db = database.GetDatabase(tableName)
var Repo = NewRepository(db)
var Srvc = NewService(Repo)
var ctrl = NewController(Srvc)

func GetRoutes(router fiber.Router) {
	r := router.Group("/users")
	r.Get("/", GetAll)
	r.Post("/", Create)
	r.Patch("/", Update)
	r.Delete("/", Delete)

	rById := r.Group("/:id")
	rById.Get("/", GetById)
}

// UserGetAll godoc
//
//	@Summary		List Users
//	@Description	get Users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		query	number	false	"id"							example(2)
//	@Param			task	query	string	false	"search by task"				example(go practice)
//	@Param			done	query	boolean	false	"search by done"				example(1)
//	@Param			page	query	string	false	"page number for pagination"	example(1)
//	@Param			items	query	string	false	"items per page for pagination"	example(10)
//	@Security		ApiKeyAuth
//	@Router			/users [get]
func GetAll(c *fiber.Ctx) error {
	return ctrl.Get(c)
}

// GetUserById godoc
//
//	@Summary		Get User by ID
//	@Description	get User by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path	int	true	"user ID"	example(12)
//	@Security		ApiKeyAuth
//	@Router			/users/{userId} [get]
func GetById(c *fiber.Ctx) error {
	return ctrl.GetById(c)
}

// UserCreate godoc
//
//	@Summary		Create new user(s)
//	@Description	Create user(s)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			User	body	string	true	"single User request json"	SchemaExample({ "task": "take shower", "done": false })
//	@Param			Users	body	string	true	"batch User request json"	SchemaExample([{ "task": "take shower", "done": false }, { "task": "go practice", "done": false }])
//	@Security		ApiKeyAuth
//	@Router			/users [post]
func Create(c *fiber.Ctx) error {
	return ctrl.Create(c)
}

// UserUpdate godoc
//
//	@Summary		Update existing user(s)
//	@Description	Update user(s)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			User	body	string	true	"single User request json"	SchemaExample({ "id":2, "task": "take shower", "done": false })
//	@Param			Users	body	string	true	"batch User request json"	SchemaExample([{ "id":2, "task": "take shower", "done": false, createdAt: "2021-01-11" }, { "id":3, "task": "go practice", "done": false }])
//	@Security		ApiKeyAuth
//	@Router			/users [patch]
func Update(c *fiber.Ctx) error {
	return ctrl.Update(c)
}

// RemoveUser godoc
//
//	@Summary		Delete user(s)
//	@Description	delete user(s)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userIds	body	string	true	"array of user IDs"	SchemaExample([1,2,3])
//	@Security		ApiKeyAuth
//	@Router			/users [delete]
func Delete(c *fiber.Ctx) error {
	return ctrl.Delete(c)
}
