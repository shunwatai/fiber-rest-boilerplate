package database

import (
	// "golang-api-starter/internal/config"
	"golang-api-starter/internal/config"
	"log"
	"github.com/jmoiron/sqlx"
)

type IDatabase interface {
	Select(map[string]interface{}) *sqlx.Rows
	Save(Records) *sqlx.Rows
	// Update(Records) *sqlx.Rows
	Delete(*[]int64) error
	// GetConnectionInfo() ConnectionInfo
}

type ConnectionInfo struct {
	Driver   string
	Host     *string
	Port     *string
	User     *string
	Pass     *string
	Database *string
}

type Records interface {
	StructToMap() []map[string]interface{}
}

// func GetDbConnection(){
// 	config := config.Cfg
// 	config.LoadEnvVariables()
// }

func GetDatabase(tableName string) IDatabase {
	config := config.Cfg
	config.LoadEnvVariables()
	log.Printf("engin: %+v\n", config.DbConf.Driver)

	if config.DbConf.Driver == "sqlite" {
		connection := config.DbConf.SqliteConf
		return &Sqlite{
			ConnectionInfo: &ConnectionInfo{
				Driver:   config.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}

	if config.DbConf.Driver == "mariadb" {
		connection := config.DbConf.MariadbConf
		return &MariaDb{
			ConnectionInfo: &ConnectionInfo{
				Driver:   config.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}

	if config.DbConf.Driver == "postgres" {
		connection := config.DbConf.PostgresConf
		return &Postgres{
			ConnectionInfo: &ConnectionInfo{
				Driver:   config.DbConf.Driver,
				Host:     connection.Host,
				Port:     connection.Port,
				User:     connection.User,
				Pass:     connection.Pass,
				Database: connection.Database,
			},
			TableName: tableName,
		}
	}
	return nil
}
