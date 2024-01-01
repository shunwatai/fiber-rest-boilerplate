package document

import (
	"github.com/gofiber/fiber/v2"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/middleware/jwtcheck"
)

var tableName = "documents"
var db = database.GetDatabase(tableName)
var Repo = NewRepository(db)
var Srvc = NewService(Repo)
var ctrl = NewController(Srvc)

func GetRoutes(router fiber.Router) {
	r := router.Group("/documents", jwtcheck.CheckFromHeader())
	r.Get("/", GetAll)
	r.Post("/", Create)
	r.Patch("/", Update)
	r.Delete("/", Delete)

	rById := r.Group("/:id")
	rById.Get("/", GetById)
}

// DocumentGetAll godoc
//
//	@Summary		List Documents
//	@Description	get Documents
//	@Tags			documents
//	@Accept			json
//	@Produce		json
//	@Param			id		query	number	false	"id"							example(2)
//	@Param			userId	query	number	false	"search by userId"				example(2)
//	@Param			task	query	string	false	"search by task"				example(go practice)
//	@Param			done	query	boolean	false	"search by done"				example(1)
//	@Param			page	query	string	false	"page number for pagination"	example(1)
//	@Param			items	query	string	false	"items per page for pagination"	example(10)
//	@Security		ApiKeyAuth
//	@Router			/documents [get]
func GetAll(c *fiber.Ctx) error {
	return ctrl.Get(c)
}

// GetDocumentById godoc
//
//	@Summary		Get Document by ID
//	@Description	get Document by ID
//	@Tags			documents
//	@Accept			json
//	@Produce		json
//	@Param			documentId	path	int	true	"document ID"	example(12)
//	@Security		ApiKeyAuth
//	@Router			/documents/{documentId} [get]
func GetById(c *fiber.Ctx) error {
	return ctrl.GetById(c)
}

// DocumentCreate godoc
//
//	@Summary		Create new document(s)
//	@Description	Create document(s)
//	@Tags			documents
//	@Accept			json
//	@Produce		json
//	@Param			Document	body	string	true	"single Document request json"	SchemaExample({ "task": "take shower", "done": false })
//	@Param			Documents	body	string	true	"batch Document request json"	SchemaExample([{ "task": "take shower", "done": false }, { "task": "go practice", "done": false }])
//	@Security		ApiKeyAuth
//	@Router			/documents [post]
func Create(c *fiber.Ctx) error {
	return ctrl.Create(c)
}

// DocumentUpdate godoc
//
//	@Summary		Update existing document(s)
//	@Description	Update document(s)
//	@Tags			documents
//	@Accept			json
//	@Produce		json
//	@Param			Document	body	string	true	"single Document request json"	SchemaExample({ "id":2, "task": "take shower", "done": false })
//	@Param			Documents	body	string	true	"batch Document request json"	SchemaExample([{ "id":2, "task": "take shower", "done": false, createdAt: "2021-01-11" }, { "id":3, "task": "go practice", "done": false }])
//	@Security		ApiKeyAuth
//	@Router			/documents [patch]
func Update(c *fiber.Ctx) error {
	return ctrl.Update(c)
}

// RemoveDocument godoc
//
//	@Summary		Delete document(s)
//	@Description	delete document(s)
//	@Tags			documents
//	@Accept			json
//	@Produce		json
//	@Param			documentIds	body	string	true	"array of document IDs"	SchemaExample([1,2,3])
//	@Security		ApiKeyAuth
//	@Router			/documents [delete]
func Delete(c *fiber.Ctx) error {
	return ctrl.Delete(c)
}
