package dbmigrate

import (
	"database/sql"
	"fmt"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var cfg = config.Cfg

func runMigration(action string, m *migrate.Migrate) error {
	if action == "up" {
		return m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	} else if action == "down" {
		return m.Steps(-1)
	}
	return fmt.Errorf("failed to run migrations....")
}

func DbMigrate(action string) error {
	var err error
	basepath := helper.RootDir()
	// log.Println(strings.Repeat("*", 50))
	// log.Printf("db migrate: %+v\n", basepath)
	// log.Println(strings.Repeat("*", 50))

	cfg.LoadEnvVariables()
	dbConn := database.GetDatabase("")

	if cfg.DbConf.Driver == "postgres" {
		connectionString := dbConn.GetConnectionString()
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("sql.Open error: %+v\n", err)
		}
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatalf("postgres.WithInstance error: %+v\n", err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			fmt.Sprintf("file:///%s/migrations/postgres", basepath),
			"postgres", driver)
		if err != nil {
			log.Fatalf("migrate.NewWithDatabaseInstance error: %+v\n", err)
		}

		err = runMigration(action, m)
		log.Println(strings.Repeat("*", 50))
		ver, dir, err := m.Version()
		log.Printf("migrated success, version: %+v, \n", ver)
		log.Printf("migrated failed, dirty: %+v, \n", dir)
		log.Printf("migrated error, error: %+v, \n", err)
		log.Println(strings.Repeat("*", 50))
	}

	return err
}
