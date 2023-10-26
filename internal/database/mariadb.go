package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MariaDb struct {
	*ConnectionInfo
	TableName string
}

func (m *MariaDb) Select(queries map[string]interface{}) *sqlx.Rows {
	fmt.Printf("select from MariaDB, table: %+v\n", m.TableName)
	return nil
}
func (m *MariaDb) Save(records Records) *sqlx.Rows {
	fmt.Printf("save from MariaDB, table: %+v\n", m.TableName)
	return nil
}
// func (m *MariaDb) Update() {
// 	fmt.Printf("update from MariaDB, table: %+v\n", m.TableName)
// }
func (m *MariaDb) Delete(ids *[]int64) error {
	fmt.Printf("delete from MariaDB, table: %+v\n", m.TableName)
	return nil
}
