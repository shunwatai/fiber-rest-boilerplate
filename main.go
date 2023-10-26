package main

import (
	"golang-api-starter/cmd/server"
)

func main() {
	api := server.Api
	api.GetApp()
	api.LoadMiddlewares()
	api.LoadAllRoutes()
	api.Start()
}
