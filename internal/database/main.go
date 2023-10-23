package database

import (
	"github.com/jmoiron/sqlx"
)

type IDatabase interface {
	Select(map[string]interface{}) *sqlx.Rows
	Save(Records) *sqlx.Rows
	Update(Records) *sqlx.Rows
	Delete()
	GetConnectionInfo() ConnectionInfo
}

type ConnectionInfo struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

type Records interface {
	StructToMap() []map[string]interface{}
}
