package main

import (
	"fmt"
	"golang-api-starter/cmd/dbmigrate"
	"golang-api-starter/cmd/gen"
	"golang-api-starter/cmd/server"
	"golang-api-starter/internal/helper/logger/zap_log"
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

	// run db migration or generate new module if args are given
	if os.Args[1] == "migrate-up" || os.Args[1] == "migrate-down" {
		fmt.Printf("db migrate\n")
		if len(os.Args) != 3 {
			logger.Errorf("please provide the target db name for 2nd arg.")
			logger.Fatalf("e.g. go run main.go migrate-[up/down] [postgres/mariadb/sqlite/mongodb]")
		}
		dbmigrate.DbMigrate(os.Args[1], os.Args[2])
	} else if os.Args[1] == "generate" {
		gen.GenerateNewModule()
	} else {
		fmt.Printf("do nothing...\n")
	}
}
