package main

import (
	"golang-api-starter/cmd/server"
)

//	@title						Golang Fiber API starter
//	@version					1.0
//	@description				This is a sample API starter by fiber.
//	@host						localhost:7000
//	@BasePath					/api
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func main() {
	api := server.Api
	api.GetApp()
	api.LoadMiddlewares()
	api.LoadSwagger()
	api.LoadAllRoutes()
	api.Start()
}
