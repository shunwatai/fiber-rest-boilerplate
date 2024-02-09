package main

import (
	"fmt"
	"golang-api-starter/cmd/dbmigrate"
	"golang-api-starter/cmd/gen"
	"golang-api-starter/cmd/server"
	"os"
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
	// default run the api server
	if len(os.Args) == 1 {
		api := server.Api
		api.GetApp()
		api.LoadMiddlewares()
		api.LoadSwagger()
		api.LoadAllRoutes()
		api.Start()
	}

	// run db migration or generate new module if args is given
	if os.Args[1] == "migrate-up" {
		dbmigrate.DbMigrate("up")
	} else if os.Args[1] == "migrate-down" {
		dbmigrate.DbMigrate("down")
	} else if os.Args[1] == "generate" {
		gen.GenerateNewModule()
	}

	fmt.Printf("do nothing...\n")
}
